package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMutationDynamicSql(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "mutation_dynamic_sql_test.yaml")
	handler := createMutationHandler(t, db, cfg)

	t.Run("update single field", func(t *testing.T) {
		body := `{"id": 1, "first_name": "Updated"}`
		req := httptest.NewRequest("PATCH", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, "Updated", result["first_name"])
		assert.Equal(t, "Doe", result["last_name"]) // unchanged
	})

	t.Run("update multiple fields", func(t *testing.T) {
		body := `{"id": 2, "first_name": "Janet", "email": "janet@example.com"}`
		req := httptest.NewRequest("PATCH", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, "Janet", result["first_name"])
		assert.Equal(t, "Smith", result["last_name"]) // unchanged
		assert.Equal(t, "janet@example.com", result["email"])
	})
}
