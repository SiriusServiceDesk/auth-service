package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/SiriusServiceDesk/auth-service/internal/helpers"
	"github.com/SiriusServiceDesk/auth-service/internal/models"
	"github.com/SiriusServiceDesk/auth-service/internal/repository"
	"github.com/SiriusServiceDesk/auth-service/internal/services"
	"github.com/SiriusServiceDesk/gateway-service/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Handler struct {
	auth_v1.UnimplementedAuthV1Server
	userService services.UserService
	redis       repository.RedisRepository
}

func (h Handler) Status(ctx context.Context, empty *emptypb.Empty) (*auth_v1.StatusResponse, error) {
	response := &auth_v1.StatusResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Time: %v, service working ok", time.Now().String()),
	}

	return response, nil
}

func (h Handler) Login(ctx context.Context, request *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	email := request.GetEmail()
	if email == "" {
		return LoginResponse(
			http.StatusBadRequest, "", "email is required", errors.New("email is required"),
		)
	}
	password := request.GetPassword()
	if password == "" {
		return LoginResponse(
			http.StatusBadRequest, "", "password is required", errors.New("password is required"),
		)
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		return LoginResponse(
			http.StatusBadRequest, "", "user not found", err,
		)
	}

	if !user.IsVerified {
		return LoginResponse(
			http.StatusUnauthorized, "", "user is not verified", errors.New("user is not verified"),
		)
	}

	if err := h.userService.ComparePassword(user.Password, password); err != nil {
		return LoginResponse(
			http.StatusBadRequest, "", "invalid password", err,
		)
	}

	token, err := h.userService.GenerateToken(user)
	if err != nil {
		return LoginResponse(
			http.StatusInternalServerError, "", "cant create token for user", err,
		)
	}

	return LoginResponse(http.StatusOK, token, "", nil)
}

func (h Handler) Registration(ctx context.Context, request *auth_v1.RegistrationRequest) (*auth_v1.RegistrationResponse, error) {
	email := request.GetEmail()
	if email == "" {
		return RegistrationResponse(
			http.StatusBadRequest, "email is required", errors.New("email is required"),
		)
	}
	password := request.GetPassword()
	if password == "" {
		return RegistrationResponse(
			http.StatusBadRequest, "password is required", errors.New("password is required"),
		)
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return RegistrationResponse(
			http.StatusInternalServerError, "failed to check if the user is in the database", err,
		)
	}

	if user != nil {
		if !user.IsVerified {
			return RegistrationResponse(
				http.StatusInternalServerError,
				"user exist but not verified", errors.New("user not verified"),
			)
		}
		return RegistrationResponse(
			http.StatusInternalServerError,
			"user is already registered", errors.New("user is already registered"),
		)
	}

	hashedPassword, err := h.userService.HashingPassword(password)
	if err != nil {
		return RegistrationResponse(
			http.StatusInternalServerError, "cant generate hashed password", err,
		)
	}
	newUser := &models.User{
		Password: string(hashedPassword),
		Email:    email,
	}

	submitCode := helpers.GenerateConfirmCode()

	if err := h.redis.Set(request.Email, submitCode); err != nil {
		return RegistrationResponse(
			http.StatusInternalServerError, "failed to generate confirm code", err,
		)
	}

	//TODO: доделать отправку сообщений на почту епта

	if _, err := h.userService.CreateUser(newUser); err != nil {
		return RegistrationResponse(
			http.StatusInternalServerError, "cant create user", err,
		)
	}

	return RegistrationResponse(http.StatusOK, "", nil)
}

func (h Handler) ConfirmEmail(ctx context.Context, request *auth_v1.ConfirmEmailRequest) (*auth_v1.ConfirmEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h Handler) User(ctx context.Context, request *auth_v1.UserRequest) (*auth_v1.UserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewHandler(userService services.UserService, redisRepository repository.RedisRepository) *Handler {
	return &Handler{
		userService: userService,
		redis:       redisRepository,
	}
}
