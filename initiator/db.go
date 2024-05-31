package initiator

import (
	"clean-architecture/utils/logger"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func InitDB(ctx context.Context, dburl string, log logger.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dburl)
	if err != nil {
		log.Fatal(ctx, "unable to parse", zap.Error(err))
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(ctx, "unable to connect", zap.Error(err))
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal(ctx, "unable to ping database", zap.Error(err))
	}

	return pool
}
