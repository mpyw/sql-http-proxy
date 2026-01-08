package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformThrowNativeErrorPost(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_throw_native_error_post_test.yaml")
	handler := createHandler(t, db, cfg)

	// Native Error in post-transform should return 500
	req := httptest.NewRequest("GET", "/post-error?id=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, "post error")
}
