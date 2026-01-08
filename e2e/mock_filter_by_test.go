package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestMockFilterBy(t *testing.T) {
	cfg, err := config.ParseFile("mock_filter_by_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("type one - found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user?id=2", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var result map[string]any
		err := json.NewDecoder(rec.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, float64(2), result["id"])
		require.Equal(t, "Bob", result["name"])
	})

	t.Run("type one - not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user?id=999", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("type many - filter by role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users?role=admin", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var result []map[string]any
		err := json.NewDecoder(rec.Body).Decode(&result)
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, "Alice", result[0]["name"])
		require.Equal(t, "Charlie", result[1]["name"])
	})

	t.Run("type many - no matches returns empty array", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users?role=guest", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var result []map[string]any
		err := json.NewDecoder(rec.Body).Decode(&result)
		require.NoError(t, err)
		require.Len(t, result, 0)
	})

	t.Run("type many - without filter returns all", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/all-users", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var result []map[string]any
		err := json.NewDecoder(rec.Body).Decode(&result)
		require.NoError(t, err)
		require.Len(t, result, 2)
	})
}
