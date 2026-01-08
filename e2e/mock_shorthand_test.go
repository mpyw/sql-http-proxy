package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestMockShorthand(t *testing.T) {
	cfg, err := config.ParseFile("mock_shorthand_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("array shorthand for type: many", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result []map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Len(t, result, 2)
		assert.Equal(t, float64(1), result[0]["id"])
		assert.Equal(t, "Alice", result[0]["name"])
		assert.Equal(t, float64(2), result[1]["id"])
		assert.Equal(t, "Bob", result[1]["name"])
	})

	t.Run("object shorthand for type: one", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, float64(42), result["id"])
		assert.Equal(t, "Charlie", result["name"])
	})
}
