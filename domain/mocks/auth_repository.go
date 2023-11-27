package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) SetRefreshToken(ctx context.Context, userID, tokenID string, expiresIn time.Duration) error {
	args := m.Called(ctx, userID, tokenID, expiresIn)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Error(0)
	}

	return r0
}

func (m *MockAuthRepository) DeleteRefreshToken(ctx context.Context, userID, prevTokenID string) error {
	args := m.Called(ctx, userID, prevTokenID)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Error(0)
	}

	return r0
}
