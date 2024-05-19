package handlers

import (
	"github.com/SiriusServiceDesk/auth-service/internal/models"
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

func ConfirmEmailResponse(status int32, message string) (*auth_v1.ConfirmEmailResponse, error) {
	return &auth_v1.ConfirmEmailResponse{
		Status:  status,
		Message: message,
	}, nil
}

func ConfirmEmailErrorResponse(code codes.Code, message string) (*auth_v1.ConfirmEmailResponse, error) {
	return nil, status.Error(code, message)
}

func GetUserIdFromTokenErrorResponse(code codes.Code, message string) (*auth_v1.GetUserIdFromTokenResponse, error) {
	return nil, status.Error(code, message)
}

func GetUserByIdErrorResponse(code codes.Code, message string) (*auth_v1.GetUserByIdResponse, error) {
	return nil, status.Error(code, message)
}

func GetUserByIdResponse(user *models.User) (*auth_v1.GetUserByIdResponse, error) {
	return &auth_v1.GetUserByIdResponse{
		Name:       user.Name,
		Surname:    user.Surname,
		SecondName: "",
		Email:      user.Email,
		Role:       string(user.Role),
	}, nil
}

func GetUserIdFromTokenResponse(userId string) (*auth_v1.GetUserIdFromTokenResponse, error) {
	return &auth_v1.GetUserIdFromTokenResponse{
		UserId: userId,
	}, nil
}
