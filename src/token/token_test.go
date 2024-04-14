package token

import (
	"testing"

	"github.com/jsfelipearaujo/lambda-register/src/entities"
)

func TestCreateJwtToken(t *testing.T) {
	type args struct {
		user entities.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create a JWT Token",
			args: args{
				user: entities.User{
					Id:          "1",
					DocumentId:  "12345678900",
					Password:    "123456",
					IsAnonymous: false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateJwtToken(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateJwtToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("CreateJwtToken() = %v, want not empty", got)
			}
		})
	}
}
