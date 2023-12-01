package middleware

import (
	"crypto/rsa"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/j03hanafi/seternak-backend/domain/apperrors"
	"github.com/j03hanafi/seternak-backend/handler/response"
	"github.com/j03hanafi/seternak-backend/utils"
	"github.com/j03hanafi/seternak-backend/utils/consts"
)

func AuthUser(publicKey *rsa.PublicKey) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    publicKey,
		},
		Claims:         &utils.IDTokenCustomClaims{},
		SuccessHandler: authSuccessHandler(),
		ErrorHandler:   authErrorHandler(),
		ContextKey:     consts.JWTContextKey,
	})
}

func authSuccessHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals(consts.JWTContextKey).(*jwt.Token).Claims.(*utils.IDTokenCustomClaims)
		c.Locals(consts.JWTUserContextKey, *claims.User)
		return c.Next()
	}
}

func authErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if err.Error() == consts.ErrBadRequestJWT {
			authErr := apperrors.NewBadRequest(err, consts.ErrBadRequestJWT)
			return c.Status(authErr.Status()).JSON(response.CustomResponse{
				HTTPStatusCode: authErr.Status(),
				ResponseData:   authErr,
			})
		}
		authErr := apperrors.NewAuthorization(err, consts.ErrUnauthorizedJWT)
		return c.Status(authErr.Status()).JSON(response.CustomResponse{
			HTTPStatusCode: authErr.Status(),
			ResponseData:   authErr,
		})
	}
}
