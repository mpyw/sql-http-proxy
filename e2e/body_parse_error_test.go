package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestBodyParseError(t *testing.T) {
	cfg, err := config.ParseFile("body_parse_error_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/json", strings.NewReader(`{invalid json`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "invalid JSON")
	})

	t.Run("invalid form encoding", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/json", strings.NewReader(`%ZZinvalid`))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "invalid form data")
	})

	t.Run("unsupported content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/json", strings.NewReader(`some data`))
		req.Header.Set("Content-Type", "application/xml")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})
}
