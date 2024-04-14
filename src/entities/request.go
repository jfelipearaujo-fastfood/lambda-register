package entities

const (
	MINIMUM_PASSWORD_LENGTH = 8
)

type Request struct {
	CPF      string `json:"cpf"`
	Password string `json:"pass"`
}

func (r Request) IsAnonymous() bool {
	return r.CPF == "" && r.Password == ""
}

func (r Request) IsPasswordWithMinimumLength() bool {
	return len(r.Password) >= MINIMUM_PASSWORD_LENGTH
}
