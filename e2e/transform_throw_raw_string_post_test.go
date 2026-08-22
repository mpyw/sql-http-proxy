package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformThrowRawStringPost(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_throw_raw_string_post_test.yaml")
	handler := createHandler(t, db, cfg)

	// Raw string throw in post-transform should return default 500
	req := httptest.NewRequest("GET", "/string-error?id=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	result := decodeJSON[string](t, w.Body)
	assert.Equal(t, "string error", result)
}
