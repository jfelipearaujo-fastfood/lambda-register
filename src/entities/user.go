package entities

type User struct {
	Id          string `json:"id"`
	DocumentId  string `json:"document_id"`
	Password    string `json:"password"`
	IsAnonymous bool   `json:"is_anonymous"`
}
