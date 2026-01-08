package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestInvalidInvalidType(t *testing.T) {
	_, err := config.ParseFile("invalid_invalid_type_test.yaml")
	require.Error(t, err)
	// Schema now uses oneOf with const, so error shows both options
	assert.Contains(t, err.Error(), "value must be")
}
