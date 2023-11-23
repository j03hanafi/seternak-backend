package controller

import "github.com/gofiber/fiber/v2"

type Version struct {
}

func NewVersion() *Version {
	return &Version{}
}

func (v *Version) GetVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"version": "1.0.0",
	})
}
