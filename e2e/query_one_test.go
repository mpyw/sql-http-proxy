package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryOne(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "query_one_test.yaml")
	handler := createHandler(t, db, cfg)

	req := httptest.NewRequest("GET", "/user?id=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[map[string]any](t, w.Body)
	assert.Equal(t, "John", result["first_name"])
}

func TestQueryOneNotFound(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "query_one_test.yaml")
	handler := createHandler(t, db, cfg)

	req := httptest.NewRequest("GET", "/user?id=999", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
