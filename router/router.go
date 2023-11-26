package router

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	configuration "github.com/j03hanafi/seternak-backend/config"
	"github.com/j03hanafi/seternak-backend/handler"
	"github.com/j03hanafi/seternak-backend/repository"
	"github.com/j03hanafi/seternak-backend/service"
	"github.com/j03hanafi/seternak-backend/utils/logger"
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

	/*
		Repository initialization
	*/
	userRepository := repository.NewPGUserRepository(config.GetDB())

	/*
		Service initialization
	*/
	userService := service.NewUserService(&service.UserServiceConfig{
		UserRepository: userRepository,
	})

	// Fiber instance
	app := fiber.New(*config.GetFiberConfig())
	app.Use(fiberzap.New(*config.GetFiberzapConfig()))
	app.Use(recover.New(*config.GetRecoverConfig()))

	// API initialization
	newAPI(&apiConfig{
		app:     app,
		baseURL: viper.GetString("API_URL"),
		version: handler.NewVersion(),
		user: handler.NewUser(&handler.UserHandlerConfig{
			UserService: userService,
		}),
	})

	return app
}
