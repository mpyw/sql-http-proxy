package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformPreError(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_pre_error_test.yaml")
	handler := createHandler(t, db, cfg)

	t.Run("success with id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user?id=1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, float64(1), result["id"])
	})

	t.Run("error without id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, "id is required", result["message"])
	})
}
