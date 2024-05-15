package auth

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type RawResponse struct {
	Status  int         `json:"status" example:"200"`
	Payload interface{} `json:"payload"`
	Details interface{} `json:"details,omitempty"`
}

func Response() *RawResponse {
	return &RawResponse{}
}

func (r *RawResponse) BadRequest(c *fiber.Ctx, payload interface{}) error {
	r.Status = http.StatusBadRequest
	r.Payload = payload
	return r.send(c)
}

func (r *RawResponse) StatusOK(c *fiber.Ctx, payload interface{}) error {
	r.Status = http.StatusOK
	r.Payload = payload
	return r.send(c)
}

func (r *RawResponse) ServerInternalError(c *fiber.Ctx, payload interface{}) error {
	r.Status = http.StatusInternalServerError
	r.Payload = payload
	return r.send(c)
}

func (r *RawResponse) WithDetails(details error) *RawResponse {
	r.Details = details
	return r
}

func (r *RawResponse) send(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	return c.Status(r.Status).JSON(r)
}
