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

func TestQueryPostMethod(t *testing.T) {
	cfg, err := config.ParseFile("query_post_method_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("POST query with body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/search", strings.NewReader(`{"query":"test"}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"query":"test","result":{"id":1}}`, rec.Body.String())
	})

	t.Run("POST query with form data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/search", strings.NewReader(`query=form-test`))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"query":"form-test","result":{"id":1}}`, rec.Body.String())
	})

	t.Run("GET query wrong method returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/get-only", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("POST query wrong method returns 404", func(t *testing.T) {
		// GET request to POST-only query endpoint
		req := httptest.NewRequest(http.MethodGet, "/search", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})
}
