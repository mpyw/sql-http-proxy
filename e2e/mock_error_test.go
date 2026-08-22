package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockError(t *testing.T) {
	cfg := loadConfig(t, "mock_error_test.yaml")
	handler := createHandler(t, nil, cfg)

	t.Run("success with id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user?id=1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		result := decodeJSON[map[string]any](t, w.Body)
		assert.Equal(t, float64(1), result["id"])
	})

	t.Run("error without id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		result := decodeJSON[map[string]any](t, w.Body)
		assert.Equal(t, "id is required", result["message"])
	})
}
