package config

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/utils/apperrors"
	"github.com/j03hanafi/seternak-backend/utils/consts"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"moul.io/zapgorm2"
)

// Config defines configuration settings for the server.
type Config struct {
	db       *gorm.DB
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
	viper.AddConfigPath("../.")    // For testing _test.go files
	viper.AddConfigPath("../../.") // For testing _test.go files in application directory
	viper.AllowEmptyEnv(false)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			log.Fatalf("Error reading config file, %v", err)
		}
	}

	// Set Config struct field values
	config.setDB()
	config.setFiberConfig()
	config.setFiberzapConfig()
	config.setRecoverConfig()

	return config
}

// setDefaults sets default values for application configuration.
func (c *Config) setDefaults() {
	// Set default App configuration
	viper.SetDefault("APP_ENV", consts.DevelopmentMode)
	viper.SetDefault("API_URL", "/api/v1")

	// Set default DB configuration
	viper.SetDefault("PG_HOST", "localhost")
	viper.SetDefault("PG_PORT", "5432")
	viper.SetDefault("PG_USER", "postgres")
	viper.SetDefault("PG_PASS", "password")
	viper.SetDefault("PG_DB", "seternak")
	viper.SetDefault("PG_SSL", "disable")

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
		StrictRouting:            true,
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

// fiberErrorHandler manages error handling in the Fiber application context.
// Returns a JSON response with the appropriate status code and error message.
func (c *Config) fiberErrorHandler(ctx *fiber.Ctx, err error) error {

	// Check for errors code
	var appErrors *apperrors.Error
	switch apperrors.Status(err) {
	case fiber.StatusNotFound:
		appErrors = apperrors.NewNotFound(err)
	default:
		appErrors = apperrors.NewInternal(err)
	}

	// Return status code with error message
	return ctx.Status(appErrors.Status()).JSON(domain.CustomResponse{
		HTTPStatusCode: appErrors.Status(),
		ResponseData:   appErrors,
	})
}

// setDB initializes and configures the database connection using GORM with PostgreSQL.
// It logs a fatal error and exits if the database initialization fails.
func (c *Config) setDB() {
	var (
		pgHost = viper.GetString("PG_HOST")
		pgPort = viper.GetString("PG_PORT")
		pgUser = viper.GetString("PG_USER")
		pgPass = viper.GetString("PG_PASS")
		pgDB   = viper.GetString("PG_DB")
		pgSSL  = viper.GetString("PG_SSL")
	)

	loggerDB := zapgorm2.New(logger.Get())
	loggerDB.SetAsDefault()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pgHost, pgPort, pgUser, pgPass, pgDB, pgSSL)
	gormPrepared, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 loggerDB,
	})
	if err != nil {
		log.Fatalf("Error initializing DB, %v", err)
	}

	if viper.GetString("APP_ENV") != consts.ProductionMode {
		gormPrepared = gormPrepared.Debug()
	}

	c.db = gormPrepared
}

// GetDB retrieves the GORM database instance from the Config struct.
// Returns a pointer to the gorm.DB instance.
func (c *Config) GetDB() *gorm.DB {
	return c.db
}
