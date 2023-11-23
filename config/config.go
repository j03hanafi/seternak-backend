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

type Config struct {
	fiber    *fiber.Config
	fiberzap *fiberzap.Config
	recover  *recover.Config
}

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

	config.setFiberConfig()

	return config
}

func (c *Config) setDefaults() {
	// Set default App configuration
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("API_URL", "/api/v1")
}

func (c *Config) setFiberConfig() {
	c.fiber = &fiber.Config{
		CaseSensitive:            true,
		DisableHeaderNormalizing: true,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
	}
}

func (c *Config) GetFiberConfig() *fiber.Config {
	return c.fiber
}

func (c *Config) setFiberzapConfig() {
	c.fiberzap = &fiberzap.Config{
		Logger: logger.Get(),
		Fields: []string{"pid", "status", "method", "path", "queryParams", "ip", "ua", "latency", "time"},
	}
}

func (c *Config) GetFiberzapConfig() *fiberzap.Config {
	return c.fiberzap
}

func (c *Config) setRecoverConfig() {
	c.recover = &recover.Config{
		EnableStackTrace: true,
	}
}

func (c *Config) GetRecoverConfig() *recover.Config {
	return c.recover
}
