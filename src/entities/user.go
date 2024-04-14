package entities

import "github.com/google/uuid"

type User struct {
	Id          string `json:"id"`
	DocumentId  string `json:"document_id"`
	Password    string `json:"password"`
	IsAnonymous bool   `json:"is_anonymous"`
}

func NewAnonymousUser() User {
	return User{
		Id:          uuid.NewString(),
		IsAnonymous: true,
	}
}

func NewUser(documentId, password string) User {
	return User{
		Id:          uuid.NewString(),
		DocumentId:  documentId,
		Password:    password,
		IsAnonymous: false,
	}
}
