package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

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
