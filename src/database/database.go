package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jsfelipearaujo/lambda-register/src/entities"
	"github.com/jsfelipearaujo/lambda-register/src/providers/interfaces"
	_ "github.com/lib/pq"
)

const (
	engine = "postgres"

	DOCUMENT_TYPE_CPF = 1
)

var (
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")
)

type Database struct {
	conn         *sql.DB
	timeProvider interfaces.TimeProvider
}

func NewDatabase(db *sql.DB, timeProvider interfaces.TimeProvider) *Database {
	return &Database{
		conn:         db,
		timeProvider: timeProvider,
	}
}

func NewDatabaseFromConnStr(timeProvider interfaces.TimeProvider) *Database {
	db, err := sql.Open(engine, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName))
	if err != nil {
		panic(err)
	}

	return &Database{
		conn:         db,
		timeProvider: timeProvider,
	}
}

func (db *Database) CheckIfCPFIsInUse(cpf string) (bool, error) {
	statement, err := db.conn.Query("SELECT COUNT(c.id) As count FROM customers c WHERE c.document_id = $1;", cpf)
	if err != nil {
		return false, err
	}

	var count int
	for statement.Next() {
		if err := statement.Scan(&count); err != nil {
			return false, err
		}
	}

	return count > 0, nil
}

func (db *Database) PersistUser(user entities.User) error {
	if user.IsAnonymous {
		_, err := db.conn.Exec("INSERT INTO customers (id, document_type, is_anonymous, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);",
			user.Id,
			DOCUMENT_TYPE_CPF,
			true,
			db.timeProvider.GetTime(),
			db.timeProvider.GetTime())

		if err != nil {
			return err
		}
	} else {
		_, err := db.conn.Exec("INSERT INTO customers (id, document_id, document_type, is_anonymous, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7);",
			user.Id,
			user.DocumentId,
			DOCUMENT_TYPE_CPF,
			false,
			user.Password,
			db.timeProvider.GetTime(),
			db.timeProvider.GetTime())

		if err != nil {
			return err
		}
	}

	return nil
}
