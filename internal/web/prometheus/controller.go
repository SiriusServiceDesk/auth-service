package prometheus

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Controller struct{}

func NewPrometheusController() *Controller {
	return &Controller{}
}

func (c *Controller) DefineRouter(app *fiber.App) {
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}
