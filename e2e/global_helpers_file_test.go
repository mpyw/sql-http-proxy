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

func TestGlobalHelpersFile(t *testing.T) {
	cfg, err := config.ParseFile("global_helpers_file_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("uses both inline and file helpers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/calc?n=5&name=Test", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, float64(10), result["doubled"])     // double(5) from inline
		assert.Equal(t, float64(15), result["tripled"])     // triple(5) from file
		assert.Equal(t, "Hello, Test!", result["greeting"]) // greet() from file
	})
}
