package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestMutationOneLastInsertId(t *testing.T) {
	cfg, err := config.ParseFile("mutation_one_last_insert_id_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(`{"name": "Alice"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var result map[string]any
	err = json.NewDecoder(rec.Body).Decode(&result)
	require.NoError(t, err)

	// In mock mode, lastInsertId and rowsAffected are not available
	require.Equal(t, false, result["hasLastInsertId"])
	require.Equal(t, false, result["hasRowsAffected"])
	require.Equal(t, "Alice", result["name"])
}
