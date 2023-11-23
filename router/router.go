package router

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	configuration "github.com/j03hanafi/seternak-backend/config"
	"github.com/j03hanafi/seternak-backend/controller"
	"github.com/j03hanafi/seternak-backend/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// New initializes and returns a new Fiber application with configured middleware.
// Returns a pointer to the fiber.App instance.
func New() *fiber.App {
	l := logger.Get()
	defer func(l *zap.Logger) {
		_ = l.Sync()
	}(l)

	config := configuration.New()

	app := fiber.New(*config.GetFiberConfig())
	app.Use(fiberzap.New(*config.GetFiberzapConfig()))
	app.Use(recover.New(*config.GetRecoverConfig()))

	newHandler(&handlerConfig{
		app:     app,
		baseURL: viper.GetString("API_URL"),
		version: controller.NewVersion(),
	})

	return app
}
