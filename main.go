package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jsfelipearaujo/fast-food-lambda-register/src/cpf"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	engine = "postgres"
)

var (
	signingKey = []byte(os.Getenv("SIGN_KEY"))

	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")

	connectionStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
)

type Request struct {
	CPF      string `json:"cpf"`
	Password string `json:"pass"`
}

type User struct {
	Id          string `json:"Id"`
	DocumentId  string `json:"DocumentId"`
	Password    string `json:"Password"`
	IsAnonymous bool   `json:"IsAnonymous"`
}

func checkIfCPFIsInUse(cpf string) (bool, error) {
	conn, err := sql.Open(engine, connectionStr)
	if err != nil {
		slog.Error("error while trying to connect to the database", "error", err)
		return false, err
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		slog.Error("error while trying to ping the database", "error", err)
		return false, err
	}

	statement, err := conn.Query("SELECT COUNT(*) As count FROM clients WHERE DocumentId = $1", cpf)
	if err != nil {
		slog.Error("error while trying to execute the query", "error", err)
		return false, err
	}

	var count int
	for statement.Next() {
		if err := statement.Scan(&count); err != nil {
			slog.Error("error while trying to scan the result", "error", err)
			return false, err
		}
	}

	return count > 0, nil
}

func persistUser(user User) error {
	conn, err := sql.Open(engine, connectionStr)
	if err != nil {
		slog.Error("error while trying to connect to the database", "error", err)
		return err
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		slog.Error("error while trying to ping the database", "error", err)
		return err
	}

	var res sql.Result

	if user.IsAnonymous {
		slog.Debug("inserting an anonymous user")
		res, err = conn.Exec("INSERT INTO clients (Id, DocumentType, IsAnonymous) VALUES ($1, $2, $3);",
			user.Id,
			"CPF",
			false)

		if err != nil {
			slog.Error("error while trying to execute the query", "error", err)
			return err
		}
	} else {
		slog.Debug("inserting an user")
		res, err = conn.Exec("INSERT INTO clients (Id, DocumentId, DocumentType, IsAnonymous, Password) VALUES ($1, $2, $3, $4, $5);",
			user.Id,
			user.DocumentId,
			"CPF",
			false,
			user.Password)

		if err != nil {
			slog.Error("error while trying to execute the query", "error", err)
			return err
		}
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		slog.Error("error while trying to get the affected rows", "error", err)
		return err
	}

	slog.Debug("persist user completed", "affected_rows", affectedRows)

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func createJwtToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString(signingKey)
}

func maskCpf(cpf string) string {
	return strings.ReplaceAll(cpf, cpf[3:(len(cpf)-2)], strings.Repeat("*", len(cpf)-5))
}

func handleCreateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request Request

	slog.Info("creating a new user")

	slog.Debug("unmarshalling the request")
	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		slog.Error("error while trying to unmarshal the request", "error", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       http.StatusText(http.StatusBadRequest),
		}, nil
	}

	var user User

	// If the request is not empty, we need to validate the CPF and the Password
	if request.CPF == "" || request.Password == "" {
		cpf := cpf.NewCPF(request.CPF)

		slog.Debug("validating the cpf")
		if !cpf.IsValid() {
			slog.Error("invalid cpf", "cpf", request.CPF)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "Invalid CPF or Password",
			}, nil
		}

		slog.Debug("validating the password")
		if len(request.Password) < 8 {
			slog.Error("invalid password", "password_length", len(request.Password))
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "Invalid CPF or Password",
			}, nil
		}

		slog.Debug("checking if the cpf is in use")
		cpfInUse, err := checkIfCPFIsInUse(request.CPF)
		if err != nil {
			slog.Error("error while trying to check if the cpf is in use", "error", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       http.StatusText(http.StatusInternalServerError),
			}, nil
		}

		slog.Debug("checking if the cpf is in use", "cpf_in_use", cpfInUse)
		if cpfInUse {
			slog.Error("cpf already in use", "cpf", maskCpf(request.CPF))
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusConflict,
				Body:       "Invalid CPF or Password",
			}, nil
		}

		slog.Debug("hashing the password")
		hashedPassword, err := hashPassword(request.Password)
		if err != nil {
			slog.Error("error while trying to hash the password", "error", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       http.StatusText(http.StatusInternalServerError),
			}, nil
		}

		user = User{
			Id:          uuid.NewString(),
			DocumentId:  request.CPF,
			Password:    hashedPassword,
			IsAnonymous: false,
		}
	} else {
		// If the request is empty, we need to create an anonymous user
		user = User{
			Id:          uuid.NewString(),
			IsAnonymous: true,
		}
	}

	slog.Debug("persisting the user")
	if err := persistUser(user); err != nil {
		slog.Error("error while trying to persist the user", "error", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	slog.Debug("creating the jwt token")
	token, err := createJwtToken(user)
	if err != nil {
		slog.Error("error while trying to create the jwt token", "error", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	slog.Debug("completed", "token", token)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       token,
	}, nil
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	slog.Info("received a request", "path", req.Path, "method", req.HTTPMethod)

	if req.Path == "/users" && req.HTTPMethod == "POST" {
		return handleCreateUser(req)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Body:       http.StatusText(http.StatusMethodNotAllowed),
	}, nil
}

func main() {
	lambda.Start(router)
}
