package router

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestInvalidRequestBody(t *testing.T) {
	tests := []struct {
		name string
		want events.APIGatewayProxyResponse
	}{
		{
			name: "InvalidRequestBody",
			want: events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       `{"status":400,"message":"error to parse the request body"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InvalidRequestBody(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvalidRequestBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidCPFOrPassword(t *testing.T) {
	tests := []struct {
		name string
		want events.APIGatewayProxyResponse
	}{
		{
			name: "InvalidCPFOrPassword",
			want: events.APIGatewayProxyResponse{
				StatusCode: 401,
				Body:       `{"status":401,"message":"invalid cpf or password"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InvalidCPFOrPassword(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvalidCPFOrPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInternalServerError(t *testing.T) {
	tests := []struct {
		name string
		want events.APIGatewayProxyResponse
	}{
		{
			name: "InternalServerError",
			want: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"status":500,"message":"internal server error"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InternalServerError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InternalServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	tests := []struct {
		name string
		want events.APIGatewayProxyResponse
	}{
		{
			name: "MethodNotAllowed",
			want: events.APIGatewayProxyResponse{
				StatusCode: 405,
				Body:       `{"status":405,"message":"method not allowed"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MethodNotAllowed(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MethodNotAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want events.APIGatewayProxyResponse
	}{
		{
			name: "Success",
			args: args{
				token: "token",
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"status":200,"message":"success","access_token":"token"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Success(tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Success() = %v, want %v", got, tt.want)
			}
		})
	}
}
