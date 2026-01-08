package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompileJS(t *testing.T) {
	t.Run("valid JS returning object", func(t *testing.T) {
		code := `return { id: parseInt(input.id), name: "Mock User" }`
		source, err := CompileJS(code)
		require.NoError(t, err)

		result, ctx, err := source.Data(
			map[string]any{},
			"SELECT * FROM users",
			map[string]any{"id": "42"},
		)
		require.NoError(t, err)
		assert.NotNil(t, ctx)

		resultMap := result.(map[string]any)
		assert.Equal(t, int64(42), resultMap["id"])
		assert.Equal(t, "Mock User", resultMap["name"])
	})

	t.Run("valid JS returning array", func(t *testing.T) {
		code := `return [{ id: 1 }, { id: 2 }]`
		source, err := CompileJS(code)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil)
		require.NoError(t, err)

		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)
	})

	t.Run("JS with ctx modification", func(t *testing.T) {
		code := `ctx.mockTime = Date.now(); return { id: 1 }`
		source, err := CompileJS(code)
		require.NoError(t, err)

		ctx := map[string]any{}
		_, newCtx, err := source.Data(ctx, "", nil)
		require.NoError(t, err)

		assert.Contains(t, newCtx, "mockTime")
	})

	t.Run("JS with sql access", func(t *testing.T) {
		code := `return { sql: sql }`
		source, err := CompileJS(code)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "SELECT * FROM users", nil)
		require.NoError(t, err)

		resultMap := result.(map[string]any)
		assert.Equal(t, "SELECT * FROM users", resultMap["sql"])
	})

	t.Run("invalid JS syntax", func(t *testing.T) {
		code := `return { invalid syntax`
		_, err := CompileJS(code)
		require.Error(t, err)
	})

	t.Run("JS throwing error", func(t *testing.T) {
		code := `throw { status: 400, body: { message: "error" } }`
		source, err := CompileJS(code)
		require.NoError(t, err)

		_, _, err = source.Data(nil, "", nil)
		require.Error(t, err)
	})
}

