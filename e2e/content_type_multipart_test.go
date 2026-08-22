package e2e

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestContentTypeMultipart(t *testing.T) {
	cfg, err := config.ParseFile("content_type_multipart_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("multipart form data", func(t *testing.T) {
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		_ = w.WriteField("id", "42")
		_ = w.WriteField("name", "Multipart User")
		_ = w.Close()

		req := httptest.NewRequest("POST", "/user", &buf)
		req.Header.Set("Content-Type", w.FormDataContentType())
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		result := decodeJSON[map[string]any](t, rec.Body)
		assert.Equal(t, float64(42), result["id"])
		assert.Equal(t, "Multipart User", result["name"])
	})
}
