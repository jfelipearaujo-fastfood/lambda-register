package interfaces

import (
	"github.com/jsfelipearaujo/lambda-register/src/entities"
)

type Token interface {
	CreateJwtToken(user entities.User) (string, error)
}
