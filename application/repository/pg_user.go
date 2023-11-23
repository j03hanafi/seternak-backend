package repository

import (
	"context"
	"errors"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/model"
	"github.com/j03hanafi/seternak-backend/utils/apperrors"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type pgUserRepository struct {
	db *gorm.DB
}

func NewPGUserRepository(db *gorm.DB) domain.UserRepository {
	return &pgUserRepository{
		db: db,
	}
}

func (p *pgUserRepository) Create(ctx context.Context, u *domain.User) error {
	l := logger.FromCtx(ctx)

	user := new(model.User)
	user.FromUser(u)

	err := p.db.WithContext(ctx).Create(user).Error
	if err != nil {
		l.Error("Could not create a user", zap.Error(err))

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.NewConflict("email", u.Email, err)
		}

		return apperrors.NewInternal(err)
	}
	return nil
}
