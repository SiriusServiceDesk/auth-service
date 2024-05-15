package swagger

import (
	_ "github.com/SiriusServiceDesk/auth-service/api"
	"github.com/SiriusServiceDesk/auth-service/internal/middleware"
	"github.com/SiriusServiceDesk/auth-service/internal/web"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

var _ web.Controller = (*Controller)(nil)

type Controller struct{}

func NewSwaggerController() *Controller {
	return &Controller{}
}

func (c *Controller) DefineRouter(app *fiber.App) {
	app.Use(middleware.SetupCORS())
	app.Get("/swagger/*", swagger.HandlerDefault)
}
