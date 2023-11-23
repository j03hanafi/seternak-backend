package config

import (
	"errors"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/j03hanafi/seternak-backend/logger"
	"github.com/spf13/viper"
	"log"
)

// Config defines configuration settings for the server.
type Config struct {
	fiber    *fiber.Config
	fiberzap *fiberzap.Config
	recover  *recover.Config
}

// New initializes a new Config struct, sets default values, and loads environment variables.
// It returns a pointer to the Config struct, or logs a fatal error if config loading fails.
func New() *Config {
	config := &Config{}

	config.setDefaults()

	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../.") // For testing _test.go files
	viper.AllowEmptyEnv(false)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			log.Fatalf("Error reading config file, %v", err)
		}
	}

	// Set Config struct field values
	config.setFiberConfig()
	config.setFiberzapConfig()
	config.setRecoverConfig()

	return config
}

// setDefaults sets default values for application configuration.
func (c *Config) setDefaults() {
	// Set default App configuration
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("API_URL", "/api/v1")
}

// setFiberConfig initializes Fiber's configuration with custom settings.
// Updates the Config struct fiber field with the new settings.
func (c *Config) setFiberConfig() {
	c.fiber = &fiber.Config{
		CaseSensitive:            true,
		DisableHeaderNormalizing: true,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
		ErrorHandler:             c.fiberErrorHandler,
	}
}

// GetFiberConfig retrieves the Fiber framework's configuration from the Config struct.
// Returns a pointer to the fiber.Config instance.
func (c *Config) GetFiberConfig() *fiber.Config {
	return c.fiber
}

// setFiberzapConfig initializes and configures the fiberzap logger settings.
// Updates the Config struct fiberzap field with the new settings.
func (c *Config) setFiberzapConfig() {
	c.fiberzap = &fiberzap.Config{
		Logger: logger.Get(),
		Fields: []string{"pid", "status", "method", "path", "queryParams", "ip", "ua", "latency", "time", "resBody", "error"},
	}
}

// GetFiberzapConfig retrieves the fiberzap logging configuration from the Config struct.
// Returns a pointer to the fiberzap.Config instance.
func (c *Config) GetFiberzapConfig() *fiberzap.Config {
	return c.fiberzap
}

// setRecoverConfig initializes and sets the configuration for the recovery middleware.
// Updates the Config struct recover field with stack trace enabled.
func (c *Config) setRecoverConfig() {
	c.recover = &recover.Config{
		EnableStackTrace: true,
	}
}

// GetRecoverConfig retrieves the recover middleware configuration from the Config struct.
// Returns a pointer to the recover.Config instance.
func (c *Config) GetRecoverConfig() *recover.Config {
	return c.recover
}

func (c *Config) fiberErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Set Content-Type: text/plain; charset=utf-8
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	// Return status code with error message
	return ctx.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
