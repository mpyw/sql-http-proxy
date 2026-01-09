package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidMockReturn(t *testing.T) {
	cfg := loadConfig(t, "invalid_mock_return_test.yaml")
	handler := createHandler(t, nil, cfg)

	t.Run("array_js returns scalar", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/array-returns-scalar", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "array")
	})

	t.Run("array_js returns object", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/array-returns-object", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "array")
	})

	t.Run("object_js returns array", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/object-returns-array", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		// Error message mentions type validation
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("object_js returns scalar", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/object-returns-scalar", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		// Error message mentions type validation
		assert.Contains(t, w.Body.String(), "error")
	})
}
