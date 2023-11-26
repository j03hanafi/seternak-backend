package repository

import (
	"context"
	"errors"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/domain/apperrors"
	"github.com/j03hanafi/seternak-backend/repository/model"
	"github.com/j03hanafi/seternak-backend/utils/logger"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// pgUserRepository is data/repository implementation of domain.UserRepository
type pgUserRepository struct {
	db *gorm.DB
}

// NewPGUserRepository is a factory for initializing domain.UserRepository
func NewPGUserRepository(db *gorm.DB) domain.UserRepository {
	return &pgUserRepository{
		db: db,
	}
}

// Create inserts a new user record into the Postgresql database and handles potential errors.
// Returns a conflict error for duplicate email or a general internal error if the operation fails.
func (p *pgUserRepository) Create(ctx context.Context, u *domain.User) error {
	l := logger.FromCtx(ctx)

	user := new(model.User)
	user.FromUser(u)

	// UID is generated by the database, so we omit it from the creation
	err := p.db.WithContext(ctx).Omit("uid").Create(user).Error
	if err != nil {
		l.Error("Could not create a user", zap.Error(err))

		// Check if the error is a duplicate email error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.NewConflict(err, map[string]any{"email": u.Email})
		}

		return apperrors.NewInternal(err)
	}
	return nil
}

func (p *pgUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	l := logger.FromCtx(ctx)

	user := new(model.User)

	err := p.db.WithContext(ctx).Where("email = ?", email).First(user).Error
	if err != nil {
		l.Error("Could not find a user", zap.Error(err), zap.String("email", email))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFound(err, map[string]any{"email": email})
		}

		return nil, apperrors.NewInternal(err)
	}

	return user.ToUser(), nil

}
