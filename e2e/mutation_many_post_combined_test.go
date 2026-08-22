package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutationManyPostCombined(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "mutation_many_post_combined_test.yaml")
	handler := createMutationHandler(t, db, cfg)

	body := `{"max_id": 2}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[map[string]any](t, w.Body)
	assert.Equal(t, float64(2), result["total"])

	items := result["items"].([]any)
	assert.Len(t, items, 2)
	// DB uppercased, post.each lowercased
	assert.Equal(t, "john@example.com", items[0].(map[string]any)["email"])
}
