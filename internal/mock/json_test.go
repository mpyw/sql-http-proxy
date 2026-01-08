package mock

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJSON(t *testing.T) {
	t.Run("array data", func(t *testing.T) {
		data := []any{
			map[string]any{"id": 1, "name": "Alice"},
			map[string]any{"id": 2, "name": "Bob"},
		}
		source, err := NewJSON(data)
		require.NoError(t, err)

		result, ctx, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		assert.Nil(t, ctx)

		// Result should be a deep copy
		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)

		row1 := resultSlice[0].(map[string]any)
		assert.Equal(t, float64(1), row1["id"]) // JSON unmarshal converts int to float64
		assert.Equal(t, "Alice", row1["name"])
	})

	t.Run("object data", func(t *testing.T) {
		data := map[string]any{"id": 1, "name": "Alice"}
		source, err := NewJSON(data)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultMap := result.(map[string]any)
		assert.Equal(t, float64(1), resultMap["id"])
		assert.Equal(t, "Alice", resultMap["name"])
	})

	t.Run("nil data", func(t *testing.T) {
		source, err := NewJSON(nil)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("string data (JSON array)", func(t *testing.T) {
		data := `[{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}]`
		source, err := NewJSON(data)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)

		row1 := resultSlice[0].(map[string]any)
		assert.Equal(t, float64(1), row1["id"])
		assert.Equal(t, "Alice", row1["name"])
	})

	t.Run("string data (JSON object)", func(t *testing.T) {
		data := `{"id": 1, "name": "Alice"}`
		source, err := NewJSON(data)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultMap := result.(map[string]any)
		assert.Equal(t, float64(1), resultMap["id"])
		assert.Equal(t, "Alice", resultMap["name"])
	})

	t.Run("invalid JSON string", func(t *testing.T) {
		data := `{invalid json}`
		_, err := NewJSON(data)
		require.Error(t, err)
	})

	t.Run("data returns deep copy", func(t *testing.T) {
		data := []any{map[string]any{"name": "Alice"}}
		source, err := NewJSON(data)
		require.NoError(t, err)

		result1, _ := lo.Must2(source.Data(nil, "", nil, nil))
		result2, _ := lo.Must2(source.Data(nil, "", nil, nil))

		// Modify first result
		slice1 := result1.([]any)
		row1 := slice1[0].(map[string]any)
		row1["name"] = "Modified"

		// Second result should be unchanged
		slice2 := result2.([]any)
		row2 := slice2[0].(map[string]any)
		assert.Equal(t, "Alice", row2["name"])
	})
}

func TestParseJSONFile(t *testing.T) {
	t.Run("valid JSON array file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.json")
		content := `[{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}]`
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)

		source, err := ParseJSONFile(path)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)
	})

	t.Run("valid JSON object file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.json")
		content := `{"id": 1, "name": "Alice"}`
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)

		source, err := ParseJSONFile(path)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultMap := result.(map[string]any)
		assert.Equal(t, float64(1), resultMap["id"])
	})

	t.Run("invalid JSON file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.json")
		content := `{invalid json}`
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)

		_, err = ParseJSONFile(path)
		require.Error(t, err)
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := ParseJSONFile("/nonexistent/file.json")
		require.Error(t, err)
	})
}

func TestParseJSONL(t *testing.T) {
	t.Run("valid JSONL", func(t *testing.T) {
		data := `{"id": 1, "name": "Alice"}
{"id": 2, "name": "Bob"}`
		source, err := ParseJSONL(data)
		require.NoError(t, err)

		result, ctx, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		assert.Nil(t, ctx)

		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)

		row1 := resultSlice[0].(map[string]any)
		assert.Equal(t, float64(1), row1["id"])
		assert.Equal(t, "Alice", row1["name"])

		row2 := resultSlice[1].(map[string]any)
		assert.Equal(t, float64(2), row2["id"])
		assert.Equal(t, "Bob", row2["name"])
	})

	t.Run("JSONL with empty lines", func(t *testing.T) {
		data := `{"id": 1}

{"id": 2}
`
		source, err := ParseJSONL(data)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2) // Empty lines skipped
	})

	t.Run("empty JSONL", func(t *testing.T) {
		source, err := ParseJSONL("")
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		// Empty JSONL returns nil (via deepCopy of nil []any)
		assert.Nil(t, result)
	})

	t.Run("invalid JSONL line", func(t *testing.T) {
		data := `{"id": 1}
{invalid json}
{"id": 2}`
		_, err := ParseJSONL(data)
		require.Error(t, err)
	})

	t.Run("JSONL with whitespace lines", func(t *testing.T) {
		data := `{"id": 1}

{"id": 2}`
		source, err := ParseJSONL(data)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)
	})
}

func TestParseJSONLFile(t *testing.T) {
	t.Run("valid JSONL file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.jsonl")
		content := `{"id": 1, "name": "Alice"}
{"id": 2, "name": "Bob"}`
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)

		source, err := ParseJSONLFile(path)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := ParseJSONLFile("/nonexistent/file.jsonl")
		require.Error(t, err)
	})
}
