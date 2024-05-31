package dbinstance

import (
	db "clean-architecture/internal/constants/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBInstance struct {
	*db.Queries
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) DBInstance {
	return DBInstance{
		pool:    pool,
		Queries: db.New(pool),
	}
}
