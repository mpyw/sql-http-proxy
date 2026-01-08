package mock

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCSV(t *testing.T) {
	t.Run("basic CSV", func(t *testing.T) {
		// Note: "1" and "0" are parsed as booleans (PHP filter_var style)
		// Use numbers like 42, 100 to get float64
		data := `id,name,active
42,Alice,true
100,Bob,false`
		source, err := ParseCSV(data)
		require.NoError(t, err)

		rows, ctx, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)
		assert.Nil(t, ctx)

		rowSlice := rows.([]map[string]any)
		assert.Len(t, rowSlice, 2)

		assert.Equal(t, float64(42), rowSlice[0]["id"])
		assert.Equal(t, "Alice", rowSlice[0]["name"])
		assert.Equal(t, true, rowSlice[0]["active"])

		assert.Equal(t, float64(100), rowSlice[1]["id"])
		assert.Equal(t, "Bob", rowSlice[1]["name"])
		assert.Equal(t, false, rowSlice[1]["active"])
	})

	t.Run("type coercion", func(t *testing.T) {
		data := `string,number,bool,null
hello,42,yes,null
world,-3.14,off,NULL`
		source, err := ParseCSV(data)
		require.NoError(t, err)

		rows, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		rowSlice := rows.([]map[string]any)
		assert.Len(t, rowSlice, 2)

		// First row
		assert.Equal(t, "hello", rowSlice[0]["string"])
		assert.Equal(t, float64(42), rowSlice[0]["number"])
		assert.Equal(t, true, rowSlice[0]["bool"])
		assert.Nil(t, rowSlice[0]["null"])

		// Second row
		assert.Equal(t, "world", rowSlice[1]["string"])
		assert.Equal(t, float64(-3.14), rowSlice[1]["number"])
		assert.Equal(t, false, rowSlice[1]["bool"])
		assert.Nil(t, rowSlice[1]["null"])
	})

	t.Run("empty data", func(t *testing.T) {
		_, err := ParseCSV("")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "empty data")
	})

	t.Run("header only", func(t *testing.T) {
		data := "id,name,active"
		source, err := ParseCSV(data)
		require.NoError(t, err)

		rows, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		rowSlice := rows.([]map[string]any)
		assert.Len(t, rowSlice, 0)
	})

	t.Run("fewer columns than header", func(t *testing.T) {
		data := `id,name,active
42,Alice`
		source, err := ParseCSV(data)
		require.NoError(t, err)

		rows, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		rowSlice := rows.([]map[string]any)
		assert.Len(t, rowSlice, 1)

		assert.Equal(t, float64(42), rowSlice[0]["id"])
		assert.Equal(t, "Alice", rowSlice[0]["name"])
		assert.Equal(t, "", rowSlice[0]["active"]) // Missing column = empty string
	})

	t.Run("more columns than header", func(t *testing.T) {
		data := `id,name
42,Alice,extra,ignored`
		source, err := ParseCSV(data)
		require.NoError(t, err)

		rows, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		rowSlice := rows.([]map[string]any)
		assert.Len(t, rowSlice, 1)

		assert.Equal(t, float64(42), rowSlice[0]["id"])
		assert.Equal(t, "Alice", rowSlice[0]["name"])
		// Extra columns are ignored
		_, exists := rowSlice[0]["extra"]
		assert.False(t, exists)
	})

	t.Run("data returns copy", func(t *testing.T) {
		data := `id,name
42,Alice`
		source, err := ParseCSV(data)
		require.NoError(t, err)

		rows1, _ := lo.Must2(source.Data(nil, "", nil, nil))
		rows2, _ := lo.Must2(source.Data(nil, "", nil, nil))

		// Modify first result
		rowSlice1 := rows1.([]map[string]any)
		rowSlice1[0]["name"] = "Modified"

		// Second result should be unchanged
		rowSlice2 := rows2.([]map[string]any)
		assert.Equal(t, "Alice", rowSlice2[0]["name"])
	})
}

func TestParseCSVFile(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		// Create temp file
		dir := t.TempDir()
		path := filepath.Join(dir, "test.csv")
		content := `id,name
1,Alice
2,Bob`
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)

		source, err := ParseCSVFile(path)
		require.NoError(t, err)

		rows, _, err := source.Data(nil, "", nil, nil)
		require.NoError(t, err)

		rowSlice := rows.([]map[string]any)
		assert.Len(t, rowSlice, 2)
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := ParseCSVFile("/nonexistent/file.csv")
		require.Error(t, err)
	})
}
