package model

import (
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"time"
)

// User defines the schema for the users table in the database.
type User struct {
	UID       ulid.ULID
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// FromUser converts the domain.User struct to a User model.
func (u *User) FromUser(user *domain.User) {
	u.UID = user.UID
	u.Email = user.Email
	u.Password = user.Password
	u.Name = user.Name
}

// ToUser converts the User model to a domain.User struct.
func (u *User) ToUser() *domain.User {
	return &domain.User{
		UID:      u.UID,
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
	}
}
