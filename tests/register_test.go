package tests

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/jfelipearaujo-org/lambda-register/internal/database"
	"github.com/jfelipearaujo-org/lambda-register/internal/handlers"
	"github.com/jfelipearaujo-org/lambda-register/internal/hashs"
	"github.com/jfelipearaujo-org/lambda-register/internal/providers"
	"github.com/jfelipearaujo-org/lambda-register/internal/token"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var opts = godog.Options{
	Format:      "pretty",
	Paths:       []string{"features"},
	Output:      colors.Colored(os.Stdout),
	Concurrency: 4,
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestFeatures(t *testing.T) {
	t.Setenv("SIGN_KEY", "key")

	o := opts
	o.TestingT = t

	status := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             &o,
	}.Run()

	if status == 2 {
		t.SkipNow()
	}

	if status != 0 {
		t.Fatalf("zero status code expected, %d received", status)
	}
}

// Context keys

type cpfCtxKey struct{}
type passwordCtxKey struct{}
type responseStatusCtxKey struct{}

func setCPF(ctx context.Context, cpf string) context.Context {
	return context.WithValue(ctx, cpfCtxKey{}, cpf)
}

func getCPF(ctx context.Context) string {
	return ctx.Value(cpfCtxKey{}).(string)
}

func setPassword(ctx context.Context, password string) context.Context {
	return context.WithValue(ctx, passwordCtxKey{}, password)
}

func getPassword(ctx context.Context) string {
	return ctx.Value(passwordCtxKey{}).(string)
}

func setResponseStatus(ctx context.Context, status int) context.Context {
	return context.WithValue(ctx, responseStatusCtxKey{}, status)
}

func getResponseStatus(ctx context.Context) int {
	return ctx.Value(responseStatusCtxKey{}).(int)
}

// Steps

type appFeature struct {
	db *sql.DB
}

func (af *appFeature) theUserCPFIs(ctx context.Context, cpf string) (context.Context, error) {
	return setCPF(ctx, cpf), nil
}

func (af *appFeature) theUserPasswordIs(ctx context.Context, password string) (context.Context, error) {
	return setPassword(ctx, password), nil
}

func (af *appFeature) theUserRequestToBeRegistered(ctx context.Context) (context.Context, error) {
	timeProvider := providers.NewTimeProvider(time.Now)
	db := database.NewDatabase(af.db, timeProvider)
	hasher := hashs.NewHasher()
	jwt := token.NewToken()
	handler := handlers.NewHandler(db, hasher, jwt)

	req := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"cpf":"%v","pass":"%v"}`, getCPF(ctx), getPassword(ctx)),
	}

	resp, err := handler.CrateUser(req)
	if err != nil {
		return ctx, err
	}

	return setResponseStatus(ctx, resp.StatusCode), nil
}

func (af *appFeature) theUserShouldBeRegisteredSuccessfully(ctx context.Context) (context.Context, error) {
	responseStatus := getResponseStatus(ctx)
	if responseStatus != http.StatusOK {
		return ctx, fmt.Errorf("expected status code 200, but got %d", responseStatus)
	}

	return ctx, nil
}

var (
	containers = make(map[string]testcontainers.Container)
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	app := &appFeature{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		postgresContainer, err := postgres.RunContainer(ctx,
			testcontainers.WithImage("postgres:16"),
			postgres.WithInitScripts(filepath.Join("testdata", "init-db.sql")),
			postgres.WithDatabase("lambda"),
			postgres.WithUsername("lambda"),
			postgres.WithPassword("123456"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Second)),
		)
		if err != nil {
			return ctx, err
		}

		connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			return ctx, err
		}

		app.db, err = sql.Open("postgres", connStr)
		if err != nil {
			return ctx, err
		}

		if err := app.db.Ping(); err != nil {
			return ctx, err
		}

		containers[sc.Id] = postgresContainer

		return ctx, nil
	})

	ctx.Step(`^the user CPF is "([^"]*)"$`, app.theUserCPFIs)
	ctx.Step(`^the user password is "([^"]*)"$`, app.theUserPasswordIs)
	ctx.Step(`^the user request to be registered$`, app.theUserRequestToBeRegistered)
	ctx.Step(`^the user should be registered successfully$`, app.theUserShouldBeRegisteredSuccessfully)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if err != nil {
			return ctx, err
		}

		if err := containers[sc.Id].Terminate(ctx); err != nil {
			return ctx, err
		}

		return ctx, nil
	})
}
