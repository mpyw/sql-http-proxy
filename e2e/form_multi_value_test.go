package e2e

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestFormMultiValue(t *testing.T) {
	cfg, err := config.ParseFile("form_multi_value_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("single value", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "test")
		form.Set("tags", "one")

		req := httptest.NewRequest(http.MethodPost, "/form", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"tags":"one","name":"test"}`, rec.Body.String())
	})

	t.Run("multiple values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "test")
		form.Add("tags", "one")
		form.Add("tags", "two")
		form.Add("tags", "three")

		req := httptest.NewRequest(http.MethodPost, "/form", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"tags":["one","two","three"],"name":"test"}`, rec.Body.String())
	})
}
