package server

import (
	"crypto/rsa"
	"github.com/gofiber/fiber/v2"
	"github.com/j03hanafi/seternak-backend/handler"
	"github.com/j03hanafi/seternak-backend/server/middleware"
	"go.uber.org/zap"
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
	zapLogger *zap.Logger
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

	// Create a group, or base url for all routes and middleware that will be used for all API
	g := c.app.Group(c.baseURL).Use(middleware.Logger(c.zapLogger), middleware.Compression())

	g.Get("", h.versionHandler.GetVersion)

	g.Post("/signup", h.userHandler.SignUp)
	g.Post("/login", h.userHandler.LogIn)

	g.Post("/tokens", middleware.AuthRefresh(c.secretKey), h.userHandler.Tokens)

	g.Post("/logout", middleware.AuthToken(c.publicKey), h.userHandler.SignOut)
}
