package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// BindParams binds named parameters and rebinds for the driver.
func BindParams(db *sqlx.DB, query string, params map[string]any) (string, []any, error) {
	boundSQL, args, err := sqlx.Named(query, params)
	if err != nil {
		return "", nil, err
	}
	return db.Rebind(boundSQL), args, nil
}

// QueryOne executes a query expecting a single row.
// Returns nil (not error) when no rows found.
func QueryOne(ctx context.Context, db *sqlx.DB, query string, args ...any) (map[string]any, error) {
	row := db.QueryRowxContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	result := map[string]any{}
	if err := row.MapScan(result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

// QueryMany executes a query expecting multiple rows.
func QueryMany(ctx context.Context, db *sqlx.DB, query string, args ...any) ([]map[string]any, error) {
	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	results := make([]map[string]any, 0)
	for rows.Next() {
		entry := map[string]any{}
		if err := rows.MapScan(entry); err != nil {
			return nil, err
		}
		results = append(results, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// ExecResult holds the result of an exec operation.
type ExecResult struct {
	LastInsertId *int64 // nil if not available (non-MySQL drivers)
	RowsAffected int64
}

// Exec executes a query without returning rows.
// For MySQL, it returns the last insert ID. For other drivers, LastInsertId is nil.
func Exec(ctx context.Context, db *sqlx.DB, query string, args ...any) (*ExecResult, error) {
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	execResult := &ExecResult{}

	// RowsAffected is generally available
	if ra, err := result.RowsAffected(); err == nil {
		execResult.RowsAffected = ra
	}

	// LastInsertId is only meaningful for MySQL
	// Other drivers either return 0 or error
	if db.DriverName() == "mysql" {
		if id, err := result.LastInsertId(); err == nil {
			execResult.LastInsertId = &id
		}
	}

	return execResult, nil
}
