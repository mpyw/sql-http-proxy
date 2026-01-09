// Package db provides database connection management for sql-http-proxy.
package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

// ConnectTimeout is the timeout for validating database connection.
const ConnectTimeout = 10 * time.Second

// Connect establishes a database connection based on the configuration.
// Returns nil if the configuration doesn't require a database connection.
// Validates the connection with a ping before returning.
// If database.init is configured, executes the init SQL after connecting.
// configDir is used to resolve relative paths in sql_files.
func Connect(cfg config.Config, configDir string) (*sqlx.DB, error) {
	if !cfg.RequiresDB() {
		return nil, nil
	}

	driverName, err := cfg.Driver()
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open(driverName, cfg.DSN())
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

	// Execute init SQL if configured
	if err := executeInit(db, cfg, configDir); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to execute database.init: %w", err)
	}

	return db, nil
}

// executeInit runs the database.init SQL if configured.
func executeInit(db *sqlx.DB, cfg config.Config, configDir string) error {
	if cfg.Database == nil || cfg.Database.Init == nil {
		return nil
	}

	init := cfg.Database.Init

	// Execute inline SQL
	if init.SQL != "" {
		if _, err := db.Exec(init.SQL); err != nil {
			return fmt.Errorf("inline sql: %w", err)
		}
	}

	// Execute SQL files
	for _, file := range init.SQLFiles {
		// Resolve relative path from config directory
		if !filepath.IsAbs(file) {
			file = filepath.Join(configDir, file)
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("sql_files[%s]: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("sql_files[%s]: %w", file, err)
		}
	}

	return nil
}
