package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutationPreTransform(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "mutation_pre_transform_test.yaml")
	handler := createMutationHandler(t, db, cfg)

	// Test successful pre-transform (email lowercased, names trimmed)
	body := `{"first_name": "  Alice  ", "last_name": "  Brown  ", "email": "ALICE@Example.COM"}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[map[string]any](t, w.Body)
	assert.Equal(t, "Alice", result["first_name"])
	assert.Equal(t, "Brown", result["last_name"])
	assert.Equal(t, "alice@example.com", result["email"])
}

func TestMutationPreTransformValidationError(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "mutation_pre_transform_test.yaml")
	handler := createMutationHandler(t, db, cfg)

	// Test validation error (invalid email)
	body := `{"first_name": "Alice", "last_name": "Brown", "email": "invalid-email"}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := decodeJSON[map[string]any](t, w.Body)
	assert.Equal(t, "invalid email", result["message"])
}
