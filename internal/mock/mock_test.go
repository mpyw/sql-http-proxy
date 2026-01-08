package mock

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestCompile(t *testing.T) {
	t.Run("nil Mock", func(t *testing.T) {
		source, err := Compile(nil, CompileOptions{})
		require.NoError(t, err)
		assert.Nil(t, source)
	})

	t.Run("ObjectJS inline", func(t *testing.T) {
		m := &config.Mock{ObjectJS: `return { id: 1 }`}
		source, err := Compile(m, CompileOptions{})
		require.NoError(t, err)
		assert.NotNil(t, source)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultMap := result.(map[string]any)
		assert.Equal(t, int64(1), resultMap["id"])
	})

	t.Run("ArrayJS inline", func(t *testing.T) {
		m := &config.Mock{ArrayJS: `return [{ id: 1 }, { id: 2 }]`}
		source, err := Compile(m, CompileOptions{})
		require.NoError(t, err)
		assert.NotNil(t, source)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)
	})

	t.Run("CSV inline", func(t *testing.T) {
		m := &config.Mock{CSV: "id,name\n42,Alice"}
		source, err := Compile(m, CompileOptions{})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultSlice := result.([]map[string]any)
		assert.Len(t, resultSlice, 1)
		assert.Equal(t, float64(42), resultSlice[0]["id"])
	})

	t.Run("CSV file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "data.csv")
		err := os.WriteFile(path, []byte("id,name\n1,Alice"), 0644)
		require.NoError(t, err)

		m := &config.Mock{CSVFile: "data.csv"}
		source, err := Compile(m, CompileOptions{ConfigDir: dir})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultSlice := result.([]map[string]any)
		assert.Len(t, resultSlice, 1)
	})

	t.Run("Object inline", func(t *testing.T) {
		data := map[string]any{"id": 1}
		m := &config.Mock{Object: data}
		source, err := Compile(m, CompileOptions{})
		require.NoError(t, err)
		assert.NotNil(t, source)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultMap := result.(map[string]any)
		// Preserves original Go type
		assert.Equal(t, 1, resultMap["id"])
	})

	t.Run("Array inline", func(t *testing.T) {
		data := []any{map[string]any{"id": 1}}
		m := &config.Mock{Array: data}
		source, err := Compile(m, CompileOptions{})
		require.NoError(t, err)
		assert.NotNil(t, source)
	})

	t.Run("ObjectJSONFile", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "data.json")
		err := os.WriteFile(path, []byte(`{"id": 1}`), 0644)
		require.NoError(t, err)

		m := &config.Mock{ObjectJSONFile: "data.json"}
		source, err := Compile(m, CompileOptions{ConfigDir: dir})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultMap := result.(map[string]any)
		assert.Equal(t, float64(1), resultMap["id"])
	})

	t.Run("ArrayJSONFile", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "data.json")
		err := os.WriteFile(path, []byte(`[{"id": 1}]`), 0644)
		require.NoError(t, err)

		m := &config.Mock{ArrayJSONFile: "data.json"}
		source, err := Compile(m, CompileOptions{ConfigDir: dir})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 1)
	})

	t.Run("JSONL inline", func(t *testing.T) {
		m := &config.Mock{JSONL: `{"id": 1}
{"id": 2}`}
		source, err := Compile(m, CompileOptions{})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)
	})

	t.Run("JSONL file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "data.jsonl")
		err := os.WriteFile(path, []byte(`{"id": 1}
{"id": 2}`), 0644)
		require.NoError(t, err)

		m := &config.Mock{JSONLFile: "data.jsonl"}
		source, err := Compile(m, CompileOptions{ConfigDir: dir})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)
	})

	t.Run("empty Mock", func(t *testing.T) {
		m := &config.Mock{}
		_, err := Compile(m, CompileOptions{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no source specified")
	})

	t.Run("validation error - multiple sources", func(t *testing.T) {
		m := &config.Mock{
			ObjectJS: `return {}`,
			CSV:      "id\n42",
		}
		_, err := Compile(m, CompileOptions{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only one source type can be specified")
	})

	t.Run("JS compile error", func(t *testing.T) {
		m := &config.Mock{ObjectJS: `return { invalid`}
		_, err := Compile(m, CompileOptions{})
		require.Error(t, err)
	})

	t.Run("CSV parse error", func(t *testing.T) {
		m := &config.Mock{CSV: ""} // Empty CSV
		_, err := Compile(m, CompileOptions{})
		require.Error(t, err)
	})

	t.Run("CSV file not found", func(t *testing.T) {
		m := &config.Mock{CSVFile: "nonexistent.csv"}
		_, err := Compile(m, CompileOptions{ConfigDir: "/tmp"})
		require.Error(t, err)
	})

	t.Run("JSON file not found", func(t *testing.T) {
		m := &config.Mock{ArrayJSONFile: "nonexistent.json"}
		_, err := Compile(m, CompileOptions{ConfigDir: "/tmp"})
		require.Error(t, err)
	})

	t.Run("JSONL parse error", func(t *testing.T) {
		m := &config.Mock{JSONL: `{invalid json}`}
		_, err := Compile(m, CompileOptions{})
		require.Error(t, err)
	})

	t.Run("JSONL file not found", func(t *testing.T) {
		m := &config.Mock{JSONLFile: "nonexistent.jsonl"}
		_, err := Compile(m, CompileOptions{ConfigDir: "/tmp"})
		require.Error(t, err)
	})

	t.Run("ObjectJSON string", func(t *testing.T) {
		m := &config.Mock{ObjectJSON: `{"id": 1, "name": "Alice"}`}
		source, err := Compile(m, CompileOptions{})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultMap := result.(map[string]any)
		assert.Equal(t, float64(1), resultMap["id"])
		assert.Equal(t, "Alice", resultMap["name"])
	})

	t.Run("ArrayJSON string", func(t *testing.T) {
		m := &config.Mock{ArrayJSON: `[{"id": 1}, {"id": 2}]`}
		source, err := Compile(m, CompileOptions{})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)
	})

	t.Run("Filter with array source", func(t *testing.T) {
		data := []any{
			map[string]any{"id": 1, "name": "Alice"},
			map[string]any{"id": 2, "name": "Bob"},
			map[string]any{"id": 3, "name": "Charlie"},
		}
		m := &config.Mock{
			Array:  data,
			Filter: `return row.id > 1`,
		}
		source, err := Compile(m, CompileOptions{})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)
	})

	t.Run("Filter with input matching", func(t *testing.T) {
		data := []any{
			map[string]any{"id": 1, "name": "Alice"},
			map[string]any{"id": 2, "name": "Bob"},
		}
		m := &config.Mock{
			Array:  data,
			Filter: `return row.id == input.id`,
		}
		source, err := Compile(m, CompileOptions{})
		require.NoError(t, err)

		input := map[string]any{"id": float64(2)}
		result, _, err := source.Data(nil, "", input, nil)
		require.NoError(t, err)
		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 1)
		assert.Equal(t, "Bob", resultSlice[0].(map[string]any)["name"])
	})
}
