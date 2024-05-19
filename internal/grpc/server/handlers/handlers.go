package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SiriusServiceDesk/auth-service/internal/grpc/client"
	"github.com/SiriusServiceDesk/auth-service/internal/helpers"
	"github.com/SiriusServiceDesk/auth-service/internal/models"
	"github.com/SiriusServiceDesk/auth-service/internal/repository"
	"github.com/SiriusServiceDesk/auth-service/internal/services"
	"github.com/SiriusServiceDesk/auth-service/pkg/logger"
	"github.com/SiriusServiceDesk/gateway-service/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type VerificationMessageData struct {
	Code int `json:"Code"`
}

type Handler struct {
	auth_v1.UnimplementedAuthV1Server
	userService services.UserService
	redis       repository.RedisRepository
}

func (h Handler) GetUserById(ctx context.Context, request *auth_v1.GetUserByIdRequest) (*auth_v1.GetUserByIdResponse, error) {
	userId := request.GetUserId()
	if userId == "" {
		logger.Debug("userId is empty")
		return GetUserByIdErrorResponse(codes.InvalidArgument, "userId is required")
	}

	user, err := h.userService.GetUserById(userId)
	if err != nil {
		logger.Debug("cant get user from db")
		return GetUserByIdErrorResponse(codes.Internal, "failed to get user")
	}

	return GetUserByIdResponse(user)
}

func (h Handler) GetUserIdFromToken(ctx context.Context, request *auth_v1.GetUserIdFromTokenRequest) (*auth_v1.GetUserIdFromTokenResponse, error) {
	header := request.GetHeader()
	if len(header) == 0 {
		return GetUserIdFromTokenErrorResponse(codes.Internal, "token is required")
	}

	token, err := helpers.GetTokenFromHeaders(header)
	if err != nil {
		return GetUserIdFromTokenErrorResponse(codes.Internal, err.Error())
	}

	userId, err := helpers.GetUserIdFromToken(token)
	if err != nil {
		return GetUserIdFromTokenErrorResponse(codes.Internal, "failed to get userId")
	}

	return GetUserIdFromTokenResponse(userId)
}

func (h Handler) User(ctx context.Context, request *auth_v1.UserRequest) (*auth_v1.UserResponse, error) {
	//TODO implement me
	panic("implement me")
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
		return LoginErrorResponse(codes.Internal, "user is not verified")
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

	if _, err := h.userService.CreateUser(newUser); err != nil {
		return RegistrationErrorResponse(codes.Internal, "cant create user")
	}

	jsonData, err := json.Marshal(VerificationMessageData{Code: submitCode})
	if err != nil {
		return RegistrationErrorResponse(codes.Internal, "cant marshal json data")
	}

	message := client.Message{
		To:           []string{request.GetEmail()},
		Data:         string(jsonData),
		Type:         "email",
		Subject:      "Verifying email for registration",
		TemplateName: "verifyingEmail",
	}

	if _, err := client.SendMessage(&message); err != nil {
		return RegistrationErrorResponse(codes.Internal, "cant send message")
	}

	return RegistrationResponse(http.StatusOK, "user created successfully")
}

func (h Handler) ConfirmEmail(ctx context.Context, request *auth_v1.ConfirmEmailRequest) (*auth_v1.ConfirmEmailResponse, error) {
	email := request.GetEmail()
	verificationCode := request.GetVerificationCode()

	if email == "" {
		return ConfirmEmailErrorResponse(codes.InvalidArgument, "email is required")
	}
	if verificationCode == 0 {
		return ConfirmEmailErrorResponse(codes.InvalidArgument, "verificationCode is required")
	}

	savedCode, err := h.redis.Get(email)
	if err != nil {
		return ConfirmEmailErrorResponse(codes.Internal, "can't get saved code from cache")
	}

	savedCodeInt, err := strconv.Atoi(savedCode)
	if err != nil {
		return ConfirmEmailErrorResponse(codes.Internal, "can't convert code type string to int")
	}

	if verificationCode != int32(savedCodeInt) {
		return ConfirmEmailErrorResponse(codes.Canceled, "the saved code doesn't match the code that came in")
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		return ConfirmEmailErrorResponse(codes.Internal, "can't get user from database")
	}

	user.IsVerified = true
	if _, err := h.userService.UpdateUser(user.Id, user); err != nil {
		return ConfirmEmailErrorResponse(codes.Internal, "can't update user")
	}

	if err := h.redis.Delete(email); err != nil {
		return ConfirmEmailErrorResponse(codes.Internal, "can't delete from cache")
	}

	return ConfirmEmailResponse(http.StatusOK, "email confirmation successful")

}

func NewHandler(userService services.UserService, redisRepository repository.RedisRepository) *Handler {
	return &Handler{
		userService: userService,
		redis:       redisRepository,
	}
}
