package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformThrowRawStringPre(t *testing.T) {
	db := setupTestDB(t)

	cfg := loadConfig(t, "transform_throw_raw_string_pre_test.yaml")
	handler := createHandler(t, db, cfg)

	// Raw string throw in pre-transform should return default 400
	req := httptest.NewRequest("GET", "/string-error?id=1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result string
	require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "string error", result)
}
