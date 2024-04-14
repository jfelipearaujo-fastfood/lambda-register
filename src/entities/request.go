package entities

type Request struct {
	CPF      string `json:"cpf"`
	Password string `json:"pass"`
}
