package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestInvalidMissingType(t *testing.T) {
	_, err := config.ParseFile("invalid_missing_type_test.yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing property 'type'")
}
