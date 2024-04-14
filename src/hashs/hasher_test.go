package hashs

import "testing"

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		notWant string
		wantErr bool
	}{
		{
			name: "Hash a password",
			args: args{
				password: "123456",
			},
			notWant: "123456",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.notWant {
				t.Errorf("HashPassword() = %v, not want %v", got, tt.notWant)
			}
		})
	}
}
