package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformThrowBoolean(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_throw_boolean_test.yaml")
	handler := createHandler(t, db, cfg)

	req := httptest.NewRequest("GET", "/user?id=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result bool
	require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, false, result)
}
