package e2e

import (
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)

	t.Cleanup(func() { _ = db.Close() })

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			email TEXT
		)
	`)
	require.NoError(t, err)

	_, err = db.Exec(`
		INSERT INTO users (id, first_name, last_name, email) VALUES
		(1, 'John', 'Doe', 'john@example.com'),
		(2, 'Jane', 'Smith', 'jane@example.com'),
		(3, 'Bob', 'Wilson', 'bob@example.com')
	`)
	require.NoError(t, err)

	return db
}

func loadConfig(t *testing.T, name string) config.Config {
	t.Helper()

	cfg, err := config.ParseFile(name)
	require.NoError(t, err)
	return cfg
}

func createHandler(t *testing.T, db *sqlx.DB, cfg config.Config) http.Handler {
	t.Helper()

	require.NotEmpty(t, cfg.Queries)

	handler, err := server.CreateHandler(db, cfg.Queries[0])
	require.NoError(t, err)
	return handler
}

func createMutationHandler(t *testing.T, db *sqlx.DB, cfg config.Config) http.Handler {
	t.Helper()

	require.NotEmpty(t, cfg.Mutations)

	handler, err := server.CreateMutationHandler(db, cfg.Mutations[0])
	require.NoError(t, err)
	return handler
}
