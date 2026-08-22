package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformThrowArray(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_throw_array_test.yaml")
	handler := createHandler(t, db, cfg)

	req := httptest.NewRequest("GET", "/user?id=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	result := decodeJSON[[]map[string]any](t, w.Body)
	assert.Len(t, result, 2)
	assert.Equal(t, "id", result[0]["field"])
	assert.Equal(t, "required", result[0]["message"])
}
