package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutationCtxSharing(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "mutation_ctx_sharing_test.yaml")
	handler := createMutationHandler(t, db, cfg)

	body := `{"first_name": "Alice", "last_name": "Brown", "email": "alice@example.com"}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[map[string]any](t, w.Body)
	assert.Equal(t, "alice@example.com", result["inputEmail"])
}
