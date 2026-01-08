package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	var result map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "Alice Brown", result["fullName"])
	assert.Equal(t, "alice@example.com", result["email"])
	assert.Equal(t, true, result["created"])
	// Original fields should not be present
	assert.Nil(t, result["first_name"])
	assert.Nil(t, result["last_name"])
}
