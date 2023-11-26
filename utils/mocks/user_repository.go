package mocks

import (
	"context"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Error(0)
	}

	return r0
}
