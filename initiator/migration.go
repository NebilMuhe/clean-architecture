package initiator

import (
	"clean-architecture/utils/logger"
	"context"

	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
)

func InitMigration(ctx context.Context, file, murl string, log logger.Logger) {
	migration, err := migrate.New(file, murl)
	if err != nil {
		log.Fatal(ctx, "unable to migrate instance", zap.Error(err))
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(ctx, "unable to migrate up", zap.Error(err))
	}
}
