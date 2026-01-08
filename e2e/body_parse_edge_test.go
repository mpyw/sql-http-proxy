package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding/japanese"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestBodyParseEdgeCases(t *testing.T) {
	cfg, err := config.ParseFile("body_parse_edge_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("invalid content-type format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/json-only", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "invalid;;content-type")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("form not accepted on json-only endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/json-only", strings.NewReader(`name=test`))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
		require.Contains(t, rec.Body.String(), "not accepted")
	})

	t.Run("json not accepted on form-only endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/form-only", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
		require.Contains(t, rec.Body.String(), "not accepted")
	})

	t.Run("missing content-type on form-only endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/form-only", strings.NewReader(`{}`))
		// No Content-Type header
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
		require.Contains(t, rec.Body.String(), "missing Content-Type")
	})

	t.Run("multipart without boundary", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/form-only", strings.NewReader(`data`))
		req.Header.Set("Content-Type", "multipart/form-data")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "boundary")
	})

	t.Run("multipart not accepted on json-only endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/json-only", strings.NewReader(`data`))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=abc123")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
	})

	t.Run("json with charset conversion", func(t *testing.T) {
		// Encode JSON in Shift_JIS
		encoder := japanese.ShiftJIS.NewEncoder()
		sjisJSON, err := encoder.String(`{"name":"日本語"}`)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/json-only", strings.NewReader(sjisJSON))
		req.Header.Set("Content-Type", "application/json; charset=shift_jis")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("json with unsupported charset", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/json-only", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json; charset=unknown-charset")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "charset")
	})
}
