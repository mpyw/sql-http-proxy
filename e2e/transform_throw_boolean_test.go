package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformThrowBoolean(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_throw_boolean_test.yaml")
	handler := createHandler(t, db, cfg)

	req := httptest.NewRequest("GET", "/user?id=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := decodeJSON[bool](t, w.Body)
	assert.Equal(t, false, result)
}
