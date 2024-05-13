package handlers

import (
	"context"
	"fmt"
	"github.com/SiriusServiceDesk/auth-service/internal/helpers"
	"github.com/SiriusServiceDesk/auth-service/internal/models"
	"github.com/SiriusServiceDesk/auth-service/internal/repository"
	"github.com/SiriusServiceDesk/auth-service/internal/services"
	"github.com/SiriusServiceDesk/gateway-service/pkg/auth_v1"
	"google.golang.org/grpc/codes"
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
		return LoginErrorResponse(codes.InvalidArgument, "email is required")
	}
	password := request.GetPassword()
	if password == "" {
		return LoginErrorResponse(codes.InvalidArgument, "password is required")
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		return LoginErrorResponse(codes.NotFound, "user is not found")
	}

	if !user.IsVerified {
		return LoginErrorResponse(codes.Aborted, "user is not verified")
	}

	if err := h.userService.ComparePassword(user.Password, password); err != nil {
		return LoginErrorResponse(codes.InvalidArgument, "invalid password")
	}

	token, err := h.userService.GenerateToken(user)
	if err != nil {
		return LoginErrorResponse(codes.Internal, "cant create token for user")
	}

	return LoginResponse(http.StatusOK, token, "user token generate successfully")
}

func (h Handler) Registration(ctx context.Context, request *auth_v1.RegistrationRequest) (*auth_v1.RegistrationResponse, error) {
	email := request.GetEmail()
	if email == "" {
		return RegistrationErrorResponse(codes.InvalidArgument, "email is required")
	}
	password := request.GetPassword()
	if password == "" {
		return RegistrationErrorResponse(codes.InvalidArgument, "email is required")
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return RegistrationErrorResponse(codes.NotFound, "failed to check if the user is in the database")
	}

	if user != nil {
		if !user.IsVerified {
			return RegistrationErrorResponse(codes.Internal, "user is not verified but created")
		}
		return RegistrationErrorResponse(codes.AlreadyExists, "user is already registered")
	}

	hashedPassword, err := h.userService.HashingPassword(password)
	if err != nil {
		return RegistrationErrorResponse(codes.Internal, "cant generate password hash")
	}
	newUser := &models.User{
		Password: string(hashedPassword),
		Email:    email,
	}

	submitCode := helpers.GenerateConfirmCode()

	if err := h.redis.Set(request.Email, submitCode); err != nil {
		return RegistrationErrorResponse(codes.Internal, "cant create confirm code")
	}

	//TODO: доделать отправку сообщений на почту епта

	if _, err := h.userService.CreateUser(newUser); err != nil {
		return RegistrationErrorResponse(codes.Internal, "cant create user")
	}

	return RegistrationResponse(http.StatusOK, "user created successfully")
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
