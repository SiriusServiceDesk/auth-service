package auth

import (
	"encoding/json"
	"github.com/SiriusServiceDesk/auth-service/internal/grpc/client"
	"github.com/SiriusServiceDesk/auth-service/internal/helpers"
	"github.com/SiriusServiceDesk/auth-service/internal/middleware"
	"github.com/SiriusServiceDesk/auth-service/internal/models"
	"github.com/SiriusServiceDesk/auth-service/internal/repository"
	"github.com/SiriusServiceDesk/auth-service/internal/services"
	"github.com/SiriusServiceDesk/auth-service/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

type Controller struct {
	userService services.UserService
	redis       repository.RedisRepository
}

func NewAuthController(userService services.UserService, redis repository.RedisRepository) *Controller {
	return &Controller{userService: userService, redis: redis}
}

func (ctrl *Controller) DefineRouter(app *fiber.App) {
	authGroup := app.Group("/v1/auth")

	authGroup.Use(middleware.BenchmarkMiddleware())
	authGroup.Use(middleware.SetupCORS())

	authGroup.Post("/login", ctrl.login)
	authGroup.Post("/registration", ctrl.registration)
	authGroup.Post("/confirmEmail", ctrl.confirmEmail)
	authGroup.Post("/resendCode", ctrl.resendCode)
	authGroup.Post("/resetPassword", ctrl.resetPassword)
	authGroup.Post("/resetPassword/confirm", ctrl.resetPasswordConfirm)

	userGroup := app.Group("/v1/user")

	userGroup.Use(middleware.BenchmarkMiddleware())
	userGroup.Use(middleware.SetupCORS())

	userGroup.Get("/user", ctrl.user)
	userGroup.Get("/user/:id", ctrl.userById)
	userGroup.Put("/user", ctrl.updateUser)
	userGroup.Post("/user", ctrl.createUser)
	userGroup.Delete("/user/:id", ctrl.deleteUser)

}

func (ctrl *Controller) login(ctx *fiber.Ctx) error {
	var request LoginRequest
	if err := ctx.BodyParser(&request); err != nil {
		logger.Debug("login.BodyParser", zap.Error(err))
		return Response().WithDetails(err).BadRequest(ctx, "cant parse data")
	}

	user, err := ctrl.userService.GetUserByEmail(request.Email)
	if err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant get user from database")
	}

	if !user.IsVerified {
		return Response().BadRequest(ctx, "user is not verified")
	}

	if err := ctrl.userService.ComparePassword(user.Password, request.Password); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "incorrect email or password")
	}

	token, err := ctrl.userService.GenerateToken(user)
	if err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant create user token")
	}

	return Response().StatusOK(ctx, LoginResponse{Token: token})
}

func (ctrl *Controller) registration(ctx *fiber.Ctx) error {
	var request RegistrationRequest

	if err := ctx.BodyParser(&request); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant parse data")
	}

	user, err := ctrl.userService.GetUserByEmail(request.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant get user from database")
	}

	if user != nil {
		if !user.IsVerified {
			return Response().BadRequest(ctx, "user is not verified")
		}
		return Response().BadRequest(ctx, "user existing")
	}

	hashedPassword, err := ctrl.userService.HashingPassword(request.Password)
	if err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant hash password")
	}

	newUser := &models.User{
		Password: string(hashedPassword),
		Email:    request.Email,
		Role:     models.UserR,
	}

	submitCode := helpers.GenerateConfirmCode()

	if err := ctrl.redis.Set(request.Email, submitCode); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant generate submit code")
	}

	if _, err := ctrl.userService.CreateUser(newUser); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant create user")
	}

	jsonData, err := json.Marshal(VerificationMessageData{Code: submitCode})
	if err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant marshal json data")
	}

	message := client.Message{
		To:           []string{request.Email},
		Data:         string(jsonData),
		Type:         "email",
		Subject:      "Verifying email for registration",
		TemplateName: "verifyingEmail",
	}

	if _, err := client.SendMessage(&message); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant send message")
	}

	return Response().StatusOK(ctx, "user created successfully")
}

