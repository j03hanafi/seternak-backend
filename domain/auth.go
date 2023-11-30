package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// AuthToken used for returning pairs of id and refresh tokens
type AuthToken struct {
	IDToken
	RefreshToken
}

// RefreshToken stores token properties that
// are accessed in multiple application layer
type RefreshToken struct {
	ID  uuid.UUID `json:"-"`
	UID uuid.UUID `json:"-"`
	SS  string    `json:"refresh_token"`
}

// IDToken stores token properties that
// are accessed in multiple application layers
type IDToken struct {
	SS string `json:"id_token"`
}

// AuthService defines methods the handler layer expects to interact
// with in regard to producing JWTs as string
type AuthService interface {

	// NewPairFromUser generates a new pair of authentication tokens (access and refresh) for a user.
	// Returns the AuthToken pair or an error if the token generation process fails.
	NewPairFromUser(ctx context.Context, u *User, prevTokenID string) (*AuthToken, error)
}

// AuthRepository defines methods it expects a repository
// it interacts with to implement
type AuthRepository interface {

	// SetRefreshToken stores or updates a refresh token for a user in the data source.
	// Returns an error if the token storage operation fails.
	SetRefreshToken(ctx context.Context, userID, tokenID string, expiresIn time.Duration) error

	// DeleteRefreshToken removes a user's refresh token from the data source.
	// Returns an error if the deletion process fails.
	DeleteRefreshToken(ctx context.Context, userID, prevTokenID string) error

	// DeleteUserRefreshTokens removes all refresh tokens associated with a specific user.
	// Returns an error if the deletion operation fails.
	DeleteUserRefreshTokens(ctx context.Context, userID string) error
}
