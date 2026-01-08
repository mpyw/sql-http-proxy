package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockOne(t *testing.T) {
	cfg := loadConfig(t, "mock_one_test.yaml")
	handler := createHandler(t, nil, cfg)

	req := httptest.NewRequest("GET", "/user?id=42", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, float64(42), result["id"])
	assert.Equal(t, "Mock User", result["name"])
	assert.Equal(t, "mock@example.com", result["email"])
}
