package main

import (
	"context"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	configuration "github.com/j03hanafi/seternak-backend/config"
	"github.com/j03hanafi/seternak-backend/logger"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := configuration.New()
	var err error

	l := logger.Get()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()

	app := fiber.New(*config.GetFiberConfig())

	app.Use(fiberzap.New(*config.GetFiberzapConfig()))
	app.Use(recover.New(*config.GetRecoverConfig()))

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]any{
			"message": "Hello, World!",
			"version": "0.1.0",
		})
	})

	app.Get("/panic", func(ctx *fiber.Ctx) error {
		panic("I'm an error")
	})

	go func() {
		if err = app.Listen(":8080"); err != nil {
			log.Fatal("Server Error", err)
		}
	}()

	l.Info("Server is starting...")

	<-ctx.Done()
	stop()
	l.Info("shutting down gracefully, press Ctrl+C again to force")

	err = app.ShutdownWithTimeout(5 * time.Second)
	if err != nil {
		l.Fatal("Server forced to shutdown", zap.Error(err))
	}

	l.Info("Server was successful shutdown.")
}
