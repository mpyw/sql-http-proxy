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

func TestGlobalHelpers(t *testing.T) {
	cfg, err := config.ParseFile("global_helpers_test.yaml")
	require.NoError(t, err)

	// Use NewServeMux to properly initialize global_helpers
	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("helper in pre-transform", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/pre?n=5", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, float64(10), result["input_doubled"]) // double(5) = 10
	})

	t.Run("helper in mock", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/mock", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, float64(20), result["value"]) // double(10) = 20
	})

	t.Run("helper in post-transform", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/post", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, "John Doe", result["fullName"])
	})

	t.Run("helper in post.each", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/each", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var result []map[string]any
		require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Len(t, result, 2)
		assert.Equal(t, float64(10), result[0]["doubled"]) // double(5) = 10
		assert.Equal(t, float64(20), result[1]["doubled"]) // double(10) = 20
	})
}
