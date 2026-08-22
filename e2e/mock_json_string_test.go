package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestMockJSONString(t *testing.T) {
	cfg := loadConfig(t, "mock_json_string_test.yaml")

	t.Run("type many with JSON array string", func(t *testing.T) {
		handler, err := server.NewQueryHandler(nil, cfg.Queries[0])
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		result := decodeJSON[[]map[string]any](t, w.Body)
		assert.Len(t, result, 2)
		assert.Equal(t, float64(1), result[0]["id"])
		assert.Equal(t, "Alice", result[0]["name"])
		assert.Equal(t, float64(2), result[1]["id"])
		assert.Equal(t, "Bob", result[1]["name"])
	})

	t.Run("type one with JSON object string", func(t *testing.T) {
		handler, err := server.NewQueryHandler(nil, cfg.Queries[1])
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/user?id=1", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		result := decodeJSON[map[string]any](t, w.Body)
		assert.Equal(t, float64(42), result["id"])
		assert.Equal(t, "Charlie", result["name"])
	})
}
