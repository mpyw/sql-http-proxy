package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestInvalidCSVValueParser(t *testing.T) {
	cfg, err := config.ParseFile("invalid_csv_value_parser_test.yaml")
	require.NoError(t, err)
	err = cfg.ValidateTransforms()
	require.Error(t, err)
	require.Contains(t, err.Error(), "csv.value_parser")
}
