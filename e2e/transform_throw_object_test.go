package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformThrowObject(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_throw_object_test.yaml")
	handler := createHandler(t, db, cfg)

	req := httptest.NewRequest("GET", "/user?id=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := decodeJSON[map[string]any](t, w.Body)
	assert.Equal(t, "object error", result["message"])
	assert.Equal(t, "INVALID", result["code"])
}
