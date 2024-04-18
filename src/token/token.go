package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jfelipearaujo-org/lambda-register/src/entities"
)

var (
	signingKey = []byte(os.Getenv("SIGN_KEY"))
)

type Token struct {
}

func NewToken() Token {
	return Token{}
}

func (t Token) CreateJwtToken(user entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString(signingKey)
}
