package providers

import (
	"testing"
	"time"

	"github.com/jsfelipearaujo/lambda-register/src/providers/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestNewTimeProvider(t *testing.T) {
	type args struct {
		funcTime interfaces.FuncTime
	}
	tests := []struct {
		name string
		args args
		want *TimeProvider
	}{
		{
			name: "Should return a new instance correctly",
			args: args{
				funcTime: time.Now,
			},
			want: &TimeProvider{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange

			// Act
			timeProvider := NewTimeProvider(tt.args.funcTime)

			// Assert
			assert.IsType(t, tt.want, timeProvider)
		})
	}
}

func parseStringToTime(t *testing.T, input string) time.Time {
	out, err := time.Parse("2006-01-02 15:04:05", input)
	assert.NoError(t, err)
	return out
}

func TestTimeProvider_GetTime(t *testing.T) {
	type fields struct {
		funcTime interfaces.FuncTime
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "Should return correctly the time",
			fields: fields{
				funcTime: func() time.Time {
					return parseStringToTime(t, "2023-02-20 17:11:00")
				},
			},
			want: parseStringToTime(t, "2023-02-20 17:11:00"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			p := NewTimeProvider(tt.fields.funcTime)

			// Act
			got := p.GetTime()

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}
