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

func TestFormCharset(t *testing.T) {
	cfg, err := config.ParseFile("form_charset_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("UTF-8 charset explicit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/form", strings.NewReader("name=test"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"name":"test"}`, rec.Body.String())
	})

	t.Run("Shift_JIS charset", func(t *testing.T) {
		// Encode "日本語" in Shift_JIS
		encoder := japanese.ShiftJIS.NewEncoder()
		sjisName, err := encoder.String("日本語")
		require.NoError(t, err)

		body := "name=" + sjisName
		req := httptest.NewRequest(http.MethodPost, "/form", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=shift_jis")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"name":"日本語"}`, rec.Body.String())
	})

	t.Run("unsupported charset", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/form", strings.NewReader("name=test"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=unsupported-charset")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "charset")
	})
}
