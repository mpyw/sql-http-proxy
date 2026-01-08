package e2e

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentTypeDefault(t *testing.T) {
	cfg := loadConfig(t, "content_type_default_test.yaml")
	handler := createMutationHandler(t, nil, cfg)

	t.Run("accepts application/json by default", func(t *testing.T) {
		body := `{"name": "Alice"}`
		req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"name": "Alice"}`, w.Body.String())
	})

	t.Run("accepts form-urlencoded by default", func(t *testing.T) {
		body := `name=Bob`
		req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"name": "Bob"}`, w.Body.String())
	})

	t.Run("accepts multipart/form-data by default", func(t *testing.T) {
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		_ = writer.WriteField("name", "Charlie")
		_ = writer.Close()

		req := httptest.NewRequest("POST", "/user", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"name": "Charlie"}`, w.Body.String())
	})

	t.Run("rejects unsupported media type", func(t *testing.T) {
		body := `<name>Dave</name>`
		req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/xml")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
	})
}
