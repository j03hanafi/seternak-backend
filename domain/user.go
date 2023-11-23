package domain

import (
	"context"
	"github.com/google/uuid"
)

type User struct {
	UID      uuid.UUID `json:"uid" db:"uid"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"-" db:"password"`
	Name     string    `json:"name" db:"name"`
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
}
