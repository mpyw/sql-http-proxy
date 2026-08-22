package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutationPostTransform(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "mutation_post_transform_test.yaml")
	handler := createMutationHandler(t, db, cfg)

	body := `{"first_name": "Alice", "last_name": "Brown", "email": "alice@example.com"}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[map[string]any](t, w.Body)
	assert.Equal(t, "Alice Brown", result["fullName"])
	assert.Equal(t, "alice@example.com", result["email"])
	assert.Equal(t, true, result["created"])
	// Original fields should not be present
	assert.Nil(t, result["first_name"])
	assert.Nil(t, result["last_name"])
}
