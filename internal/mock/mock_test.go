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
	t.Run("nil MockTransform", func(t *testing.T) {
		source, err := Compile(nil, CompileOptions{})
		require.NoError(t, err)
		assert.Nil(t, source)
	})

	t.Run("JS inline", func(t *testing.T) {
		mt := &config.MockTransform{JS: `return { id: 1 }`}
		source, err := Compile(mt, CompileOptions{})
		require.NoError(t, err)
		assert.NotNil(t, source)

		result, _, err := source.Data(nil, "", nil)
		require.NoError(t, err)
		resultMap := result.(map[string]any)
		assert.Equal(t, int64(1), resultMap["id"])
	})

	t.Run("JS file with relative path", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "mock.js")
		err := os.WriteFile(path, []byte(`return { id: 2 }`), 0644)
		require.NoError(t, err)

		mt := &config.MockTransform{JSFile: "mock.js"}
		source, err := Compile(mt, CompileOptions{ConfigDir: dir})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil)
		require.NoError(t, err)
		resultMap := result.(map[string]any)
		assert.Equal(t, int64(2), resultMap["id"])
	})

	t.Run("JS file with absolute path", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "mock.js")
		err := os.WriteFile(path, []byte(`return { id: 3 }`), 0644)
		require.NoError(t, err)

		mt := &config.MockTransform{JSFile: path}
		source, err := Compile(mt, CompileOptions{ConfigDir: "/different/config/dir"})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil)
		require.NoError(t, err)
		resultMap := result.(map[string]any)
		assert.Equal(t, int64(3), resultMap["id"])
	})

	t.Run("CSV inline", func(t *testing.T) {
		mt := &config.MockTransform{CSV: "id,name\n42,Alice"}
		source, err := Compile(mt, CompileOptions{})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil)
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

		mt := &config.MockTransform{CSVFile: "data.csv"}
		source, err := Compile(mt, CompileOptions{ConfigDir: dir})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil)
		require.NoError(t, err)
		resultSlice := result.([]map[string]any)
		assert.Len(t, resultSlice, 1)
	})

	t.Run("JSON inline", func(t *testing.T) {
		data := []any{map[string]any{"id": 1}}
		mt := &config.MockTransform{JSON: data}
		source, err := Compile(mt, CompileOptions{})
		require.NoError(t, err)
		assert.NotNil(t, source)
	})

	t.Run("JSON file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "data.json")
		err := os.WriteFile(path, []byte(`[{"id": 1}]`), 0644)
		require.NoError(t, err)

		mt := &config.MockTransform{JSONFile: "data.json"}
		source, err := Compile(mt, CompileOptions{ConfigDir: dir})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil)
		require.NoError(t, err)
		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 1)
	})

	t.Run("JSONL inline", func(t *testing.T) {
		mt := &config.MockTransform{JSONL: `{"id": 1}
{"id": 2}`}
		source, err := Compile(mt, CompileOptions{})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil)
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

		mt := &config.MockTransform{JSONLFile: "data.jsonl"}
		source, err := Compile(mt, CompileOptions{ConfigDir: dir})
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil)
		require.NoError(t, err)
		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)
	})

	t.Run("empty MockTransform", func(t *testing.T) {
		mt := &config.MockTransform{}
		_, err := Compile(mt, CompileOptions{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no source specified")
	})

	t.Run("validation error - multiple sources", func(t *testing.T) {
		mt := &config.MockTransform{
			JS:  `return {}`,
			CSV: "id\n42",
		}
		_, err := Compile(mt, CompileOptions{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only one source type can be specified")
	})

	t.Run("JS compile error", func(t *testing.T) {
		mt := &config.MockTransform{JS: `return { invalid`}
		_, err := Compile(mt, CompileOptions{})
		require.Error(t, err)
	})

	t.Run("JS file not found", func(t *testing.T) {
		mt := &config.MockTransform{JSFile: "nonexistent.js"}
		_, err := Compile(mt, CompileOptions{ConfigDir: "/tmp"})
		require.Error(t, err)
	})

	t.Run("CSV parse error", func(t *testing.T) {
		mt := &config.MockTransform{CSV: ""} // Empty CSV
		_, err := Compile(mt, CompileOptions{})
		require.Error(t, err)
	})

	t.Run("CSV file not found", func(t *testing.T) {
		mt := &config.MockTransform{CSVFile: "nonexistent.csv"}
		_, err := Compile(mt, CompileOptions{ConfigDir: "/tmp"})
		require.Error(t, err)
	})

	t.Run("JSON file not found", func(t *testing.T) {
		mt := &config.MockTransform{JSONFile: "nonexistent.json"}
		_, err := Compile(mt, CompileOptions{ConfigDir: "/tmp"})
		require.Error(t, err)
	})

	t.Run("JSONL parse error", func(t *testing.T) {
		mt := &config.MockTransform{JSONL: `{invalid json}`}
		_, err := Compile(mt, CompileOptions{})
		require.Error(t, err)
	})

	t.Run("JSONL file not found", func(t *testing.T) {
		mt := &config.MockTransform{JSONLFile: "nonexistent.jsonl"}
		_, err := Compile(mt, CompileOptions{ConfigDir: "/tmp"})
		require.Error(t, err)
	})
}
