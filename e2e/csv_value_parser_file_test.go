package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestCSVValueParserFile(t *testing.T) {
	cfg, err := config.ParseFile("csv_value_parser_file_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/data", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := decodeJSON[[]map[string]any](t, w.Body)
	require.Len(t, result, 2)

	// Custom parser uppercases strings
	assert.Equal(t, float64(1), result[0]["id"])
	assert.Equal(t, "ALICE", result[0]["name"]) // uppercased by custom parser
	assert.Equal(t, true, result[0]["active"])

	assert.Equal(t, float64(2), result[1]["id"])
	assert.Equal(t, "BOB", result[1]["name"]) // uppercased by custom parser
	assert.Equal(t, false, result[1]["active"])
}
