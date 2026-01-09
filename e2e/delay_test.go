package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestDelay(t *testing.T) {
	cfg, err := config.ParseFile("delay_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("object with delay", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slow-user", nil)
		rec := httptest.NewRecorder()

		start := time.Now()
		mux.ServeHTTP(rec, req)
		elapsed := time.Since(start)

		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"id":1,"name":"Slow User"}`, rec.Body.String())
		// Should take at least 50ms
		require.GreaterOrEqual(t, elapsed.Milliseconds(), int64(40), "expected delay of at least 40ms")
	})

	t.Run("array with delay", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slow-users", nil)
		rec := httptest.NewRecorder()

		start := time.Now()
		mux.ServeHTTP(rec, req)
		elapsed := time.Since(start)

		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`, rec.Body.String())
		// Should take at least 50ms
		require.GreaterOrEqual(t, elapsed.Milliseconds(), int64(40), "expected delay of at least 40ms")
	})

	t.Run("filtered array with delay", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slow-filtered-user?id=1", nil)
		rec := httptest.NewRecorder()

		start := time.Now()
		mux.ServeHTTP(rec, req)
		elapsed := time.Since(start)

		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"id":1,"name":"Alice"}`, rec.Body.String())
		// Should take at least 50ms
		require.GreaterOrEqual(t, elapsed.Milliseconds(), int64(40), "expected delay of at least 40ms")
	})

	t.Run("JS source with delay", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slow-js-user?id=42", nil)
		rec := httptest.NewRecorder()

		start := time.Now()
		mux.ServeHTTP(rec, req)
		elapsed := time.Since(start)

		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"id":42,"name":"JS User"}`, rec.Body.String())
		// Should take at least 50ms
		require.GreaterOrEqual(t, elapsed.Milliseconds(), int64(40), "expected delay of at least 40ms")
	})

	t.Run("no delay is fast", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/no-delay-user", nil)
		rec := httptest.NewRecorder()

		start := time.Now()
		mux.ServeHTTP(rec, req)
		elapsed := time.Since(start)

		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"id":1,"name":"Fast User"}`, rec.Body.String())
		// Should complete quickly (less than 30ms typically)
		require.Less(t, elapsed.Milliseconds(), int64(30), "expected fast response without delay")
	})
}
