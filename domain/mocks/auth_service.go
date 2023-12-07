package mocks

import (
	"context"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) NewPairFromUser(ctx context.Context, u *domain.User, prevTokenID string) (*domain.AuthToken, error) {
	args := m.Called(ctx, u, prevTokenID)

	var r0 *domain.AuthToken
	if args.Get(0) != nil {
		r0 = args.Get(0).(*domain.AuthToken)
	}

	var r1 error
	if args.Get(1) != nil {
		r1 = args.Error(1)
	}

	return r0, r1
}
