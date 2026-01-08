package e2e

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding/japanese"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestMultipartCharset(t *testing.T) {
	cfg, err := config.ParseFile("multipart_charset_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("Shift_JIS charset", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)

		// Encode value in Shift_JIS
		encoder := japanese.ShiftJIS.NewEncoder()
		sjisValue, err := encoder.String("日本語")
		require.NoError(t, err)

		_ = w.WriteField("name", sjisValue)
		_ = w.Close()

		req := httptest.NewRequest(http.MethodPost, "/multipart", &body)
		req.Header.Set("Content-Type", fmt.Sprintf("%s; charset=shift_jis", w.FormDataContentType()))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"name":"日本語"}`, rec.Body.String())
	})

	t.Run("unsupported charset", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		_ = w.WriteField("name", "test")
		_ = w.Close()

		req := httptest.NewRequest(http.MethodPost, "/multipart", &body)
		req.Header.Set("Content-Type", fmt.Sprintf("%s; charset=unknown-charset", w.FormDataContentType()))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "charset")
	})
}
