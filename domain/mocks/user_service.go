package mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Get(ctx context.Context, uid uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, uid)

	var r0 *domain.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.User)
	}

	var r1 error
	if args.Get(1) != nil {
		r1 = args.Error(1)
	}

	return r0, r1
}

func (m *MockUserService) LogOut(ctx context.Context, uid uuid.UUID) error {
	args := m.Called(ctx, uid)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Error(0)
	}

	return r0
}

func (m *MockUserService) SignUp(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Error(0)
	}

	return r0
}

func (m *MockUserService) LogIn(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Error(0)
	}

	return r0

}
