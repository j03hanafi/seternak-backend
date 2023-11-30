package domain

import (
	"context"
	"github.com/google/uuid"
)

// User defines domain model and its json representation
type User struct {
	UID      uuid.UUID `json:"uid"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Name     string    `json:"name"`
}

// UserService defines methods the handler layer expects
// any service it interacts with to implement
type UserService interface {

	// SignUp registers a new user into the system.
	// Returns an error if the user registration process fails.
	SignUp(ctx context.Context, u *User) error

	// SignIn authenticates a user based on provided credentials.
	// Returns an error if authentication or any part of the sign-in process fails.
	SignIn(ctx context.Context, u *User) error

	// SignOut handles the user sign-out process.
	// Returns an error if the sign-out process encounters any issues.
	SignOut(ctx context.Context, uid uuid.UUID) error
}

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement
type UserRepository interface {

	// Create inserts a new User record into the database.
	// Returns an error if the user creation process fails.
	Create(ctx context.Context, u *User) error

	// FindByEmail retrieves a user by their email address from the database.
	// Returns a User object or an error if the user is not found.
	FindByEmail(ctx context.Context, email string) (*User, error)
}
