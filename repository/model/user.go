package model

import (
	"github.com/j03hanafi/seternak-backend/domain"
	"gorm.io/gorm"
	"time"
)

// User defines the schema for the users table in the database.
type User struct {
	UID       string
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// FromUser converts the domain.User struct to a User model.
func (u *User) FromUser(user *domain.User) {
	u.UID = user.UID.String()
	u.Email = user.Email
	u.Password = user.Password
	u.Name = user.Name
}
