package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestInvalidHandleNotFoundWithMany(t *testing.T) {
	_, err := config.ParseFile("invalid_handle_not_found_with_many_test.yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "additional properties 'handle_not_found' not allowed")
}
