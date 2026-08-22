package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformManyPostCombined(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_many_post_combined_test.yaml")
	handler := createHandler(t, db, cfg)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[map[string]any](t, w.Body)
	assert.Equal(t, float64(2), result["total"])

	items := result["items"].([]any)
	assert.Len(t, items, 2)
	assert.Equal(t, "JOHN@EXAMPLE.COM", items[0].(map[string]any)["email"])
}
