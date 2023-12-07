package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/j03hanafi/seternak-backend/utils/id"
)

func RequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Generator: func() string {
			rid := id.New()
			return rid.String()
		},
	})
}
