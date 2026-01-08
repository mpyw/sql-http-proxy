package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestInvalidInvalidPath(t *testing.T) {
	_, err := config.ParseFile("invalid_invalid_path_test.yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not match pattern")
}
