package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/j03hanafi/seternak-backend/config"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/domain/apperrors"
	"github.com/j03hanafi/seternak-backend/utils/consts"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPgUserRepository_Create(t *testing.T) {
	t.Parallel()

	r := initPGUser(t)

	t.Run("Success create a user", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		ctx = logger.WithCtx(ctx, logger.Get())
		uid, _ := uuid.NewRandom()
		u := &domain.User{
			UID:      uid,
			Email:    "j03hanafi@email.com",
			Password: "password",
			Name:     "Joenathan Hanafi",
		}

		err := r.Create(ctx, u)
		assert.NoError(t, err)
	})

	t.Run("Error create a duplicate user", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		ctx = logger.WithCtx(ctx, logger.Get())
		uid, _ := uuid.NewRandom()
		u := &domain.User{
			UID:      uid,
			Email:    "j03hanafi@email.com",
			Password: "password",
			Name:     "Joenathan Hanafi",
		}

		err := r.Create(ctx, u)
		assert.Error(t, err)

		assert.EqualError(t, err, apperrors.NewConflict(err).Error())
	})
}

func initPGUser(t testing.TB) domain.UserRepository {
	t.Helper()

	viper.Set("APP_ENV", consts.TestMode)
	viper.Set("PG_HOST", "localhost")
	cfg := config.New()

	return NewPGUser(cfg.GetDB())

}
