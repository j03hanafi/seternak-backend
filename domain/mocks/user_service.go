package mocks

import (
	"context"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Get(ctx context.Context, uid ulid.ULID) (*domain.User, error) {
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

func (m *MockUserService) LogOut(ctx context.Context, uid ulid.ULID) error {
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
