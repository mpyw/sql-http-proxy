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

func TestMutationMock(t *testing.T) {
	cfg := loadConfig(t, "mutation_mock_test.yaml")
	handler := createMutationHandler(t, nil, cfg)

	body := `{"first_name": "Mock", "last_name": "User", "email": "mock@example.com"}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, float64(999), result["id"])
	assert.Equal(t, "Mock", result["first_name"])
	assert.Equal(t, "User", result["last_name"])
	assert.Equal(t, "mock@example.com", result["email"])
}
