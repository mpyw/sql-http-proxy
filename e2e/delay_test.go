package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

// configuredDelay matches the `delay:` value on the delayed endpoints in
// delay_test.yaml.
const configuredDelay = 50 * time.Millisecond

func TestDelay(t *testing.T) {
	cfg, err := config.ParseFile("delay_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	// Each case runs inside a synctest bubble, so the handler's time.Sleep is
	// served by a fake clock: the assertions can require the exact delay
	// instead of a tolerance window, and the test does not actually wait.
	tests := []struct {
		name  string
		path  string
		body  string
		delay time.Duration
	}{
		{"object with delay", "/slow-user", `{"id":1,"name":"Slow User"}`, configuredDelay},
		{"array with delay", "/slow-users", `[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`, configuredDelay},
		{"filtered array with delay", "/slow-filtered-user?id=1", `{"id":1,"name":"Alice"}`, configuredDelay},
		{"JS source with delay", "/slow-js-user?id=42", `{"id":42,"name":"JS User"}`, configuredDelay},
		{"no delay", "/no-delay-user", `{"id":1,"name":"Fast User"}`, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, tt.path, nil)
				rec := httptest.NewRecorder()

				start := time.Now()
				mux.ServeHTTP(rec, req)
				elapsed := time.Since(start)

				require.Equal(t, http.StatusOK, rec.Code)
				require.JSONEq(t, tt.body, rec.Body.String())
				require.Equal(t, tt.delay, elapsed)
			})
		})
	}
}
