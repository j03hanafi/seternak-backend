package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/j03hanafi/seternak-backend/handler"
)

// api struct holds required handlers for api to function
type api struct {
	versionHandler *handler.Version
	userHandler    *handler.User
}

// apiConfig will hold handlers that will eventually be injected into this
// api layer on api initialization
type apiConfig struct {
	app     *fiber.App
	baseURL string
	version *handler.Version
	user    *handler.User
}

// newAPI initializes the api with required injected handlers along with http routes
func newAPI(c *apiConfig) {
	h := &api{
		versionHandler: c.version,
		userHandler:    c.user,
	}

	// Create a group, or base url for all routes
	g := c.app.Group(c.baseURL)

	g.Get("", h.versionHandler.GetVersion)

	g.Post("/signup", h.userHandler.SignUp)
}
