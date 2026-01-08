package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformThrowObjectNoStatusBody(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_throw_object_no_status_body_test.yaml")
	handler := createHandler(t, db, cfg)

	// Object without status/body should return 500 (treated like native Error)
	req := httptest.NewRequest("GET", "/obj-error?id=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
