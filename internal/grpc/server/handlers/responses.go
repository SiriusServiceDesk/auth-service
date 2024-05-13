package handlers

import "github.com/SiriusServiceDesk/gateway-service/pkg/auth_v1"

func LoginResponse(status int32, token, stringError string, err error) (*auth_v1.LoginResponse, error) {
	return &auth_v1.LoginResponse{
		Status: status,
		Token:  token,
		Error:  stringError,
	}, err
}

func RegistrationResponse(status int32, stringError string, err error) (*auth_v1.RegistrationResponse, error) {
	return &auth_v1.RegistrationResponse{
		Status: status,
		Error:  stringError,
	}, err
}
