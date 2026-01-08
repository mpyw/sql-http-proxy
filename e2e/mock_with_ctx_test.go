package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockWithCtx(t *testing.T) {
	cfg := loadConfig(t, "mock_with_ctx_test.yaml")
	handler := createHandler(t, nil, cfg)

	req := httptest.NewRequest("GET", "/user?id=42", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, float64(42), result["id"])
	assert.Equal(t, "42", result["requestedId"])
	assert.Equal(t, true, result["hadMockTime"])
}
