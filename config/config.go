package config

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
	"github.com/j03hanafi/seternak-backend/domain/apperrors"
	"github.com/j03hanafi/seternak-backend/handler/response"
	"github.com/j03hanafi/seternak-backend/utils/consts"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"moul.io/zapgorm2"
	"os"
	"time"
)

// Config defines configuration settings for the server.
type Config struct {
	db          *gorm.DB
	fiber       *fiber.Config
	fiberzap    *fiberzap.Config
	recover     *recover.Config
	redisClient *redis.Client
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	zapLogger   *zap.Logger
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
	config.setLogger()
	config.setDB()
	config.setFiberConfig()
	config.setFiberzapConfig()
	config.setRecoverConfig()
	config.setRedis()
	config.setRSAKeys()

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

	// Set default Redis configuration
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")

	// Set default Auth configuration
	viper.SetDefault("REFRESH_TOKEN_SECRET", "refresh_token_secret")
	viper.SetDefault("PRIVATE_KEY_FILE", "./rsa_private.pem")
	viper.SetDefault("PUBLIC_KEY_FILE", "./rsa_public.pem")
	viper.SetDefault("ID_TOKEN_EXP", "900")         // 15 minutes
	viper.SetDefault("REFRESH_TOKEN_EXP", "259200") // 3 days

}

// setLogger initializes from logger utils package.
func (c *Config) setLogger() {
	c.zapLogger = logger.Get()
}

// GetLogger retrieves the zap logger instance from the Config struct.
// Returns a pointer to the zap.Logger instance.
func (c *Config) GetLogger() *zap.Logger {
	return c.zapLogger
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
		Logger: c.zapLogger,
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
	case fiber.StatusConflict:
		appErrors = apperrors.NewConflict(err)
	case fiber.StatusNotFound:
		appErrors = apperrors.NewNotFound(err)
	case fiber.StatusUnauthorized:
		appErrors = apperrors.NewAuthorization(err)
	default:
		appErrors = apperrors.NewInternal(err)
	}

	// Return status code with error message
	return ctx.Status(appErrors.Status()).JSON(response.CustomResponse{
		HTTPStatusCode: appErrors.Status(),
		ResponseData:   appErrors,
	})
}

// setDB initializes and configures the database connection using GORM with PostgreSQL.
// It logs a fatal error and exits if the database initialization fails.
func (c *Config) setDB() {
	var (
		l      = c.zapLogger
		pgHost = viper.GetString("PG_HOST")
		pgPort = viper.GetString("PG_PORT")
		pgUser = viper.GetString("PG_USER")
		pgPass = viper.GetString("PG_PASS")
		pgDB   = viper.GetString("PG_DB")
		pgSSL  = viper.GetString("PG_SSL")
	)

	loggerDB := zapgorm2.New(l)
	loggerDB.SetAsDefault()

	l.Info("Connecting to Postgres...")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pgHost, pgPort, pgUser, pgPass, pgDB, pgSSL)
	gormPrepared, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 loggerDB,
	})
	if err != nil {
		log.Fatalf("Error initializing DB, %v", err)
	}

	// Register db resolver
	err = gormPrepared.Use(
		dbresolver.Register(dbresolver.Config{}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)
	if err != nil {
		log.Fatalf("Error setting up db resolver, %v", err)
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

// setRedis initializes and configures the Redis connection.
// It logs a fatal error and exits if the Redis initialization fails.
func (c *Config) setRedis() {
	var (
		l         = c.zapLogger
		redisHost = viper.GetString("REDIS_HOST")
		redisPort = viper.GetString("REDIS_PORT")
	)

	l.Info("Connecting to Redis...")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		l.Fatal("Error initializing Redis", zap.Error(err))
	}

	c.redisClient = redisClient
}

// GetRedis retrieves the Redis client instance from the Config struct.
// Returns a pointer to the redis.Client instance.
func (c *Config) GetRedis() *redis.Client {
	return c.redisClient
}

// setRSAKeys initializes and sets the RSA keys for JWT signing and verification.
// It logs a fatal error and exits if the RSA key initialization fails.
func (c *Config) setRSAKeys() {
	l := c.zapLogger

	privateKeyFile, err := os.ReadFile(viper.GetString("PRIVATE_KEY_FILE"))
	if err != nil {
		l.Fatal("Error reading private key file", zap.Error(err))
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		l.Fatal("Error parsing private key", zap.Error(err))
	}

	c.privateKey = privateKey

	publicKeyFile, err := os.ReadFile(viper.GetString("PUBLIC_KEY_FILE"))
	if err != nil {
		l.Fatal("Error reading public key file", zap.Error(err))
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	if err != nil {
		l.Fatal("Error parsing public key", zap.Error(err))
	}

	c.publicKey = publicKey
}

// GetPrivateKey retrieves the RSA private key from the Config struct.
// Returns a pointer to the rsa.PrivateKey instance.
func (c *Config) GetPrivateKey() *rsa.PrivateKey {
	return c.privateKey
}

// GetPublicKey retrieves the RSA public key from the Config struct.
// Returns a pointer to the rsa.PublicKey instance.
func (c *Config) GetPublicKey() *rsa.PublicKey {
	return c.publicKey
}

// Close to be used in graceful server shutdown
func (c *Config) Close() error {
	l := c.zapLogger

	db, _ := c.db.DB()
	if err := db.Close(); err != nil {
		l.Error("Error closing DB connection", zap.Error(err))
		return err
	}

	if err := c.redisClient.Close(); err != nil {
		l.Error("Error closing Redis connection", zap.Error(err))
		return err
	}

	return nil
}
