package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"strings"
)

var methods = []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete, fiber.MethodOptions}
var headers = []string{fiber.HeaderAccept, fiber.HeaderAuthorization, fiber.HeaderContentType,
	fiber.HeaderContentLength, fiber.HeaderAcceptEncoding, "X-CSRF-Token"}

func SetupCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Join([]string{"*"}, ", "),
		AllowMethods:     strings.Join(methods, ", "),
		AllowHeaders:     strings.Join(headers, ", "),
		AllowCredentials: true,
		MaxAge:           300,
	})
}
