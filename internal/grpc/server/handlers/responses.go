package handlers

import (
	"github.com/SiriusServiceDesk/gateway-service/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoginResponse(status int32, token, message string) (*auth_v1.LoginResponse, error) {
	return &auth_v1.LoginResponse{
		Status:  status,
		Token:   token,
		Message: message,
	}, nil
}

func LoginErrorResponse(code codes.Code, message string) (*auth_v1.LoginResponse, error) {
	return nil, status.Error(code, message)
}

func RegistrationResponse(status int32, message string) (*auth_v1.RegistrationResponse, error) {
	return &auth_v1.RegistrationResponse{
		Status:  status,
		Message: message,
	}, nil
}

func RegistrationErrorResponse(code codes.Code, message string) (*auth_v1.RegistrationResponse, error) {
	return nil, status.Error(code, message)
}
