package cpf

import (
	"testing"
)

func Test_sumDigit(t *testing.T) {
	type args struct {
		s     string
		table []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Sum only zeroes",
			args: args{
				s:     "000",
				table: []int{1, 2, 3},
			},
			want: 0,
		},
		{
			name: "Sum only ones",
			args: args{
				s:     "111",
				table: []int{1, 2, 3},
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sumDigit(tt.args.s, tt.args.table); got != tt.want {
				t.Errorf("sumDigit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClean(t *testing.T) {
	type args struct {
		cpf string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Clean all special characters",
			args: args{
				cpf: " 123.456.789-00",
			},
			want: "12345678900",
		},
		{
			name: "Clean nothing",
			args: args{
				cpf: "12345678900",
			},
			want: "12345678900",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clean(tt.args.cpf); got != tt.want {
				t.Errorf("Clean() = %v, want %v", got, tt.want)
			}
		})
	}
}
