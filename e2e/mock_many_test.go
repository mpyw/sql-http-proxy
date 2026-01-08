package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockMany(t *testing.T) {
	cfg := loadConfig(t, "mock_many_test.yaml")
	handler := createHandler(t, nil, cfg)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Len(t, result, 3)
	assert.Equal(t, float64(1), result[0]["id"])
	assert.Equal(t, "User 1", result[0]["name"])
}
