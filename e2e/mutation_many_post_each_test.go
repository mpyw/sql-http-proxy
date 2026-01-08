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

func TestMutationManyPostEach(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "mutation_many_post_each_test.yaml")
	handler := createMutationHandler(t, db, cfg)

	body := `{"max_id": 2}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Len(t, result, 2)

	// Check transformed fields
	assert.Equal(t, "John Doe", result[0]["fullName"])
	assert.Equal(t, "JOHN@EXAMPLE.COM", result[0]["emailUpper"])
	assert.Equal(t, "Jane Smith", result[1]["fullName"])
	assert.Equal(t, "JANE@EXAMPLE.COM", result[1]["emailUpper"])

	// Original fields should not be present
	assert.Nil(t, result[0]["first_name"])
	assert.Nil(t, result[0]["last_name"])
}
