package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutationMockNone(t *testing.T) {
	cfg := loadConfig(t, "mutation_mock_none_test.yaml")
	handler := createMutationHandler(t, nil, cfg)

	body := `{"id": 1}`
	req := httptest.NewRequest("DELETE", "/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}
