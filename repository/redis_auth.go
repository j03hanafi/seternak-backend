package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/domain/apperrors"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

// redisAuthRepository is a repository that implements domain.AuthRepository interface
type redisAuthRepository struct {
	redis *redis.Client
}

// NewRedisAuth is a factory for initializing Auth Repositories
func NewRedisAuth(redis *redis.Client) domain.AuthRepository {
	return &redisAuthRepository{
		redis: redis,
	}
}

// SetRefreshToken stores a user's refresh token in Redis with a specified expiration time.
// Returns an error if storing the token in Redis fails.
func (r *redisAuthRepository) SetRefreshToken(ctx context.Context, userID, tokenID string, expiresIn time.Duration) error {
	l := logger.FromCtx(ctx)

	// We'll store userID with token id, so we can scan (non-blocking)
	// over the user's tokens and delete them in case of token leakage
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		l.Error("Could not SET refresh token to redis",
			zap.String("userID", userID),
			zap.String("tokenID", tokenID),
			zap.Error(err),
		)
		return apperrors.NewInternal(err)
	}

	return nil
}

// DeleteRefreshToken removes a specified refresh token for a user from Redis.
// Returns an authorization error if the token does not exist, or another error if the deletion process fails.
func (r *redisAuthRepository) DeleteRefreshToken(ctx context.Context, userID, prevTokenID string) error {
	l := logger.FromCtx(ctx)

	key := fmt.Sprintf("%s:%s", userID, prevTokenID)

	result := r.redis.Del(ctx, key)
	if err := result.Err(); err != nil {
		l.Error("Could not delete refresh token to redis",
			zap.String("userID", userID),
			zap.String("prevTokenID", prevTokenID),
			zap.Error(err),
		)
		return apperrors.NewInternal(err)
	}

	// Val returns count of deleted keys
	// If no key was deleted, the refresh token is invalid
	if result.Val() < 1 {
		l.Error("Refresh token does not exist in redis",
			zap.String("userID", userID),
			zap.String("prevTokenID", prevTokenID),
		)
		return apperrors.NewAuthorization(result.Err(), "Invalid refresh token")
	}

	return nil
}

func (r *redisAuthRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	l := logger.FromCtx(ctx)

	pattenr := fmt.Sprintf("%s*", userID)

	iter := r.redis.Scan(ctx, 0, pattenr, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := r.redis.Del(ctx, iter.Val()).Err(); err != nil {
			l.Error("Failed to delete refresh token",
				zap.String("tokenID", iter.Val()),
				zap.Error(err),
			)
			failCount++
		}
	}

	// checks for any errors that might have occurred during iteration
	if err := iter.Err(); err != nil {
		l.Error("Failed to delete refresh token",
			zap.String("tokenID", iter.Val()),
			zap.Error(err),
		)
	}

	if failCount > 0 {
		return apperrors.NewInternal(errors.New("failed to delete refresh token"))
	}

	return nil
}
