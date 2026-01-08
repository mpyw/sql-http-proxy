package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestMutationNoneWithOutput(t *testing.T) {
	db := setupTestDB(t)
	cfg, err := config.ParseFile("mutation_none_with_output_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(db, cfg, ".")
	require.NoError(t, err)

	// type: none returns 204 No Content regardless of SQL result
	req := httptest.NewRequest(http.MethodPost, "/execute", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNoContent, rec.Code)
}
