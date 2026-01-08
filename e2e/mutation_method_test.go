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

func TestMutationMethod(t *testing.T) {
	cfg, err := config.ParseFile("mutation_method_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("default POST method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/post-default", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"default_post":true}`, rec.Body.String())
	})

	t.Run("explicit POST method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/post-explicit", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"explicit_post":true}`, rec.Body.String())
	})

	t.Run("PUT method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/put", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"success":true}`, rec.Body.String())
	})

	t.Run("DELETE method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/delete", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"deleted":true}`, rec.Body.String())
	})

	t.Run("PATCH method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/patch", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"patched":true}`, rec.Body.String())
	})

	t.Run("wrong method returns 405", func(t *testing.T) {
		// GET request to PUT-only endpoint returns 405 (Method Not Allowed)
		req := httptest.NewRequest(http.MethodGet, "/put", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})
}
