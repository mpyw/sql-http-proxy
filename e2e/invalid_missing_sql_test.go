package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestInvalidMissingSql(t *testing.T) {
	_, err := config.ParseFile("invalid_missing_sql_test.yaml")
	require.Error(t, err)
	// Improved error message for sql/mock exclusivity
	assert.Contains(t, err.Error(), "must have either 'sql' or 'mock'")
}
