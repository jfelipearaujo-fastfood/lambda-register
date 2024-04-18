package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	db_interface "github.com/jfelipearaujo-org/lambda-register/src/database/interfaces"
	db_interface_mock "github.com/jfelipearaujo-org/lambda-register/src/database/interfaces/mocks"
	hash_interface "github.com/jfelipearaujo-org/lambda-register/src/hashs/interfaces"
	hash_interface_mock "github.com/jfelipearaujo-org/lambda-register/src/hashs/interfaces/mocks"
	token_interface "github.com/jfelipearaujo-org/lambda-register/src/token/interfaces"
	token_interface_mock "github.com/jfelipearaujo-org/lambda-register/src/token/interfaces/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHandler(t *testing.T) {
	type args struct {
		db     db_interface.Database
		hasher hash_interface.Hasher
		jwt    token_interface.Token
	}
	tests := []struct {
		name string
		args args
		want Handler
	}{
		{
			name: "Should return a new instance correctly",
			args: args{
				db:     db_interface_mock.NewMockDatabase(t),
				hasher: hash_interface_mock.NewMockHasher(t),
				jwt:    token_interface_mock.NewMockToken(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange

			// Act
			got := NewHandler(tt.args.db, tt.args.hasher, tt.args.jwt)

			// Assert
			assert.IsType(t, tt.want, got)
		})
	}
}

func TestHandler_CrateUser(t *testing.T) {
	t.Run("Should return a success response when creating a non anonymous user", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		db_mock.On("CheckIfCPFIsInUse", "218.486.310-65").
			Return(false, nil).
			Once()

		hasher_mock.On("HashPassword", "12345678").
			Return("abc123", nil).
			Once()

		db_mock.On("PersistUser", mock.AnythingOfType("entities.User")).
			Return(nil).
			Once()

		jwt_mock.On("CreateJwtToken", mock.AnythingOfType("entities.User")).
			Return("token", nil).
			Once()

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"218.486.310-65","pass":"12345678"}`,
		}

		// Act
		got, err := h.CrateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, got.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return a success response when creating a anonymous user", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		db_mock.On("PersistUser", mock.AnythingOfType("entities.User")).
			Return(nil).
			Once()

		jwt_mock.On("CreateJwtToken", mock.AnythingOfType("entities.User")).
			Return("token", nil).
			Once()

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"","pass":""}`,
		}

		// Act
		got, err := h.CrateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, got.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return an error when CPF is invalid", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"123","pass":"12345678"}`,
		}

		// Act
		got, err := h.CrateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, got.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return an error when password is invalid", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"784.655.630-47","pass":"123"}`,
		}

		// Act
		got, err := h.CrateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, got.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return an error when CPF is in use", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		db_mock.On("CheckIfCPFIsInUse", "218.486.310-65").
			Return(true, nil).
			Once()

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"218.486.310-65","pass":"12345678"}`,
		}

		// Act
		got, err := h.CrateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, got.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return an error when something got wrong when check if CPF is in use", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		db_mock.On("CheckIfCPFIsInUse", "218.486.310-65").
			Return(false, errors.New("error")).
			Once()

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"218.486.310-65","pass":"12345678"}`,
		}

		// Act
		got, err := h.CrateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, got.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return an error when something got wrong when the password is hashed", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		db_mock.On("CheckIfCPFIsInUse", "218.486.310-65").
			Return(false, nil).
			Once()

		hasher_mock.On("HashPassword", "12345678").
			Return("abc123", errors.New("error")).
			Once()

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"218.486.310-65","pass":"12345678"}`,
		}

		// Act
		got, err := h.CrateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, got.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return an error when something got wrong when try to persist the user", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		db_mock.On("CheckIfCPFIsInUse", "218.486.310-65").
			Return(false, nil).
			Once()

		hasher_mock.On("HashPassword", "12345678").
			Return("abc123", nil).
			Once()

		db_mock.On("PersistUser", mock.AnythingOfType("entities.User")).
			Return(errors.New("error")).
			Once()

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"218.486.310-65","pass":"12345678"}`,
		}

		// Act
		got, err := h.CrateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, got.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})

	t.Run("Should return an error when something got wrong when generate the token", func(t *testing.T) {
		// Arrange
		db_mock := db_interface_mock.NewMockDatabase(t)
		hasher_mock := hash_interface_mock.NewMockHasher(t)
		jwt_mock := token_interface_mock.NewMockToken(t)

		h := NewHandler(
			db_mock,
			hasher_mock,
			jwt_mock,
		)

		db_mock.On("CheckIfCPFIsInUse", "218.486.310-65").
			Return(false, nil).
			Once()

		hasher_mock.On("HashPassword", "12345678").
			Return("abc123", nil).
			Once()

		db_mock.On("PersistUser", mock.AnythingOfType("entities.User")).
			Return(nil).
			Once()

		jwt_mock.On("CreateJwtToken", mock.AnythingOfType("entities.User")).
			Return("token", errors.New("error")).
			Once()

		req := events.APIGatewayProxyRequest{
			Body: `{"cpf":"218.486.310-65","pass":"12345678"}`,
		}

		// Act
		got, err := h.CrateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, got.StatusCode)

		db_mock.AssertExpectations(t)
		hasher_mock.AssertExpectations(t)
		jwt_mock.AssertExpectations(t)
	})
}
