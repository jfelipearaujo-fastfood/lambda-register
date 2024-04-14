package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewAnonymousUser(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Should return a new anonymous user",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAnonymousUser()

			if got.IsAnonymous != tt.want {
				t.Errorf("NewAnonymousUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	type args struct {
		documentId string
		password   string
	}
	tests := []struct {
		name string
		args args
		want User
	}{
		{
			name: "Should return a new user",
			args: args{
				documentId: "123",
				password:   "123456",
			},
			want: User{
				DocumentId:  "123",
				Password:    "123456",
				IsAnonymous: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUser(tt.args.documentId, tt.args.password)

			if err := uuid.Validate(got.Id); err != nil {
				t.Errorf("NewUser().Id = %v, want a valid UUIDm got %v", got.Id, err)
			}

			if got.DocumentId != tt.want.DocumentId {
				t.Errorf("NewUser().DocumentId = %v, want %v", got.DocumentId, tt.want.DocumentId)
			}

			if got.Password != tt.want.Password {
				t.Errorf("NewUser().Password = %v, want %v", got.Password, tt.want.Password)
			}

			if got.IsAnonymous != tt.want.IsAnonymous {
				t.Errorf("NewUser().IsAnonymous = %v, want %v", got.IsAnonymous, tt.want.IsAnonymous)
			}
		})
	}
}
