package utils

import (
	"crypto/rsa"
	"github.com/google/uuid"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"go.uber.org/zap"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/j03hanafi/seternak-backend/domain"
)

// IDTokenCustomClaims holds structure of jwt claims of idToken
type IDTokenCustomClaims struct {
	User *domain.User `json:"user"`
	jwt.RegisteredClaims
}

// GenerateIDToken generates an IDToken which is a jwt with myCustomClaims
// Could call this GenerateIDTokenString, but the signature makes this fairly clear
func GenerateIDToken(u *domain.User, key *rsa.PrivateKey, exp int64) (string, error) {
	l := logger.Get()

	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp) * time.Second) // 15 minutes

	claims := IDTokenCustomClaims{
		User: u,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(tokenExp),
			NotBefore: jwt.NewNumericDate(currentTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedString, err := token.SignedString(key)
	if err != nil {
		l.Error("Error signing token",
			zap.Error(err),
		)
		return "", err
	}

	return signedString, nil
}

// RefreshTokenData holds the actual signed jwt string along with the ID
// We return the id, so it can be used without re-parsing the JWT from signed string
type RefreshTokenData struct {
	SS        string
	ID        uuid.UUID
	ExpiresIn time.Duration
}

// RefreshTokenCustomClaims holds the payload of a refresh token
// This can be used to extract user id for subsequent
// application operations (IE, fetch user in Redis)
type RefreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.RegisteredClaims
}

// GenerateRefreshToken creates a refresh token
// The refresh token stores only the user's ID, a string
func GenerateRefreshToken(uid uuid.UUID, key string, exp int64) (*RefreshTokenData, error) {
	l := logger.Get()

	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp) * time.Second) // 3 days
	tokenID, err := uuid.NewRandom()
	if err != nil {
		l.Error("Error generating token id",
			zap.Error(err),
		)
	}

	claims := RefreshTokenCustomClaims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(tokenExp),
			ID:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(key))
	if err != nil {
		l.Error("Error signing token",
			zap.Error(err),
		)
		return nil, err
	}

	return &RefreshTokenData{
		SS:        signedString,
		ID:        tokenID,
		ExpiresIn: tokenExp.Sub(currentTime),
	}, nil
}
