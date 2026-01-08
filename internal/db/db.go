// Package db provides database connection management for sql-http-proxy.
package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

// Connect establishes a database connection based on the configuration.
// Returns nil if the configuration doesn't require a database connection.
func Connect(cfg config.Config) (*sqlx.DB, error) {
	if !cfg.RequiresDB() {
		return nil, nil
	}

	driverName, err := cfg.Driver()
	if err != nil {
		return nil, err
	}

	return sqlx.Open(driverName, cfg.DSN)
}
