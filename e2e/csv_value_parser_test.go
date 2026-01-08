package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestCSVValueParser(t *testing.T) {
	cfg, err := config.ParseFile("csv_value_parser_test.yaml")
	require.NoError(t, err)

	// Use NewServeMux to properly initialize csv.value_parser and global_helpers
	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/data", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&result))
	require.Len(t, result, 2)

	// Row 1: 1,Alice,true,95.5,2024-01-15,
	row1 := result[0]
	assert.Equal(t, float64(1), row1["id"]) // int parsed as number
	assert.Equal(t, "Alice", row1["name"])  // string unchanged
	assert.Equal(t, true, row1["active"])   // boolean true
	assert.Equal(t, 95.5, row1["score"])    // float parsed as number
	assert.Nil(t, row1["empty"])            // empty string -> null

	// Check custom date parsing (returns object with type and value)
	date1 := row1["date"].(map[string]any)
	assert.Equal(t, "date", date1["type"])
	assert.Equal(t, "2024-01-15", date1["value"])

	// Row 2: 2,Bob,false,87,2024-02-20,null
	row2 := result[1]
	assert.Equal(t, float64(2), row2["id"])     // int parsed as number
	assert.Equal(t, "Bob", row2["name"])        // string unchanged
	assert.Equal(t, false, row2["active"])      // boolean false
	assert.Equal(t, float64(87), row2["score"]) // int parsed as number
	assert.Nil(t, row2["empty"])                // "null" string -> null

	// Check custom date parsing
	date2 := row2["date"].(map[string]any)
	assert.Equal(t, "date", date2["type"])
	assert.Equal(t, "2024-02-20", date2["value"])
}
