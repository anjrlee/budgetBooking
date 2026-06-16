package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// New creates a pgx connection pool using the DATABASE_URL environment variable.
func New(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
}
