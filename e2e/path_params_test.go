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

func TestPathParameters(t *testing.T) {
	cfg, err := config.ParseFile("path_params_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("basic path parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"id":1,"name":"Alice"}`, rec.Body.String())
	})

	t.Run("path parameter not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("multiple path parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1/posts/10", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"user_id":1,"post_id":10,"title":"Post 10"}`, rec.Body.String())
	})

	t.Run("multiple path parameters user 2", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/2/posts/30", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"user_id":2,"post_id":30,"title":"Post 30"}`, rec.Body.String())
	})

	t.Run("path parameter takes priority over query parameter", func(t *testing.T) {
		// Path param is "path_value", query param tries to override with "query_value"
		req := httptest.NewRequest(http.MethodGet, "/priority/path_value?id=query_value", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		// Path parameter should win
		require.JSONEq(t, `{"received_id":"path_value"}`, rec.Body.String())
	})

	t.Run("path parameter in mutation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/42", strings.NewReader(`{"name":"Updated User"}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"id":42,"name":"Updated User","updated":true}`, rec.Body.String())
	})
}
