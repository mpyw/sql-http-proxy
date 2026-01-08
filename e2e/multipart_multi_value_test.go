package e2e

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestMultipartMultiValue(t *testing.T) {
	cfg, err := config.ParseFile("multipart_multi_value_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("multiple values for same field", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)

		_ = w.WriteField("name", "test")
		_ = w.WriteField("tags", "one")
		_ = w.WriteField("tags", "two")
		_ = w.WriteField("tags", "three")
		_ = w.Close()

		req := httptest.NewRequest(http.MethodPost, "/multipart", &body)
		req.Header.Set("Content-Type", w.FormDataContentType())
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"tags":["one","two","three"],"name":"test"}`, rec.Body.String())
	})

	t.Run("file upload is skipped", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)

		_ = w.WriteField("name", "test")
		// Create a file field - should be skipped
		fw, err := w.CreateFormFile("document", "test.txt")
		require.NoError(t, err)
		_, _ = fw.Write([]byte("file content"))
		_ = w.Close()

		req := httptest.NewRequest(http.MethodPost, "/multipart", &body)
		req.Header.Set("Content-Type", w.FormDataContentType())
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		// File should be skipped, only name should be present
		// tags is undefined in input, so it becomes null
		require.JSONEq(t, `{"name":"test","tags":null}`, rec.Body.String())
	})
}
