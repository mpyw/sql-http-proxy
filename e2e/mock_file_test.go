package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestMockFile(t *testing.T) {
	cfg, err := config.ParseFile("mock_file_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("csv_file", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/csv", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		result := decodeJSON[[]map[string]any](t, w.Body)
		assert.Len(t, result, 2)
		assert.Equal(t, "CSV User 1", result[0]["name"])
		assert.Equal(t, true, result[0]["active"])
	})

	t.Run("json_file", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/json", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		result := decodeJSON[[]map[string]any](t, w.Body)
		assert.Len(t, result, 2)
		assert.Equal(t, "JSON User 1", result[0]["name"])
	})

	t.Run("jsonl_file", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/jsonl", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		result := decodeJSON[[]map[string]any](t, w.Body)
		assert.Len(t, result, 2)
		assert.Equal(t, "JSONL User 1", result[0]["name"])
	})
}
