package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestInvalidJSONLWithOne(t *testing.T) {
	_, err := config.ParseFile("invalid_jsonl_with_one_test.yaml")
	require.Error(t, err)
	// jsonl is not allowed for type: one (mockTransformOne only allows js/json)
	assert.Contains(t, err.Error(), "additional properties 'jsonl' not allowed")
}
