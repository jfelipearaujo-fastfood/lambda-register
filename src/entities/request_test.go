package entities

import (
	"testing"
)

func TestRequest_IsAnonymous(t *testing.T) {
	type fields struct {
		CPF      string
		Password string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Should return true when cpf and password are empty",
			fields: fields{
				CPF:      "",
				Password: "",
			},
			want: true,
		},
		{
			name: "Should return false when cpf is not empty",
			fields: fields{
				CPF:      "123",
				Password: "",
			},
			want: false,
		},
		{
			name: "Should return false when password is not empty",
			fields: fields{
				CPF:      "",
				Password: "123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Request{
				CPF:      tt.fields.CPF,
				Password: tt.fields.Password,
			}
			if got := r.IsAnonymous(); got != tt.want {
				t.Errorf("Request.IsAnonymous() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_IsPasswordWithMinimumLength(t *testing.T) {
	type fields struct {
		CPF      string
		Password string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Should return true when password is valid",
			fields: fields{
				CPF:      "",
				Password: "12345678",
			},
			want: true,
		},
		{
			name: "Should return false when password is empty",
			fields: fields{
				CPF:      "",
				Password: "",
			},
			want: false,
		},
		{
			name: "Should return false when password is less than minimum length",
			fields: fields{
				CPF:      "",
				Password: "123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Request{
				CPF:      tt.fields.CPF,
				Password: tt.fields.Password,
			}
			if got := r.IsPasswordWithMinimumLength(); got != tt.want {
				t.Errorf("Request.IsPasswordWithMinimumLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
