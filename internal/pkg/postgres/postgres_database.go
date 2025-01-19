package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	*pgxpool.Pool
}

const URL = "postgres://%s:%s@%s:%s/%s?sslmode=disable"

func NewDatabase(ctx context.Context, postgresURL string) (Database, error) {
	pool, err := pgxpool.New(ctx, postgresURL)
	if err != nil {
		return Database{}, fmt.Errorf("connect to database: %w", err)
	}

	return Database{
		Pool: pool,
	}, nil
}
