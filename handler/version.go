package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/j03hanafi/seternak-backend/handler/response"
)

// Version struct holds required services for handler to function
type Version struct {
}

// NewVersion is a factory function for initializing a Version Handler
// with its service layer dependencies
func NewVersion() *Version {
	return &Version{}
}

// GetVersion handler
func (v *Version) GetVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(response.CustomResponse{
		HTTPStatusCode: fiber.StatusOK,
		ResponseData: fiber.Map{
			"version": "1.0.0",
		},
	})
}

func (v *Version) getGitVersion() {

}
