package service

import (
	"context"
	"crypto/rsa"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/domain/apperrors"
	"github.com/j03hanafi/seternak-backend/utils"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"go.uber.org/zap"
)

// authService acts as a struct for injecting an implementation of repositories
// for use in service methods
type authService struct {
	authRepository             domain.AuthRepository
	privateKey                 *rsa.PrivateKey
	refreshTokenSecret         string
	idTokenExpirationSecs      int64
	refreshTokenExpirationSecs int64
}

// AuthServiceConfig will hold repositories that will eventually be injected into this
// service layer
type AuthServiceConfig struct {
	AuthRepository             domain.AuthRepository
	PrivateKey                 *rsa.PrivateKey
	RefreshTokenSecret         string
	IDTokenExpirationSecs      int64
	RefreshTokenExpirationSecs int64
}

// NewAuth is a factory function for
// initializing a authService with its repository layer dependencies
func NewAuth(c *AuthServiceConfig) domain.AuthService {
	service := new(authService)

	if c.AuthRepository != nil {
		service.authRepository = c.AuthRepository
	}

	if c.PrivateKey != nil {
		service.privateKey = c.PrivateKey
	}

	if c.RefreshTokenSecret != "" {
		service.refreshTokenSecret = c.RefreshTokenSecret
	}

	if c.IDTokenExpirationSecs != 0 {
		service.idTokenExpirationSecs = c.IDTokenExpirationSecs
	}

	if c.RefreshTokenExpirationSecs != 0 {
		service.refreshTokenExpirationSecs = c.RefreshTokenExpirationSecs
	}

	return service
}

// NewPairFromUser creates a new set of ID and refresh tokens for a user, replacing any previous refresh token.
// Returns an AuthToken pair or an error if the token creation or storage process fails.
func (s *authService) NewPairFromUser(ctx context.Context, u *domain.User, prevTokenID string) (*domain.AuthToken, error) {
	l := logger.FromCtx(ctx)

	// delete user's current refresh token (used when refreshing idToken)
	if prevTokenID != "" {
		if err := s.authRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			l.Error("Error deleting previous refresh token for user",
				zap.Error(err),
			)
			return nil, apperrors.NewInternal(err)
		}
	}

	// No need to use a repository for idToken as it is unrelated to any data source
	idToken, err := utils.GenerateIDToken(u, s.privateKey, s.idTokenExpirationSecs)
	if err != nil {
		l.Error("Error generating ID Token for user",
			zap.Error(err),
		)
		return nil, apperrors.NewInternal(err)
	}

	refreshToken, err := utils.GenerateRefreshToken(u.UID, s.refreshTokenSecret, s.refreshTokenExpirationSecs)
	if err != nil {
		l.Error("Error generating Refresh Token for user",
			zap.Error(err),
		)
		return nil, apperrors.NewInternal(err)
	}

	// set refresh tokens by calling TokenRepository methods
	if err = s.authRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID.String(), refreshToken.ExpiresIn); err != nil {
		l.Error("Error saving refresh token for user",
			zap.Error(err),
		)
		return nil, apperrors.NewInternal(err)
	}

	return &domain.AuthToken{
		IDToken:      domain.IDToken{SS: idToken},
		RefreshToken: domain.RefreshToken{SS: refreshToken.SS, ID: refreshToken.ID, UID: u.UID},
	}, nil

}
