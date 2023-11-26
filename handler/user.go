package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/handler/request"
	"github.com/j03hanafi/seternak-backend/handler/response"
	"github.com/j03hanafi/seternak-backend/utils/apperrors"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"go.uber.org/zap"
)

// User struct holds required services for handler to function
type User struct {
	userService domain.UserService
}

// UserHandlerConfig will hold services that will eventually be injected into this
// handler layer
type UserHandlerConfig struct {
	UserService domain.UserService
}

// NewUser is a factory function for initializing a User Handler
// with its service layer dependencies
func NewUser(c *UserHandlerConfig) *User {
	u := new(User)

	if c.UserService != nil {
		u.userService = c.UserService
	}

	return u
}

// SignUp handler
func (u *User) SignUp(c *fiber.Ctx) error {
	ctx := c.UserContext()
	l := logger.Get()
	ctx = logger.WithCtx(ctx, l)

	// bind request body to SignUp struct
	req := new(request.SignUp)
	if err := c.BodyParser(req); err != nil {
		l.Error("error binding data",
			zap.Error(err),
		)
		return apperrors.NewInternal(err)
	}

	// validate request body
	if err := req.Validate(); err != nil {
		l.Error("error validating data",
			zap.Error(err),
		)
		return apperrors.NewInternal(err)
	}

	// create user domain object
	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	// sign up user
	if err := u.userService.SignUp(ctx, user); err != nil {
		l.Info("Unable to sign up user",
			zap.Error(err),
		)
		return c.Status(apperrors.Status(err)).JSON(response.CustomResponse{
			HTTPStatusCode: apperrors.Status(err),
			ResponseData:   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.CustomResponse{
		HTTPStatusCode: fiber.StatusCreated,
	})
}
