package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryPostMethod(t *testing.T) {
	cfg := loadConfig(t, "query_post_method_test.yaml")
	handler := createHandler(t, nil, cfg)

	t.Run("accepts POST with JSON body", func(t *testing.T) {
		body := `{"id": 1, "name": "Alice"}`
		req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"id": 1, "name": "Alice"}`, w.Body.String())
	})

	t.Run("rejects GET when method is POST", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user?id=1&name=Bob", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
