package js

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompilePre(t *testing.T) {
	t.Run("valid pre-transform", func(t *testing.T) {
		transformer, err := CompilePre(`return { id: input.id * 2 }`)
		require.NoError(t, err)
		assert.NotNil(t, transformer)
	})

	t.Run("syntax error", func(t *testing.T) {
		_, err := CompilePre(`return {invalid`)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to compile pre-transform")
	})
}

func TestCompilePost(t *testing.T) {
	t.Run("valid post-transform", func(t *testing.T) {
		transformer, err := CompilePost(`return { ...output, processed: true }`)
		require.NoError(t, err)
		assert.NotNil(t, transformer)
	})

	t.Run("syntax error", func(t *testing.T) {
		_, err := CompilePost(`return {invalid`)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to compile post-transform")
	})
}
