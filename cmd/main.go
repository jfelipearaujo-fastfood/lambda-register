package main

import (
	_ "embed"
	"log/slog"
	"os"
	"time"

	"github.com/jfelipearaujo-org/lambda-register/internal/database"
	"github.com/jfelipearaujo-org/lambda-register/internal/handlers"
	"github.com/jfelipearaujo-org/lambda-register/internal/hashs"
	"github.com/jfelipearaujo-org/lambda-register/internal/providers"
	"github.com/jfelipearaujo-org/lambda-register/internal/router"
	"github.com/jfelipearaujo-org/lambda-register/internal/token"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	log := slog.New(handler)

	slog.SetDefault(log)
}

func routerReq(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	slog.Info("received a request", "path", req.Path, "method", req.HTTPMethod)

	timeProvider := providers.NewTimeProvider(time.Now)
	db := database.NewDatabaseFromConnStr(timeProvider)
	hasher := hashs.NewHasher()
	jwt := token.NewToken()

	handler := handlers.NewHandler(db, hasher, jwt)

	if req.Path == "/register" && req.HTTPMethod == "POST" {
		return handler.CrateUser(req)
	}

	return router.MethodNotAllowed(), nil
}

func main() {
	lambda.Start(routerReq)
}
