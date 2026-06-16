package handlers

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Handlers holds shared dependencies for all route handlers.
type Handlers struct {
	DB             *pgxpool.Pool
	GoogleClientID string
	JWTSecret      string
	JWTExpiry      time.Duration
}

func New(pool *pgxpool.Pool, googleClientID, jwtSecret string, jwtExpiry time.Duration) *Handlers {
	return &Handlers{
		DB:             pool,
		GoogleClientID: googleClientID,
		JWTSecret:      jwtSecret,
		JWTExpiry:      jwtExpiry,
	}
}
