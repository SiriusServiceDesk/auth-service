package status

import (
	"github.com/SiriusServiceDesk/auth-service/internal/web"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

var (
	_ web.Controller = (*Controller)(nil)
)

type Controller struct{}

func NewStatusController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) DefineRouter(app *fiber.App) {
	app.Get("/api/v1/status/", ctrl.Status)
}

// Status
// @Summary Get the status
// @Description Get the status of the API
// @ID Status
// @Produce json
// @Success 200 {object} Response
// @Router /api/v1/status/ [get]
func (ctrl *Controller) Status(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(Response{
		Code:    http.StatusOK,
		Message: "Success",
	})
}
