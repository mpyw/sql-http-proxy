package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestExampleFileValidation(t *testing.T) {
	// Ensure the example file is valid and can be parsed
	cfg, err := config.ParseFile("../sql-http-proxy.example.yaml")
	require.NoError(t, err)

	// Basic sanity checks
	require.NotNil(t, cfg.GlobalHelpers)
	require.NotEmpty(t, cfg.GlobalHelpers.JS)

	require.NotNil(t, cfg.CSV)
	require.NotEmpty(t, cfg.CSV.ValueParser)

	require.NotEmpty(t, cfg.Queries)
	require.NotEmpty(t, cfg.Mutations)
}
