package database

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jfelipearaujo-org/lambda-register/internal/entities"
	"github.com/jfelipearaujo-org/lambda-register/internal/providers/interfaces/mocks"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestDatabase_CheckIfCPFIsInUse_IsInUse(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeProviderMock := mocks.NewMockTimeProvider(t)

	database := NewDatabase(db, timeProviderMock)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery("SELECT (.+) FROM customers c").
		WithArgs("123").
		WillReturnRows(rows)

	expectedResult := true

	// Act
	result, err := database.CheckIfCPFIsInUse("123")
	if err != nil {
		t.Errorf("error while checking if CPF is in use: %v", err)
	}

	// Assert
	if result != expectedResult {
		t.Errorf("expected result to be %v, got %v", expectedResult, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDatabase_CheckIfCPFIsInUse_IsNotInUse(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeProviderMock := mocks.NewMockTimeProvider(t)

	database := NewDatabase(db, timeProviderMock)

	mock.ExpectQuery("SELECT (.+) FROM customers c").
		WithArgs("123").
		WillReturnRows(sqlmock.NewRows([]string{}))

	expectedResult := false

	// Act
	result, err := database.CheckIfCPFIsInUse("123")
	if err != nil {
		t.Errorf("error while checking if CPF is in use: %v", err)
	}

	// Assert
	if result != expectedResult {
		t.Errorf("expected result to be %v, got %v", expectedResult, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func parseStringToTime(t *testing.T, input string) time.Time {
	out, err := time.Parse("2006-01-02 15:04:05", input)
	assert.NoError(t, err)
	return out
}

func TestDatabase_PersistUser_NonAnonymous(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeProviderMock := mocks.NewMockTimeProvider(t)

	now := parseStringToTime(t, "2024-04-13 23:37:11")

	timeProviderMock.On("GetTime").
		Return(now).
		Times(2)

	database := NewDatabase(db, timeProviderMock)

	mock.ExpectExec("INSERT INTO customers").
		WithArgs("1", "123", DOCUMENT_TYPE_CPF, false, "123456", now, now).
		WillReturnResult(sqlmock.NewResult(1, 1))

	user := entities.User{
		Id:          "1",
		DocumentId:  "123",
		Password:    "123456",
		IsAnonymous: false,
	}

	// Act
	err = database.PersistUser(user)

	// Assert
	if err != nil {
		t.Errorf("error while persisting the user: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDatabase_PersistUser_Anonymous(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeProviderMock := mocks.NewMockTimeProvider(t)

	now := parseStringToTime(t, "2024-04-13 23:37:11")

	timeProviderMock.On("GetTime").
		Return(now).
		Times(2)

	database := NewDatabase(db, timeProviderMock)

	mock.ExpectExec("INSERT INTO customers").
		WithArgs("1", DOCUMENT_TYPE_CPF, true, now, now).
		WillReturnResult(sqlmock.NewResult(1, 1))

	user := entities.User{
		Id:          "1",
		DocumentId:  "123",
		Password:    "123456",
		IsAnonymous: true,
	}

	// Act
	err = database.PersistUser(user)

	// Assert
	if err != nil {
		t.Errorf("error while persisting the user: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewDatabase(t *testing.T) {
	// Arrange
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timeProviderMock := mocks.NewMockTimeProvider(t)

	// Act
	database := NewDatabase(db, timeProviderMock)

	// Assert
	assert.NotNil(t, database)
}

func TestNewDatabaseFromConnStr(t *testing.T) {
	// Arrange
	timeProviderMock := mocks.NewMockTimeProvider(t)

	// Act
	database := NewDatabaseFromConnStr(timeProviderMock)

	// Assert
	assert.NotNil(t, database)
}
