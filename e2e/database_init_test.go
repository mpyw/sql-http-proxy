package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/db"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestDatabaseInit_InlineSQL(t *testing.T) {
	cfg := loadConfig(t, "database_init_test.yaml")

	// Connect to database - this should execute init SQL
	configDir, _ := filepath.Abs(".")
	conn, err := db.Connect(cfg, configDir)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer func() { _ = conn.Close() }()

	// Create server handler
	mux, err := server.NewServeMux(conn, cfg, configDir)
	require.NoError(t, err)

	// Test that data was initialized
	req := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var result []map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&result))

	// Verify init SQL created table and inserted data
	require.Len(t, result, 2)
	assert.Equal(t, float64(1), result[0]["id"])
	assert.Equal(t, "Widget", result[0]["name"])
	assert.Equal(t, 19.99, result[0]["price"])
	assert.Equal(t, float64(2), result[1]["id"])
	assert.Equal(t, "Gadget", result[1]["name"])
	assert.Equal(t, 29.99, result[1]["price"])
}

func TestDatabaseInit_SQLFile(t *testing.T) {
	// Create test config and SQL file in temp directory
	tmpDir := t.TempDir()

	// Create SQL init file
	sqlFile := filepath.Join(tmpDir, "init.sql")
	err := os.WriteFile(sqlFile, []byte(`
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		);
		DELETE FROM categories;
		INSERT INTO categories (id, name) VALUES (1, 'Electronics');
		INSERT INTO categories (id, name) VALUES (2, 'Books');
	`), 0644)
	require.NoError(t, err)

	// Create config file
	configFile := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configFile, []byte(`
database:
  dsn: ":memory:"
  init:
    sql_files:
      - ./init.sql

queries:
  - type: many
    path: /categories
    sql: SELECT * FROM categories ORDER BY id
`), 0644)
	require.NoError(t, err)

	// Parse config
	cfg, err := config.ParseFile(configFile)
	require.NoError(t, err)

	// Connect to database - this should execute init SQL files
	conn, err := db.Connect(cfg, tmpDir)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer func() { _ = conn.Close() }()

	// Create server handler
	mux, err := server.NewServeMux(conn, cfg, tmpDir)
	require.NoError(t, err)

	// Test that data was initialized
	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var result []map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&result))

	// Verify init SQL file created table and inserted data
	require.Len(t, result, 2)
	assert.Equal(t, float64(1), result[0]["id"])
	assert.Equal(t, "Electronics", result[0]["name"])
	assert.Equal(t, float64(2), result[1]["id"])
	assert.Equal(t, "Books", result[1]["name"])
}
