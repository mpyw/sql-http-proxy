package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutationNone(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "mutation_none_test.yaml")
	handler := createMutationHandler(t, db, cfg)

	body := `{"id": 1}`
	req := httptest.NewRequest("DELETE", "/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())

	// Verify the user was actually deleted
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE id = 1")
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
