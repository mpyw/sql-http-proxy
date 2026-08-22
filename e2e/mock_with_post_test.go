package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockWithPost(t *testing.T) {
	cfg := loadConfig(t, "mock_with_post_test.yaml")
	handler := createHandler(t, nil, cfg)

	req := httptest.NewRequest("GET", "/user?id=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[map[string]any](t, w.Body)
	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, "Mock User", result["fullName"])
}
