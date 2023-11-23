package main

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()

	app := fiber.New(fiber.Config{
		CaseSensitive:            true,
		DisableHeaderNormalizing: true,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
	})

	app.Use(logger.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

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

	log.Println("Server is starting...")

	<-ctx.Done()
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	err = app.ShutdownWithTimeout(5 * time.Second)
	if err != nil {
		log.Fatal("Server forced to shutdown", err)
	}

	log.Println("Server was successful shutdown.")
}
