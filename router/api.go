package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/j03hanafi/seternak-backend/application/handler"
)

// api struct holds required handlers for api to function
type api struct {
	versionController *handler.Version
}

// apiConfig will hold handlers that will eventually be injected into this
// api layer on api initialization
type apiConfig struct {
	app     *fiber.App
	baseURL string
	version *handler.Version
}

// newAPI initializes the api with required injected handlers along with http routes
func newAPI(c *apiConfig) {
	h := &api{
		versionController: c.version,
	}

	// Create a group, or base url for all routes
	g := c.app.Group(c.baseURL)

	g.Get("", h.versionController.GetVersion)
}
