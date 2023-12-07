package repository

import (
	"github.com/j03hanafi/seternak-backend/config"
	"github.com/j03hanafi/seternak-backend/domain"
	"github.com/j03hanafi/seternak-backend/utils/consts"
	"github.com/spf13/viper"
	"testing"
)

func initPGUser(t testing.TB) domain.UserRepository {
	t.Helper()

	viper.Set("APP_ENV", consts.TestMode)
	viper.Set("PG_HOST", "localhost")
	cfg := config.New()

	return NewPGUser(cfg.GetDB())

}
