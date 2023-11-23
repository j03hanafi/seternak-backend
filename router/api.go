package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/j03hanafi/seternak-backend/controller"
)

// Handler struct holds required services for handler to function
type handler struct {
	versionController *controller.Version
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type handlerConfig struct {
	app     *fiber.App
	baseURL string
	version *controller.Version
}

// NewHandler initializes the handler with required injected services along with http routes
// Does not return as it deals directly with a reference to the gin Engine
func newHandler(c *handlerConfig) {
	h := &handler{
		versionController: c.version,
	}

	// Create a group, or base url for all routes
	g := c.app.Group(c.baseURL)

	g.Get("", h.versionController.GetVersion)
}
