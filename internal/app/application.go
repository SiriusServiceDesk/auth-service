package app

import (
	"github.com/SiriusServiceDesk/auth-service/internal/app/dependencies"
	"github.com/SiriusServiceDesk/auth-service/internal/app/initializers"
	"github.com/SiriusServiceDesk/auth-service/internal/config"
	"github.com/SiriusServiceDesk/auth-service/internal/repository"
	"github.com/SiriusServiceDesk/auth-service/internal/services"
	"github.com/SiriusServiceDesk/auth-service/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Application struct{}

func InitApplication(app *fiber.App) {
	cfg := config.GetConfig()

	userRepository := repository.NewUserRepository()
	userService := services.NewUserService(userRepository)

	redisRepository := repository.NewRedisClient(cfg)

	container := &dependencies.Container{
		UserService: userService,
		Redis:       redisRepository,
	}

	grpcListener, err := initializers.InitializeGRPCListener()
	if err != nil {
		logger.Fatal("failed initializing grpc listener", zap.Error(err))
	}

	grpcServer := initializers.InitializeGRPCServer(grpcListener, container)

	initializers.SetupRoutes(app, container)

	initializers.StartGRPCServer(grpcServer)
}
