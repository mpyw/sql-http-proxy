package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestMockMixed(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "mock_mixed_test.yaml")

	t.Run("real database query", func(t *testing.T) {
		handler, err := server.CreateHandler(db, cfg.Queries[0])
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/user?id=1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, float64(1), result["id"])
		assert.Equal(t, "John", result["first_name"])
		assert.NotContains(t, result, "source")
	})

	t.Run("mock query", func(t *testing.T) {
		handler, err := server.CreateHandler(db, cfg.Queries[1])
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/mock/user?id=1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, float64(1), result["id"])
		assert.Equal(t, "Mock User", result["name"])
		assert.Equal(t, "mock", result["source"])
	})
}
