package middleware

import (
	"crypto/rsa"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/domain/apperrors"
	"github.com/j03hanafi/seternak-backend/handler/response"
	"github.com/j03hanafi/seternak-backend/utils"
	"github.com/j03hanafi/seternak-backend/utils/consts"
)

/*
	Authentication for ID Token
*/

func AuthToken(publicKey *rsa.PublicKey) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    publicKey,
		},
		Claims:         &utils.IDTokenCustomClaims{},
		SuccessHandler: authTokenSuccessHandler(),
		ErrorHandler:   authTokenErrorHandler(),
		ContextKey:     consts.JWTContextKey,
	})
}

func authTokenSuccessHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals(consts.JWTContextKey).(*jwt.Token).Claims.(*utils.IDTokenCustomClaims)
		c.Locals(consts.JWTUserContextKey, *claims.User)
		return c.Next()
	}
}

func authTokenErrorHandler() fiber.ErrorHandler {
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

/*
	Authentication for Refresh Token
*/

func AuthRefresh(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(secret),
		},
		Claims:         &utils.RefreshTokenCustomClaims{},
		SuccessHandler: AuthRefreshSuccessHandler(),
		ErrorHandler:   AuthRefreshErrorHandler(),
		ContextKey:     consts.JWTContextKey,
	})
}

func AuthRefreshSuccessHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Locals(consts.JWTContextKey).(*jwt.Token).Raw
		claims := c.Locals(consts.JWTContextKey).(*jwt.Token).Claims.(*utils.RefreshTokenCustomClaims)

		// Registered claims store ID as a string
		// parse claims.ID as a uuid
		tokenID, err := uuid.Parse(claims.ID)
		if err != nil {
			authErr := apperrors.NewBadRequest(err, consts.ErrBadRequestJWT)
			return c.Status(authErr.Status()).JSON(response.CustomResponse{
				HTTPStatusCode: authErr.Status(),
				ResponseData:   authErr,
			})
		}

		c.Locals(consts.JWTRefreshTokenContextKey, &domain.RefreshToken{
			ID:  tokenID,
			UID: claims.UID,
			SS:  tokenString,
		})
		return c.Next()
	}
}

func AuthRefreshErrorHandler() fiber.ErrorHandler {
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
