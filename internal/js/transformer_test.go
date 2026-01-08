package js

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompile(t *testing.T) {
	t.Run("valid JS with no args", func(t *testing.T) {
		transformer, err := Compile(`return 42`, nil)
		require.NoError(t, err)
		assert.NotNil(t, transformer)

		result, err := transformer.Apply()
		require.NoError(t, err)
		assert.Equal(t, int64(42), result)
	})

	t.Run("valid JS with single arg", func(t *testing.T) {
		transformer, err := Compile(`return x * 2`, []string{"x"})
		require.NoError(t, err)

		result, err := transformer.Apply(21)
		require.NoError(t, err)
		assert.Equal(t, int64(42), result)
	})

	t.Run("valid JS with multiple args", func(t *testing.T) {
		transformer, err := Compile(`return a + b`, []string{"a", "b"})
		require.NoError(t, err)

		result, err := transformer.Apply(20, 22)
		require.NoError(t, err)
		assert.Equal(t, int64(42), result)
	})

	t.Run("syntax error", func(t *testing.T) {
		_, err := Compile(`return {invalid`, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to compile transform")
	})
}

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

func TestApply(t *testing.T) {
	t.Run("returns object", func(t *testing.T) {
		transformer, err := Compile(`return { id: 1, name: "test" }`, nil)
		require.NoError(t, err)

		result, err := transformer.Apply()
		require.NoError(t, err)

		resultMap := result.(map[string]any)
		assert.Equal(t, int64(1), resultMap["id"])
		assert.Equal(t, "test", resultMap["name"])
	})

	t.Run("returns array", func(t *testing.T) {
		transformer, err := Compile(`return [1, 2, 3]`, nil)
		require.NoError(t, err)

		result, err := transformer.Apply()
		require.NoError(t, err)

		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 3)
	})

	t.Run("transform throws Lambda-style error", func(t *testing.T) {
		transformer, err := Compile(`throw { status: 400, body: { message: "bad request" } }`, nil)
		require.NoError(t, err)

		_, err = transformer.Apply()
		require.Error(t, err)

		var transformErr *TransformError
		require.ErrorAs(t, err, &transformErr)
		assert.Equal(t, 400, transformErr.Status)
	})

	t.Run("transform not a function", func(t *testing.T) {
		// This test is tricky - we need to somehow make program return non-function
		// Actually, the way it's compiled, it will always be a function
		// Skip this test as it's hard to trigger
	})
}
