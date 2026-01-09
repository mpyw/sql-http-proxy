package e2e

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/db"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestEmptyBody_AcceptsEmpty(t *testing.T) {
	cfg := loadConfig(t, "empty_body_test.yaml")

	configDir, _ := filepath.Abs(".")
	conn, err := db.Connect(cfg, configDir)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	mux, err := server.NewServeMux(conn, cfg, configDir)
	require.NoError(t, err)

	t.Run("DELETE with empty body succeeds", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/resources/1", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		require.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("DELETE with non-empty body rejected when accepts=[]", func(t *testing.T) {
		// Non-empty body without Content-Type should fail with "missing Content-Type"
		req := httptest.NewRequest(http.MethodDelete, "/resources/1", strings.NewReader(`{"key":"value"}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		// Should reject because body is non-empty but no Content-Type accepted
		require.Equal(t, http.StatusUnsupportedMediaType, w.Code)
	})

	t.Run("DELETE with Content-Type rejected when accepts=[]", func(t *testing.T) {
		// Any Content-Type should fail when accepts=[]
		req := httptest.NewRequest(http.MethodDelete, "/resources/1", strings.NewReader(`{"key":"value"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnsupportedMediaType, w.Code)
	})

	t.Run("DELETE with body accepted when using default accepts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/resources-default/1", strings.NewReader(`{"key":"value"}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		// Should succeed because default accepts allows json
		require.Equal(t, http.StatusNoContent, w.Code)
	})
}
