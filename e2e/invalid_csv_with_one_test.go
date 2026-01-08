package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestInvalidCSVWithOne(t *testing.T) {
	_, err := config.ParseFile("invalid_csv_with_one_test.yaml")
	require.Error(t, err)
	// csv requires filter for type: one
	assert.Contains(t, err.Error(), "type 'one' with 'csv' requires 'filter'")
}
