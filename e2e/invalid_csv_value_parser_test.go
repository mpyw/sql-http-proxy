package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestInvalidCSVValueParser(t *testing.T) {
	_, err := config.ParseFile("invalid_csv_value_parser_test.yaml")
	require.Error(t, err)
	require.Contains(t, err.Error(), "csv.value_parser")
}
