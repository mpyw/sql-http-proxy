package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestExecuteInit_InlineSQL(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	cfg := config.Config{
		Database: &config.DatabaseConfig{
			DSN: "sqlite::memory:",
			Init: &config.DatabaseInit{
				SQL: "CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)",
			},
		},
	}

	err = executeInit(db, cfg, "")
	require.NoError(t, err)

	// Verify table was created
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM test")
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestExecuteInit_SQLFiles(t *testing.T) {
	// Create a temporary directory for SQL files
	tmpDir := t.TempDir()

	// Create SQL file
	sqlFile := filepath.Join(tmpDir, "init.sql")
	err := os.WriteFile(sqlFile, []byte(`
		CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT);
		INSERT INTO test (id, name) VALUES (1, 'Alice');
	`), 0644)
	require.NoError(t, err)

	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	cfg := config.Config{
		Database: &config.DatabaseConfig{
			DSN: "sqlite::memory:",
			Init: &config.DatabaseInit{
				SQLFiles: []string{"init.sql"},
			},
		},
	}

	err = executeInit(db, cfg, tmpDir)
	require.NoError(t, err)

	// Verify data was inserted
	var name string
	err = db.Get(&name, "SELECT name FROM test WHERE id = 1")
	require.NoError(t, err)
	assert.Equal(t, "Alice", name)
}

func TestExecuteInit_BothSQLAndFiles(t *testing.T) {
	// Create a temporary directory for SQL files
	tmpDir := t.TempDir()

	// Create SQL file
	sqlFile := filepath.Join(tmpDir, "seed.sql")
	err := os.WriteFile(sqlFile, []byte(`
		INSERT INTO test (id, name) VALUES (2, 'Bob');
	`), 0644)
	require.NoError(t, err)

	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	cfg := config.Config{
		Database: &config.DatabaseConfig{
			DSN: "sqlite::memory:",
			Init: &config.DatabaseInit{
				SQL:      "CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT); INSERT INTO test (id, name) VALUES (1, 'Alice');",
				SQLFiles: []string{"seed.sql"},
			},
		},
	}

	err = executeInit(db, cfg, tmpDir)
	require.NoError(t, err)

	// Verify both inline SQL and file SQL were executed
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM test")
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestExecuteInit_NoInit(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	cfg := config.Config{
		Database: &config.DatabaseConfig{
			DSN: "sqlite::memory:",
		},
	}

	err = executeInit(db, cfg, "")
	require.NoError(t, err)
}

func TestExecuteInit_NilDatabase(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	cfg := config.Config{}

	err = executeInit(db, cfg, "")
	require.NoError(t, err)
}

func TestExecuteInit_InvalidSQL(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	cfg := config.Config{
		Database: &config.DatabaseConfig{
			DSN: "sqlite::memory:",
			Init: &config.DatabaseInit{
				SQL: "INVALID SQL SYNTAX",
			},
		},
	}

	err = executeInit(db, cfg, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "inline sql")
}

func TestExecuteInit_MissingSQLFile(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	cfg := config.Config{
		Database: &config.DatabaseConfig{
			DSN: "sqlite::memory:",
			Init: &config.DatabaseInit{
				SQLFiles: []string{"nonexistent.sql"},
			},
		},
	}

	err = executeInit(db, cfg, "/tmp")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "sql_files")
}
