package db

import (
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
