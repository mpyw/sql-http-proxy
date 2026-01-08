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

func TestMutationOneNullResult(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "mutation_one_null_result_test.yaml")
	handler := createMutationHandler(t, db, cfg)

	t.Run("returns data when row exists", func(t *testing.T) {
		body := `{"id": 1}`
		req := httptest.NewRequest("DELETE", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, float64(1), result["id"])
		assert.Equal(t, "John", result["first_name"])
	})

	t.Run("returns 204 when row does not exist", func(t *testing.T) {
		body := `{"id": 999}`
		req := httptest.NewRequest("DELETE", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
	})
}
