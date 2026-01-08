package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestInvalidGlobalHelpers(t *testing.T) {
	cfg, err := config.ParseFile("invalid_global_helpers_test.yaml")
	if err != nil {
		// Error during parsing (e.g., JS compilation in ParseFile)
		require.Contains(t, err.Error(), "global_helpers")
		return
	}

	_, err = server.NewServeMux(nil, cfg, ".")
	require.Error(t, err)
	require.Contains(t, err.Error(), "global_helpers")
}
