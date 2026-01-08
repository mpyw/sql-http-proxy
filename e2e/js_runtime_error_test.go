package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestJSRuntimeError(t *testing.T) {
	cfg, err := config.ParseFile("js_runtime_error_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("pre transform error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pre-error", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "pre error")
	})

	t.Run("mock transform error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/mock-error", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "mock error")
	})

	t.Run("post transform error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/post-error", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "post error")
	})
}
