package router

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jfelipearaujo-org/lambda-register/internal/entities"
)

func InvalidRequestBody() events.APIGatewayProxyResponse {
	return buildResponse(http.StatusBadRequest, "error to parse the request body", "")
}

func InvalidCPFOrPassword() events.APIGatewayProxyResponse {
	return buildResponse(http.StatusUnauthorized, "invalid cpf or password", "")
}

func InternalServerError() events.APIGatewayProxyResponse {
	return buildResponse(http.StatusInternalServerError, "internal server error", "")
}

func MethodNotAllowed() events.APIGatewayProxyResponse {
	return buildResponse(http.StatusMethodNotAllowed, "method not allowed", "")
}

func Success(token string) events.APIGatewayProxyResponse {
	return buildResponse(http.StatusOK, "success", token)
}

func buildResponse(status int, message string, token string) events.APIGatewayProxyResponse {
	response := entities.Response{
		Status:      status,
		Message:     message,
		AccessToken: token,
	}

	body, err := json.Marshal(response)
	if err != nil {
		slog.Error("error while trying to marshal the response", "error", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
