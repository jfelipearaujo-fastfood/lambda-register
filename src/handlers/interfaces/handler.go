package interfaces

import (
	"github.com/aws/aws-lambda-go/events"
)

type Handler interface {
	CrateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
