package main

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		CaseSensitive:            true,
		DisableHeaderNormalizing: true,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
		Prefork:                  true,
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World 👋!")
	})

	app.Get("/panic", func(ctx *fiber.Ctx) error {
		panic("I'm an error")
	})

	app.Listen(":8080")
}
