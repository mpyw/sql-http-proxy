package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockManyPostEach(t *testing.T) {
	cfg := loadConfig(t, "mock_many_post_each_test.yaml")
	handler := createHandler(t, nil, cfg)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[[]map[string]any](t, w.Body)
	assert.Len(t, result, 2)
	assert.Equal(t, "ALICE", result[0]["name"])
	assert.Equal(t, "BOB", result[1]["name"])
}
