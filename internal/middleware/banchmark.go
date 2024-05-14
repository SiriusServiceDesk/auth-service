package middleware

import (
	"github.com/SiriusServiceDesk/auth-service/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

func BenchmarkMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		logger.Info("Request completed successfully",
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
			zap.Any("duration", duration),
		)
		return err
	}
}
