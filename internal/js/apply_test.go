package js

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyPre(t *testing.T) {
	t.Run("transforms input", func(t *testing.T) {
		transformer, err := CompilePre(`return { ...input, transformed: true }`)
		require.NoError(t, err)

		result, err := transformer.ApplyPre(
			map[string]any{},
			"SELECT * FROM users",
			map[string]any{"id": 1},
			nil,
		)
		require.NoError(t, err)

		assert.Equal(t, int64(1), result.Output["id"])
		assert.Equal(t, true, result.Output["transformed"])
	})

	t.Run("modifies SQL via free variable", func(t *testing.T) {
		transformer, err := CompilePre(`sql = sql + " WHERE id = :id"; return input`)
		require.NoError(t, err)

		result, err := transformer.ApplyPre(
			map[string]any{},
			"SELECT * FROM users",
			map[string]any{"id": 1},
			nil,
		)
		require.NoError(t, err)

		assert.Equal(t, "SELECT * FROM users WHERE id = :id", result.SQL)
	})

	t.Run("modifies ctx via free variable", func(t *testing.T) {
		transformer, err := CompilePre(`ctx.startTime = Date.now(); return input`)
		require.NoError(t, err)

		result, err := transformer.ApplyPre(
			map[string]any{},
			"SELECT * FROM users",
			map[string]any{},
			nil,
		)
		require.NoError(t, err)

		assert.Contains(t, result.Ctx, "startTime")
	})

	t.Run("must return object", func(t *testing.T) {
		transformer, err := CompilePre(`return 42`)
		require.NoError(t, err)

		_, err = transformer.ApplyPre(
			map[string]any{},
			"SELECT * FROM users",
			map[string]any{},
			nil,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must return an object")
	})

	t.Run("throws error", func(t *testing.T) {
		transformer, err := CompilePre(`throw { status: 400, body: "error" }`)
		require.NoError(t, err)

		_, err = transformer.ApplyPre(
			map[string]any{},
			"SELECT * FROM users",
			map[string]any{},
			nil,
		)
		require.Error(t, err)

		var transformErr *TransformError
		require.ErrorAs(t, err, &transformErr)
		assert.Equal(t, 400, transformErr.Status)
	})
}

func TestApplyMock(t *testing.T) {
	t.Run("returns mock data", func(t *testing.T) {
		transformer, err := CompilePre(`return { id: parseInt(input.id), name: "mock" }`)
		require.NoError(t, err)

		result, err := transformer.ApplyMock(
			map[string]any{},
			"SELECT * FROM users WHERE id = :id",
			map[string]any{"id": "42"},
			nil,
		)
		require.NoError(t, err)

		resultMap := result.Output.(map[string]any)
		assert.Equal(t, int64(42), resultMap["id"])
		assert.Equal(t, "mock", resultMap["name"])
	})

	t.Run("returns array", func(t *testing.T) {
		transformer, err := CompilePre(`return [{ id: 1 }, { id: 2 }]`)
		require.NoError(t, err)

		result, err := transformer.ApplyMock(
			map[string]any{},
			"SELECT * FROM users",
			map[string]any{},
			nil,
		)
		require.NoError(t, err)

		resultSlice := result.Output.([]any)
		assert.Len(t, resultSlice, 2)
	})

	t.Run("can read sql free variable", func(t *testing.T) {
		transformer, err := CompilePre(`return { sql: sql }`)
		require.NoError(t, err)

		result, err := transformer.ApplyMock(
			map[string]any{},
			"SELECT * FROM users",
			map[string]any{},
			nil,
		)
		require.NoError(t, err)

		resultMap := result.Output.(map[string]any)
		assert.Equal(t, "SELECT * FROM users", resultMap["sql"])
	})

	t.Run("can modify ctx", func(t *testing.T) {
		transformer, err := CompilePre(`ctx.mockTime = Date.now(); return { id: 1 }`)
		require.NoError(t, err)

		result, err := transformer.ApplyMock(
			map[string]any{},
			"SELECT * FROM users",
			map[string]any{},
			nil,
		)
		require.NoError(t, err)

		assert.Contains(t, result.Ctx, "mockTime")
	})
}

func TestApplyPost(t *testing.T) {
	t.Run("transforms output", func(t *testing.T) {
		transformer, err := CompilePost(`return { ...output, processed: true }`)
		require.NoError(t, err)

		result, err := transformer.ApplyPost(
			map[string]any{},
			map[string]any{"requestId": "123"},
			map[string]any{"id": 1, "name": "test"},
			nil,
		)
		require.NoError(t, err)

		resultMap := result.Output.(map[string]any)
		assert.Equal(t, int64(1), resultMap["id"])
		assert.Equal(t, true, resultMap["processed"])
	})

	t.Run("can access input", func(t *testing.T) {
		transformer, err := CompilePost(`return { ...output, requestId: input.requestId }`)
		require.NoError(t, err)

		result, err := transformer.ApplyPost(
			map[string]any{},
			map[string]any{"requestId": "abc"},
			map[string]any{"id": 1},
			nil,
		)
		require.NoError(t, err)

		resultMap := result.Output.(map[string]any)
		assert.Equal(t, "abc", resultMap["requestId"])
	})

	t.Run("ctx is free variable", func(t *testing.T) {
		transformer, err := CompilePost(`return { ...output, startTime: ctx.startTime }`)
		require.NoError(t, err)

		result, err := transformer.ApplyPost(
			map[string]any{"startTime": int64(12345)},
			map[string]any{},
			map[string]any{"id": 1},
			nil,
		)
		require.NoError(t, err)

		resultMap := result.Output.(map[string]any)
		assert.Equal(t, int64(12345), resultMap["startTime"])
	})

	t.Run("can modify ctx", func(t *testing.T) {
		transformer, err := CompilePost(`ctx.postTime = Date.now(); return output`)
		require.NoError(t, err)

		result, err := transformer.ApplyPost(
			map[string]any{},
			map[string]any{},
			map[string]any{"id": 1},
			nil,
		)
		require.NoError(t, err)

		assert.Contains(t, result.Ctx, "postTime")
	})
}

func TestApplyPostToEachRow(t *testing.T) {
	t.Run("transforms each row", func(t *testing.T) {
		transformer, err := CompilePost(`return { ...output, processed: true }`)
		require.NoError(t, err)

		rows := []map[string]any{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
		}

		results, newCtx, err := transformer.ApplyPostToEachRow(
			map[string]any{},
			map[string]any{},
			rows,
			nil,
		)
		require.NoError(t, err)
		assert.NotNil(t, newCtx)
		assert.Len(t, results, 2)

		row1 := results[0].(map[string]any)
		assert.Equal(t, true, row1["processed"])
	})

	t.Run("ctx is passed through rows", func(t *testing.T) {
		transformer, err := CompilePost(`ctx.count = (ctx.count || 0) + 1; return output`)
		require.NoError(t, err)

		rows := []map[string]any{
			{"id": 1},
			{"id": 2},
			{"id": 3},
		}

		_, newCtx, err := transformer.ApplyPostToEachRow(
			map[string]any{},
			map[string]any{},
			rows,
			nil,
		)
		require.NoError(t, err)

		// After 3 rows, count should be 3
		assert.Equal(t, int64(3), newCtx["count"])
	})
}

func TestApplyPostToAllRows(t *testing.T) {
	t.Run("transforms entire array", func(t *testing.T) {
		transformer, err := CompilePost(`return { total: output.length, items: output }`)
		require.NoError(t, err)

		rows := []map[string]any{
			{"id": 1},
			{"id": 2},
		}

		result, err := transformer.ApplyPostToAllRows(
			map[string]any{},
			map[string]any{},
			rows,
			nil,
		)
		require.NoError(t, err)

		resultMap := result.Output.(map[string]any)
		assert.Equal(t, int64(2), resultMap["total"])
	})
}

func TestJSTimeout(t *testing.T) {
	t.Run("infinite loop times out", func(t *testing.T) {
		// This test verifies that infinite loops are interrupted
		// Note: JSTimeout is 5 seconds by default
		transformer, err := CompilePre(`while(true) {} return input`)
		require.NoError(t, err)

		_, err = transformer.ApplyPre(
			map[string]any{},
			"SELECT 1",
			map[string]any{},
			nil,
		)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrJSTimeout), "expected timeout error, got: %v", err)
	})
}
