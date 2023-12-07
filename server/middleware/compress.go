package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func Compression() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelDefault,
	})
}
