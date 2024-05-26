package initializers

import (
	"github.com/SiriusServiceDesk/auth-service/internal/app/dependencies"
	"github.com/SiriusServiceDesk/auth-service/internal/web"
	"github.com/SiriusServiceDesk/auth-service/internal/web/auth"
	"github.com/SiriusServiceDesk/auth-service/internal/web/prometheus"
	"github.com/SiriusServiceDesk/auth-service/internal/web/status"
	"github.com/SiriusServiceDesk/auth-service/internal/web/swagger"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container *dependencies.Container) {
	ctrl := buildRouters(container)

	for i := range ctrl {
		ctrl[i].DefineRouter(app)
	}
}

func buildRouters(container *dependencies.Container) []web.Controller {
	return []web.Controller{
		status.NewStatusController(),
		swagger.NewSwaggerController(),
		auth.NewAuthController(container.UserService, container.Redis),
		prometheus.NewPrometheusController(),
	}
}
