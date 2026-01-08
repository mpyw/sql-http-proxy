package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestInvalidEachWithOne(t *testing.T) {
	_, err := config.ParseFile("invalid_each_with_one_test.yaml")
	require.Error(t, err)
	// For "one" type, post must be a string, not an object
	assert.Contains(t, err.Error(), "got object, want string")
}
