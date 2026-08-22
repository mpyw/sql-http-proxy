package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformPostError(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_post_error_test.yaml")
	handler := createHandler(t, db, cfg)

	t.Run("success for allowed user", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user?id=2", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		result := decodeJSON[map[string]any](t, w.Body)
		assert.Equal(t, float64(2), result["id"])
	})

	t.Run("error for john", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user?id=1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		result := decodeJSON[map[string]any](t, w.Body)
		assert.Equal(t, "access denied", result["message"])
	})
}
