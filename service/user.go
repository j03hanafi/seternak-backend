package service

import (
	"context"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/utils"
	"github.com/j03hanafi/seternak-backend/utils/apperrors"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"go.uber.org/zap"
)

// userService acts as a struct for injecting an implementation of repositories
// for use in service methods
type userService struct {
	UserRepository domain.UserRepository
}

// UserServiceConfig will hold repositories that will eventually be injected into this
// service layer
type UserServiceConfig struct {
	UserRepository domain.UserRepository
}

// NewUserService is a factory function for
// initializing a userService with its repository layer dependencies
func NewUserService(c *UserServiceConfig) domain.UserService {
	service := new(userService)

	if c.UserRepository != nil {
		service.UserRepository = c.UserRepository
	}

	return service
}

// SignUp handles the user registration process including password hashing and user creation in the repository.
// Returns an error if the sign-up process or any of its steps fail.
func (u *userService) SignUp(ctx context.Context, user *domain.User) error {
	l := logger.FromCtx(ctx)

	pw, err := utils.HashPassword(user.Password)
	if err != nil {
		l.Error("error hashing password",
			zap.Error(err),
		)
		return apperrors.NewInternal(err)
	}

	user.Password = pw

	if err := u.UserRepository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}