func (ctrl *Controller) confirmEmail(ctx *fiber.Ctx) error {
	var request ConfirmEmailRequest

	if err := ctx.BodyParser(&request); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant parse data")
	}

	savedCode, err := ctrl.redis.Get(request.Email)
	if err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "can't get saved code from cache")
	}

	savedCodeInt, err := strconv.Atoi(savedCode)
	if err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "can't convert code type string to int")
	}

	if request.VerificationCode != savedCodeInt {
		return Response().BadRequest(ctx, "the saved code doesn't match the code that came in")
	}

	user, err := ctrl.userService.GetUserByEmail(request.Email)
	if err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "can't get user from database")
	}

	user.IsVerified = true
	if _, err := ctrl.userService.UpdateUser(user.Id, user); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "can't update user")
	}

	if err := ctrl.redis.Delete(request.Email); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "can't delete from cache")
	}

	return Response().StatusOK(ctx, "email confirmation successful")
}

func (ctrl *Controller) resendCode(ctx *fiber.Ctx) error {
	var request ResendCodeRequest
	if err := ctx.BodyParser(&request); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant parse data")
	}

	submitCode := helpers.GenerateConfirmCode()

	if err := ctrl.redis.Set(request.Email, submitCode); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant generate submit code")
	}

	jsonData, err := json.Marshal(VerificationMessageData{Code: submitCode})
	if err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant marshal json data")
	}

	message := client.Message{
		To:           []string{request.Email},
		Data:         string(jsonData),
		Type:         "email",
		Subject:      "Resend new code",
		TemplateName: "verifyingEmail",
	}

	if _, err := client.SendMessage(&message); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant send message")
	}

	return Response().StatusOK(ctx, "new code sent successfully")
}

func (ctrl *Controller) resetPassword(ctx *fiber.Ctx) error {
	var request ResetPasswordRequest
	if err := ctx.BodyParser(&request); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant parse request data")
	}

	message := client.Message{
		To:           []string{request.Email},
		Data:         "",
		Type:         "email",
		Subject:      "Reset password",
		TemplateName: "resetPassword",
	}

	if _, err := client.SendMessage(&message); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant send message to ")
	}

	return Response().StatusOK(ctx, "reset password message sent successfully")
}

func (ctrl *Controller) resetPasswordConfirm(ctx *fiber.Ctx) error {
	var request ResetPasswordConfirmRequest
	if err := ctx.BodyParser(&request); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant parse request data")
	}

	user, err := ctrl.userService.GetUserByEmail(request.Email)
	if err != nil {
		return Response().WithDetails(err).BadRequest(ctx, "user not found")
	}

	if err := ctrl.userService.ComparePassword(user.Password, request.NewPassword); err == nil {
		return Response().WithDetails(err).BadRequest(ctx, "passwords are same")
	}

	newPasswordHash, err := ctrl.userService.HashingPassword(request.NewPassword)
	if err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant hash password")
	}

	if err := ctrl.userService.UpdatePassword(user.Id, string(newPasswordHash)); err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "failed to update user")
	}

	return Response().StatusOK(ctx, "user updated successfully")
}

func (ctrl *Controller) user(ctx *fiber.Ctx) error {
	userId, err := helpers.GetUserIdFromToken(ctx)
	if err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "cant get user id from header")
	}

	user, err := ctrl.userService.GetUserById(userId)
	if err != nil {
		return Response().WithDetails(err).ServerInternalError(ctx, "failed to get user")
	}

	return Response().StatusOK(ctx, MappingUserModelToResponse(user))
}

func (ctrl *Controller) userById(ctx *fiber.Ctx) error {
	panic("implement me pls")
}

func (ctrl *Controller) updateUser(ctx *fiber.Ctx) error {
	panic("implement me pls")
}

func (ctrl *Controller) createUser(ctx *fiber.Ctx) error {
	panic("implement me pls")
}

func (ctrl *Controller) deleteUser(ctx *fiber.Ctx) error {
	panic("implement me pls")
}
