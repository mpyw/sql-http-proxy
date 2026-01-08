package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestInvalidCSVValueParserEmpty(t *testing.T) {
	_, err := config.ParseFile("invalid_csv_value_parser_both_test.yaml")
	require.Error(t, err)
	require.Contains(t, err.Error(), "value_parser")
}
