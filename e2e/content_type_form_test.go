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

func TestContentTypeForm(t *testing.T) {
	cfg := loadConfig(t, "content_type_form_test.yaml")
	handler := createMutationHandler(t, nil, cfg)

	t.Run("accepts application/x-www-form-urlencoded", func(t *testing.T) {
		body := `name=Alice&age=30`
		req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"name": "Alice", "age": "30"}`, w.Body.String())
	})

	t.Run("accepts multipart/form-data", func(t *testing.T) {
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		_ = writer.WriteField("name", "Bob")
		_ = writer.WriteField("age", "25")
		_ = writer.Close()

		req := httptest.NewRequest("POST", "/user", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"name": "Bob", "age": "25"}`, w.Body.String())
	})

	t.Run("rejects application/json when accepts is form only", func(t *testing.T) {
		body := `{"name": "Charlie", "age": 35}`
		req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
	})
}
