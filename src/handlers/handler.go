package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jsfelipearaujo/lambda-register/src/cpf"
	db_interface "github.com/jsfelipearaujo/lambda-register/src/database/interfaces"
	"github.com/jsfelipearaujo/lambda-register/src/entities"
	hash_interface "github.com/jsfelipearaujo/lambda-register/src/hashs/interfaces"
	"github.com/jsfelipearaujo/lambda-register/src/router"
	token_interface "github.com/jsfelipearaujo/lambda-register/src/token/interfaces"
)

type Handler struct {
	db     db_interface.Database
	hasher hash_interface.Hasher
	jwt    token_interface.Token
}

func NewHandler(
	db db_interface.Database,
	hasher hash_interface.Hasher,
	jwt token_interface.Token,
) Handler {
	return Handler{
		db:     db,
		hasher: hasher,
		jwt:    jwt,
	}
}

func (h Handler) CrateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request entities.Request
	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		return router.InvalidRequestBody(), err
	}

	var user entities.User

	if request.IsAnonymous() {
		user = entities.NewAnonymousUser()
	} else {
		cpf := cpf.NewCPF(request.CPF)

		if !cpf.IsValid() {
			return router.InvalidCPFOrPassword(), nil
		}

		if !request.IsPasswordWithMinimumLength() {
			return router.InvalidCPFOrPassword(), nil
		}

		cpfInUse, err := h.db.CheckIfCPFIsInUse(cpf.String())
		if err != nil {
			return router.InternalServerError(), err
		}

		if cpfInUse {
			return router.InvalidCPFOrPassword(), nil
		}

		hashedPassword, err := h.hasher.HashPassword(request.Password)
		if err != nil {
			return router.InternalServerError(), err
		}

		user = entities.NewUser(cpf.String(), hashedPassword)
	}

	if err := h.db.PersistUser(user); err != nil {
		return router.InternalServerError(), err
	}

	token, err := h.jwt.CreateJwtToken(user)
	if err != nil {
		return router.InternalServerError(), err
	}

	return router.Success(token), nil
}
