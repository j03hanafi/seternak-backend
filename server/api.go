package server

import (
	"crypto/rsa"
	"github.com/gofiber/fiber/v2"
	"github.com/j03hanafi/seternak-backend/handler"
	"github.com/j03hanafi/seternak-backend/server/middleware"
)

// api struct holds required handlers for api to function
type api struct {
	versionHandler *handler.Version
	userHandler    *handler.User
}

// apiConfig will hold handlers that will eventually be injected into this
// api layer on api initialization
type apiConfig struct {
	app       *fiber.App
	baseURL   string
	version   *handler.Version
	user      *handler.User
	publicKey *rsa.PublicKey
	secretKey string
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
	g.Post("/login", h.userHandler.LogIn)

	g.Post("/tokens", middleware.AuthRefresh(c.secretKey), h.userHandler.Tokens)

	g.Post("/logout", middleware.AuthToken(c.publicKey), h.userHandler.SignOut)
}
