package server

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	configuration "github.com/j03hanafi/seternak-backend/config"
	"github.com/j03hanafi/seternak-backend/handler"
	"github.com/j03hanafi/seternak-backend/repository"
	"github.com/j03hanafi/seternak-backend/service"
	"github.com/spf13/viper"
)

// New initializes and returns a new Fiber application with configured middleware.
// Returns a pointer to the fiber.App instance.
func New() (*fiber.App, func() error) {
	config := configuration.New()

	/*
		Repository initialization
	*/
	userRepository := repository.NewPGUser(config.GetDB())
	authRepository := repository.NewRedisAuth(config.GetRedis())

	/*
		Service initialization
	*/
	userService := service.NewUser(&service.UserServiceConfig{
		UserRepository: userRepository,
		AuthRepository: authRepository,
	})
	authService := service.NewAuth(&service.AuthServiceConfig{
		AuthRepository:             authRepository,
		PrivateKey:                 config.GetPrivateKey(),
		RefreshTokenSecret:         viper.GetString("REFRESH_TOKEN_SECRET"),
		IDTokenExpirationSecs:      viper.GetInt64("ID_TOKEN_EXP"),
		RefreshTokenExpirationSecs: viper.GetInt64("REFRESH_TOKEN_EXP"),
	})

	// Fiber instance
	app := fiber.New(*config.GetFiberConfig())
	app.Use(fiberzap.New(*config.GetFiberzapConfig()))
	app.Use(recover.New(*config.GetRecoverConfig()))

	/*
		API initialization
	*/
	newAPI(&apiConfig{
		app:     app,
		baseURL: viper.GetString("API_URL"),
		version: handler.NewVersion(),
		user: handler.NewUser(&handler.UserHandlerConfig{
			UserService: userService,
			AuthService: authService,
		}),
		publicKey: config.GetPublicKey(),
		zapLogger: config.GetLogger(),
		secretKey: viper.GetString("REFRESH_TOKEN_SECRET"),
	})

	return app, config.Close
}
