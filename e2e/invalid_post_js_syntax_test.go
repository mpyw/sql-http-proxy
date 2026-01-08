package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestInvalidPostJSSyntax(t *testing.T) {
	cfg, err := config.ParseFile("invalid_post_js_syntax_test.yaml")
	require.NoError(t, err)

	_, err = server.NewServeMux(nil, cfg, ".")
	require.Error(t, err)
	require.Contains(t, err.Error(), "post")
}
