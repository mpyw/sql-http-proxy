package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutationManyPostAll(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "mutation_many_post_all_test.yaml")
	handler := createMutationHandler(t, db, cfg)

	body := `{"max_id": 2}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[[]map[string]any](t, w.Body)
	assert.Len(t, result, 2)
	// DB uppercased, then post.all lowercased
	assert.Equal(t, "john@example.com", result[0]["email"])
	assert.Equal(t, "jane@example.com", result[1]["email"])
	// Only id and email should be present (transformed)
	assert.Nil(t, result[0]["first_name"])
}
