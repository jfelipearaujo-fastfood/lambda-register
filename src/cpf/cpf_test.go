package cpf

import (
	"testing"
)

func TestNewCPF(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name           string
		args           args
		wantString     string
		wantValidation bool
	}{
		{
			name: "Create a new CPF with a valid value",
			args: args{
				s: "231.654.140-25",
			},
			wantString:     "231.654.140-25",
			wantValidation: true,
		},
		{
			name: "Create a new CPF with an invalid value",
			args: args{
				s: "762.000.810-11",
			},
			wantString:     "76200081011",
			wantValidation: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCPF(tt.args.s)

			if got.IsValid() != tt.wantValidation {
				t.Errorf("IsValid() = %v, want %v", got, tt.wantValidation)
				return
			}

			if got.String() != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}
		})
	}
}

func TestCPF_Mask(t *testing.T) {
	tests := []struct {
		name string
		c    CPF
		want string
	}{
		{
			name: "Mask a valid CPF",
			c:    NewCPF("231.654.140-25"),
			want: "231******25",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Mask(); got != tt.want {
				t.Errorf("CPF.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
