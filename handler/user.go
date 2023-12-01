package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/domain/apperrors"
	"github.com/j03hanafi/seternak-backend/handler/request"
	"github.com/j03hanafi/seternak-backend/handler/response"
	"github.com/j03hanafi/seternak-backend/utils/consts"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"go.uber.org/zap"
)

// User struct holds required services for handler to function
type User struct {
	userService domain.UserService
	authService domain.AuthService
}

// UserHandlerConfig will hold services that will eventually be injected into this
// handler layer
type UserHandlerConfig struct {
	UserService domain.UserService
	AuthService domain.AuthService
}

// NewUser is a factory function for initializing a User Handler
// with its service layer dependencies
func NewUser(c *UserHandlerConfig) *User {
	u := new(User)

	if c.UserService != nil {
		u.userService = c.UserService
	}

	if c.AuthService != nil {
		u.authService = c.AuthService
	}

	return u
}

// SignUp handler
func (u *User) SignUp(c *fiber.Ctx) error {
	ctx := c.UserContext()
	l := logger.FromCtx(ctx)

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

func (u *User) SignIn(c *fiber.Ctx) error {
	ctx := c.UserContext()
	l := logger.FromCtx(ctx)

	// bind request body to SignIn struct
	req := new(request.SignIn)
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
	}

	// sign in user
	if err := u.userService.SignIn(ctx, user); err != nil {
		l.Info("Unable to sign in user",
			zap.Error(err),
		)
		return c.Status(apperrors.Status(err)).JSON(response.CustomResponse{
			HTTPStatusCode: apperrors.Status(err),
			ResponseData:   err,
		})
	}

	// create token pair as strings
	authToken, err := u.authService.NewPairFromUser(ctx, user, "")
	if err != nil {
		l.Info("Unable to create token pair for user",
			zap.Error(err),
		)

		// may eventually implement rollback logic here
		// meaning, if we fail to create tokens after creating a user,
		// we make sure to clear/delete the created user in the database

		return c.Status(apperrors.Status(err)).JSON(response.CustomResponse{
			HTTPStatusCode: apperrors.Status(err),
			ResponseData:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CustomResponse{
		HTTPStatusCode: fiber.StatusOK,
		ResponseData: fiber.Map{
			"tokens": authToken,
		},
	})

}

func (u *User) SignOut(c *fiber.Ctx) error {
	ctx := c.UserContext()
	l := logger.FromCtx(ctx)

	user := c.Locals(consts.JWTUserContextKey).(domain.User)

	// sign out user
	if err := u.userService.SignOut(ctx, user.UID); err != nil {
		l.Info("Unable to sign out user",
			zap.Error(err),
		)
		return c.Status(apperrors.Status(err)).JSON(response.CustomResponse{
			HTTPStatusCode: apperrors.Status(err),
			ResponseData:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CustomResponse{
		HTTPStatusCode: fiber.StatusOK,
		ResponseData:   "Successfully signed out",
	})

}

func (u *User) Tokens(c *fiber.Ctx) error {
	ctx := c.UserContext()
	l := logger.FromCtx(ctx)

	refreshToken := c.Locals(consts.JWTRefreshTokenContextKey).(*domain.RefreshToken)

	// get up-to-date user data
	user, err := u.userService.Get(ctx, refreshToken.UID)
	if err != nil {
		l.Error("Unable to get user",
			zap.Error(err),
		)
		return c.Status(apperrors.Status(err)).JSON(response.CustomResponse{
			HTTPStatusCode: apperrors.Status(err),
			ResponseData:   err,
		})
	}

	// create token pair as strings
	authToken, err := u.authService.NewPairFromUser(ctx, user, refreshToken.ID.String())
	if err != nil {
		l.Info("Unable to create token pair for user",
			zap.Error(err),
		)

		// may eventually implement rollback logic here
		// meaning, if we fail to create tokens after creating a user,
		// we make sure to clear/delete the created user in the database

		return c.Status(apperrors.Status(err)).JSON(response.CustomResponse{
			HTTPStatusCode: apperrors.Status(err),
			ResponseData:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CustomResponse{
		HTTPStatusCode: fiber.StatusOK,
		ResponseData: fiber.Map{
			"tokens": authToken,
		},
	})
}
