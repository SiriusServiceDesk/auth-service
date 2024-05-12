package app

import (
	"github.com/SiriusServiceDesk/auth-service/internal/app/dependencies"
	"github.com/SiriusServiceDesk/auth-service/internal/app/initializers"
	"github.com/SiriusServiceDesk/auth-service/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type Application struct{}

func InitApplication(app *fiber.App) {
	repository.NewExampleRepository()

	container := &dependencies.Container{}

	initializers.SetupRoutes(app, container)
}
