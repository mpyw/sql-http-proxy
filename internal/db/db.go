// Package db provides database connection management for sql-http-proxy.
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

// ConnectTimeout is the timeout for validating database connection.
const ConnectTimeout = 10 * time.Second

// Connect establishes a database connection based on the configuration.
// Returns nil if the configuration doesn't require a database connection.
// Validates the connection with a ping before returning.
func Connect(cfg config.Config) (*sqlx.DB, error) {
	if !cfg.RequiresDB() {
		return nil, nil
	}

	driverName, err := cfg.Driver()
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open(driverName, cfg.DSN)
	if err != nil {
		return nil, err
	}

	// Validate connection with a ping
	ctx, cancel := context.WithTimeout(context.Background(), ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
