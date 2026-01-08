package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentTypeJSON(t *testing.T) {
	cfg := loadConfig(t, "content_type_json_test.yaml")
	handler := createMutationHandler(t, nil, cfg)

	t.Run("accepts application/json", func(t *testing.T) {
		body := `{"name": "Alice", "age": 30}`
		req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"name": "Alice", "age": 30}`, w.Body.String())
	})

	t.Run("accepts application/json with charset utf-8", func(t *testing.T) {
		body := `{"name": "Bob", "age": 25}`
		req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"name": "Bob", "age": 25}`, w.Body.String())
	})

	t.Run("rejects form-urlencoded when accepts is json only", func(t *testing.T) {
		body := `name=Charlie&age=35`
		req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
	})
}
