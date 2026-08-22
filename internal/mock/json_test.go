package mock

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		source, err := newJSON(data)
		require.NoError(t, err)

		result, ctx, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		assert.Nil(t, ctx)

		// Result should be a deep copy
		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)

		row1 := resultSlice[0].(map[string]any)
		assert.Equal(t, 1, row1["id"]) // Preserves original Go type
		assert.Equal(t, "Alice", row1["name"])
	})

	t.Run("object data", func(t *testing.T) {
		data := map[string]any{"id": 1, "name": "Alice"}
		source, err := newJSON(data)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultMap := result.(map[string]any)
		assert.Equal(t, 1, resultMap["id"]) // Preserves original Go type
		assert.Equal(t, "Alice", resultMap["name"])
	})

	t.Run("nil data", func(t *testing.T) {
		source, err := newJSON(nil)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("string data (JSON array)", func(t *testing.T) {
		data := `[{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}]`
		source, err := newJSON(data)
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
		source, err := newJSON(data)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultMap := result.(map[string]any)
		assert.Equal(t, float64(1), resultMap["id"])
		assert.Equal(t, "Alice", resultMap["name"])
	})

	t.Run("invalid JSON string", func(t *testing.T) {
		data := `{invalid json}`
		_, err := newJSON(data)
		require.Error(t, err)
	})

	t.Run("data returns deep copy", func(t *testing.T) {
		data := []any{map[string]any{"name": "Alice"}}
		source, err := newJSON(data)
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

		source, err := parseJSONFile(path)
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

		source, err := parseJSONFile(path)
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

		_, err = parseJSONFile(path)
		require.Error(t, err)
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := parseJSONFile("/nonexistent/file.json")
		require.Error(t, err)
	})
}

func TestParseJSONL(t *testing.T) {
	t.Run("valid JSONL", func(t *testing.T) {
		data := `{"id": 1, "name": "Alice"}
{"id": 2, "name": "Bob"}`
		source, err := parseJSONL(data)
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
		source, err := parseJSONL(data)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2) // Empty lines skipped
	})

	t.Run("empty JSONL", func(t *testing.T) {
		source, err := parseJSONL("")
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		// Empty JSONL returns nil (nil slice preserved by deepCopy)
		assert.Nil(t, result)
	})

	t.Run("invalid JSONL line", func(t *testing.T) {
		data := `{"id": 1}
{invalid json}
{"id": 2}`
		_, err := parseJSONL(data)
		require.Error(t, err)
	})

	t.Run("JSONL with whitespace lines", func(t *testing.T) {
		data := `{"id": 1}

{"id": 2}`
		source, err := parseJSONL(data)
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

		source, err := parseJSONLFile(path)
		require.NoError(t, err)

		result, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		resultSlice := result.([]any)
		assert.Len(t, resultSlice, 2)
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := parseJSONLFile("/nonexistent/file.jsonl")
		require.Error(t, err)
	})
}

func TestParseJSONL_LargeRow(t *testing.T) {
	// The previous bufio.Scanner-based reader capped a line at 64KiB and failed
	// with "token too long". A jsontext.Decoder has no such limit.
	name := strings.Repeat("x", 128*1024)
	source, err := parseJSONL(fmt.Sprintf(`{"name":%q}`+"\n", name))
	require.NoError(t, err)

	rows := source.data.([]any)
	require.Len(t, rows, 1)
	assert.Equal(t, name, rows[0].(map[string]any)["name"])
}

func TestParseJSONL_RejectsIllFormed(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"duplicate object member", `{"id":1,"id":2}`},
		{"invalid UTF-8", "{\"name\":\"\xff\xfe\"}"},
		{"truncated value", `{"id":1`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseJSONL(tt.data)
			require.Error(t, err)
		})
	}
}

func TestParseJSONString_RejectsDuplicateMember(t *testing.T) {
	_, err := parseJSONString(`{"id":1,"id":2}`)
	require.Error(t, err)
}
