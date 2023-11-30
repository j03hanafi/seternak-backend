package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"go.uber.org/zap"
)

func Logger(l *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		ctx = logger.WithCtx(ctx, l)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
