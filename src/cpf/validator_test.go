package cpf

import "testing"

func TestValidateCPF(t *testing.T) {
	type args struct {
		cpf string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Validate a valid CPF",
			args: args{
				cpf: "16115279020",
			},
			want: true,
		},
		{
			name: "Validate an invalid CPF",
			args: args{
				cpf: "16115273020",
			},
			want: false,
		},
		{
			name: "Validate an invalid CPF with wrong length",
			args: args{
				cpf: "123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateCPF(tt.args.cpf); got != tt.want {
				t.Errorf("ValidateCPF() = %v, want %v", got, tt.want)
			}
		})
	}
}
