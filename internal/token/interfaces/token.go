package interfaces

import (
	"github.com/jfelipearaujo-org/lambda-register/internal/entities"
)

type Token interface {
	CreateJwtToken(user entities.User) (string, error)
}
