package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutationMockMany(t *testing.T) {
	cfg := loadConfig(t, "mutation_mock_many_test.yaml")
	handler := createMutationHandler(t, nil, cfg)

	body := `{}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[[]map[string]any](t, w.Body)
	assert.Len(t, result, 2)
	assert.Equal(t, float64(1), result[0]["id"])
	assert.Equal(t, "mock1@example.com", result[0]["email"])
	assert.Equal(t, float64(2), result[1]["id"])
	assert.Equal(t, "mock2@example.com", result[1]["email"])
}
