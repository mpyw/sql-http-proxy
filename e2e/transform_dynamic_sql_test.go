package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformDynamicSql(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_dynamic_sql_test.yaml")
	handler := createHandler(t, db, cfg)

	t.Run("no search query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/search", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result []map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Len(t, result, 3)
	})

	t.Run("search by keyword", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/search?q=john", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result []map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Len(t, result, 1)
		assert.Equal(t, "john@example.com", result[0]["email"])
	})

	t.Run("search with multiple keywords", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/search?q=jane+smith", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result []map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Len(t, result, 1)
		assert.Equal(t, "jane@example.com", result[0]["email"])
	})
}
