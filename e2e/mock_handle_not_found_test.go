package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockHandleNotFound(t *testing.T) {
	cfg := loadConfig(t, "mock_handle_not_found_test.yaml")
	handler := createHandler(t, nil, cfg)

	t.Run("found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user?id=1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, true, result["found"])
		assert.Equal(t, float64(1), result["id"])
	})

	t.Run("not found with post transform", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user?id=999", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, false, result["found"])
		assert.Equal(t, "guest", result["default"])
	})
}
