package service

import (
	"context"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/domain/apperrors"
	"github.com/j03hanafi/seternak-backend/utils"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"go.uber.org/zap"
)

// userService acts as a struct for injecting an implementation of repositories
// for use in service methods
type userService struct {
	userRepository domain.UserRepository
}

// UserServiceConfig will hold repositories that will eventually be injected into this
// service layer
type UserServiceConfig struct {
	UserRepository domain.UserRepository
}

// NewUser is a factory function for
// initializing a userService with its repository layer dependencies
func NewUser(c *UserServiceConfig) domain.UserService {
	service := new(userService)

	if c.UserRepository != nil {
		service.userRepository = c.UserRepository
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

	if err := u.userRepository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *userService) SignIn(ctx context.Context, user *domain.User) error {
	l := logger.FromCtx(ctx)

	userFetch, err := u.userRepository.FindByEmail(ctx, user.Email)
	if err != nil {
		l.Error("error fetching user",
			zap.Error(err),
			zap.Any("user", user),
		)
		return apperrors.NewAuthorization(err, err.Error())
	}

	match, err := utils.ComparePasswords(userFetch.Password, user.Password)
	if err != nil {
		l.Error("error comparing passwords",
			zap.Error(err),
		)
		return apperrors.NewInternal(err)
	}

	if !match {
		l.Error("passwords do not match",
			zap.Error(err),
		)
		return apperrors.NewAuthorization(nil, "Invalid email/password combination")
	}

	*user = *userFetch

	return nil

}
