package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryMany(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "query_many_test.yaml")
	handler := createHandler(t, db, cfg)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[[]map[string]any](t, w.Body)
	assert.Len(t, result, 3)
}
