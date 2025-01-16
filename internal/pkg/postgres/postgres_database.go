package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	*pgxpool.Pool
}

func NewDatabase(ctx context.Context, connString string) (Database, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return Database{}, fmt.Errorf("connect to database: %w", err)
	}

	return Database{
		Pool: pool,
	}, nil
}
