package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestMock_IsEmpty(t *testing.T) {
	t.Run("empty mock", func(t *testing.T) {
		m := &Mock{}
		assert.True(t, m.IsEmpty())
	})

	t.Run("with object", func(t *testing.T) {
		m := &Mock{Object: map[string]any{"id": 1}}
		assert.False(t, m.IsEmpty())
	})

	t.Run("with array", func(t *testing.T) {
		m := &Mock{Array: []any{map[string]any{"id": 1}}}
		assert.False(t, m.IsEmpty())
	})

	t.Run("with CSV", func(t *testing.T) {
		m := &Mock{CSV: "id,name\n1,Alice"}
		assert.False(t, m.IsEmpty())
	})
}

func TestMock_IsArraySource(t *testing.T) {
	t.Run("empty mock", func(t *testing.T) {
		m := &Mock{}
		assert.False(t, m.IsArraySource())
	})

	t.Run("with object", func(t *testing.T) {
		m := &Mock{Object: map[string]any{"id": 1}}
		assert.False(t, m.IsArraySource())
	})

	t.Run("with array", func(t *testing.T) {
		m := &Mock{Array: []any{map[string]any{"id": 1}}}
		assert.True(t, m.IsArraySource())
	})

	t.Run("with array_json", func(t *testing.T) {
		m := &Mock{ArrayJSON: `[{"id": 1}]`}
		assert.True(t, m.IsArraySource())
	})

	t.Run("with csv", func(t *testing.T) {
		m := &Mock{CSV: "id\n1"}
		assert.True(t, m.IsArraySource())
	})

	t.Run("with csv_file", func(t *testing.T) {
		m := &Mock{CSVFile: "test.csv"}
		assert.True(t, m.IsArraySource())
	})

	t.Run("with jsonl", func(t *testing.T) {
		m := &Mock{JSONL: `{"id": 1}`}
		assert.True(t, m.IsArraySource())
	})

	t.Run("with jsonl_file", func(t *testing.T) {
		m := &Mock{JSONLFile: "test.jsonl"}
		assert.True(t, m.IsArraySource())
	})

	t.Run("with array_js", func(t *testing.T) {
		m := &Mock{ArrayJS: "return []"}
		assert.True(t, m.IsArraySource())
	})

	t.Run("with array_json_file", func(t *testing.T) {
		m := &Mock{ArrayJSONFile: "test.json"}
		assert.True(t, m.IsArraySource())
	})
}

func TestMock_IsObjectSource(t *testing.T) {
	t.Run("empty mock", func(t *testing.T) {
		m := &Mock{}
		assert.False(t, m.IsObjectSource())
	})

	t.Run("with array", func(t *testing.T) {
		m := &Mock{Array: []any{map[string]any{"id": 1}}}
		assert.False(t, m.IsObjectSource())
	})

	t.Run("with object", func(t *testing.T) {
		m := &Mock{Object: map[string]any{"id": 1}}
		assert.True(t, m.IsObjectSource())
	})

	t.Run("with object_json", func(t *testing.T) {
		m := &Mock{ObjectJSON: `{"id": 1}`}
		assert.True(t, m.IsObjectSource())
	})

	t.Run("with object_json_file", func(t *testing.T) {
		m := &Mock{ObjectJSONFile: "test.json"}
		assert.True(t, m.IsObjectSource())
	})

	t.Run("with object_js", func(t *testing.T) {
		m := &Mock{ObjectJS: "return {}"}
		assert.True(t, m.IsObjectSource())
	})
}

func TestMock_GetJS(t *testing.T) {
	t.Run("empty mock", func(t *testing.T) {
		m := &Mock{}
		assert.Equal(t, "", m.GetJS())
	})

	t.Run("with object_js", func(t *testing.T) {
		m := &Mock{ObjectJS: "return {id: 1}"}
		assert.Equal(t, "return {id: 1}", m.GetJS())
	})

	t.Run("with array_js", func(t *testing.T) {
		m := &Mock{ArrayJS: "return [{id: 1}]"}
		assert.Equal(t, "return [{id: 1}]", m.GetJS())
	})

	t.Run("object_js takes priority", func(t *testing.T) {
		m := &Mock{ObjectJS: "return {}", ArrayJS: "return []"}
		assert.Equal(t, "return {}", m.GetJS())
	})
}

func TestMock_Validate(t *testing.T) {
	t.Run("empty mock", func(t *testing.T) {
		m := &Mock{}
		assert.NoError(t, m.Validate())
	})

	t.Run("single source", func(t *testing.T) {
		m := &Mock{Object: map[string]any{"id": 1}}
		assert.NoError(t, m.Validate())
	})

	t.Run("multiple sources", func(t *testing.T) {
		m := &Mock{
			Object: map[string]any{"id": 1},
			Array:  []any{map[string]any{"id": 1}},
		}
		err := m.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only one source type can be specified")
	})

	t.Run("multiple array sources", func(t *testing.T) {
		m := &Mock{
			CSV:   "id\n1",
			JSONL: `{"id": 1}`,
		}
		err := m.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only one source type can be specified")
	})
}

func TestMock_ValidateForTypeOne(t *testing.T) {
	t.Run("object source", func(t *testing.T) {
		m := &Mock{Object: map[string]any{"id": 1}}
		assert.NoError(t, m.ValidateForTypeOne())
	})

	t.Run("array source without filter", func(t *testing.T) {
		m := &Mock{Array: []any{map[string]any{"id": 1}}}
		err := m.ValidateForTypeOne()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "filter is required")
	})

	t.Run("array source with filter", func(t *testing.T) {
		m := &Mock{
			Array:  []any{map[string]any{"id": 1}},
			Filter: "return row.id === parseInt(input.id)",
		}
		assert.NoError(t, m.ValidateForTypeOne())
	})

	t.Run("csv without filter", func(t *testing.T) {
		m := &Mock{CSV: "id\n1"}
		err := m.ValidateForTypeOne()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "filter is required")
	})

	t.Run("csv with filter", func(t *testing.T) {
		m := &Mock{
			CSV:    "id\n1",
			Filter: "return true",
		}
		assert.NoError(t, m.ValidateForTypeOne())
	})

	t.Run("multiple sources error", func(t *testing.T) {
		m := &Mock{
			Object: map[string]any{"id": 1},
			Array:  []any{map[string]any{"id": 1}},
		}
		err := m.ValidateForTypeOne()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only one source type can be specified")
	})
}

func TestMock_ValidateForTypeMany(t *testing.T) {
	t.Run("array source", func(t *testing.T) {
		m := &Mock{Array: []any{map[string]any{"id": 1}}}
		assert.NoError(t, m.ValidateForTypeMany())
	})

	t.Run("object source", func(t *testing.T) {
		m := &Mock{Object: map[string]any{"id": 1}}
		err := m.ValidateForTypeMany()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "object sources are not allowed")
	})

	t.Run("object_json source", func(t *testing.T) {
		m := &Mock{ObjectJSON: `{"id": 1}`}
		err := m.ValidateForTypeMany()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "object sources are not allowed")
	})

	t.Run("multiple sources error", func(t *testing.T) {
		m := &Mock{
			Array: []any{map[string]any{"id": 1}},
			CSV:   "id\n1",
		}
		err := m.ValidateForTypeMany()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only one source type can be specified")
	})
}

func TestQuery_GetMethod(t *testing.T) {
	t.Run("default method", func(t *testing.T) {
		q := Query{}
		assert.Equal(t, "GET", q.GetMethod())
	})

	t.Run("custom method", func(t *testing.T) {
		q := Query{Method: "POST"}
		assert.Equal(t, "POST", q.GetMethod())
	})
}

func TestQuery_GetAccepts(t *testing.T) {
	t.Run("default accepts when nil", func(t *testing.T) {
		q := Query{}
		assert.Equal(t, DefaultAcceptTypes, q.GetAccepts())
	})

	t.Run("custom accepts", func(t *testing.T) {
		accepts := AcceptTypes{AcceptJSON}
		q := Query{Accepts: &accepts}
		assert.Equal(t, AcceptTypes{AcceptJSON}, q.GetAccepts())
	})

	t.Run("empty accepts explicitly set", func(t *testing.T) {
		accepts := AcceptTypes{}
		q := Query{Accepts: &accepts}
		assert.Equal(t, AcceptTypes{}, q.GetAccepts())
		assert.Empty(t, q.GetAccepts())
	})
}

func TestQuery_HasMock(t *testing.T) {
	t.Run("no mock", func(t *testing.T) {
		q := Query{}
		assert.False(t, q.HasMock())
	})

	t.Run("nil mock", func(t *testing.T) {
		q := Query{Mock: nil}
		assert.False(t, q.HasMock())
	})

	t.Run("empty mock", func(t *testing.T) {
		q := Query{Mock: &Mock{}}
		assert.False(t, q.HasMock())
	})

	t.Run("with mock", func(t *testing.T) {
		q := Query{Mock: &Mock{Object: map[string]any{"id": 1}}}
		assert.True(t, q.HasMock())
	})
}

func TestMutation_GetMethod(t *testing.T) {
	t.Run("default method", func(t *testing.T) {
		m := Mutation{}
		assert.Equal(t, "POST", m.GetMethod())
	})

	t.Run("custom method", func(t *testing.T) {
		m := Mutation{Method: "PUT"}
		assert.Equal(t, "PUT", m.GetMethod())
	})
}

func TestMutation_GetAccepts(t *testing.T) {
	t.Run("default accepts when nil", func(t *testing.T) {
		m := Mutation{}
		assert.Equal(t, DefaultAcceptTypes, m.GetAccepts())
	})

	t.Run("custom accepts", func(t *testing.T) {
		accepts := AcceptTypes{AcceptForm}
		m := Mutation{Accepts: &accepts}
		assert.Equal(t, AcceptTypes{AcceptForm}, m.GetAccepts())
	})

	t.Run("empty accepts explicitly set", func(t *testing.T) {
		accepts := AcceptTypes{}
		m := Mutation{Accepts: &accepts}
		assert.Equal(t, AcceptTypes{}, m.GetAccepts())
		assert.Empty(t, m.GetAccepts())
	})
}

func TestMutation_HasMock(t *testing.T) {
	t.Run("no mock", func(t *testing.T) {
		m := Mutation{}
		assert.False(t, m.HasMock())
	})

	t.Run("with mock", func(t *testing.T) {
		m := Mutation{Mock: &Mock{Object: map[string]any{"id": 1}}}
		assert.True(t, m.HasMock())
	})
}

func TestAcceptTypes_UnmarshalYAML(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		var a AcceptTypes
		err := yaml.Unmarshal([]byte(`json`), &a)
		require.NoError(t, err)
		assert.Equal(t, AcceptTypes{AcceptJSON}, a)
	})

	t.Run("array value", func(t *testing.T) {
		var a AcceptTypes
		err := yaml.Unmarshal([]byte(`[json, form]`), &a)
		require.NoError(t, err)
		assert.Equal(t, AcceptTypes{AcceptJSON, AcceptForm}, a)
	})
}

func TestPostTransform_UnmarshalYAML(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		var p PostTransform
		err := yaml.Unmarshal([]byte(`return output`), &p)
		require.NoError(t, err)
		assert.Equal(t, "return output", p.All)
		assert.Empty(t, p.Each)
	})

	t.Run("object with all", func(t *testing.T) {
		var p PostTransform
		err := yaml.Unmarshal([]byte("all: return output"), &p)
		require.NoError(t, err)
		assert.Equal(t, "return output", p.All)
		assert.Empty(t, p.Each)
	})

	t.Run("object with each", func(t *testing.T) {
		var p PostTransform
		err := yaml.Unmarshal([]byte("each: return output"), &p)
		require.NoError(t, err)
		assert.Equal(t, "return output", p.Each)
		assert.Empty(t, p.All)
	})

	t.Run("object with both", func(t *testing.T) {
		var p PostTransform
		yamlStr := `
each: "return output"
all: "return {data: output}"
`
		err := yaml.Unmarshal([]byte(yamlStr), &p)
		require.NoError(t, err)
		assert.Equal(t, "return output", p.Each)
		assert.Equal(t, "return {data: output}", p.All)
	})
}

func TestGlobalHelpers_UnmarshalYAML(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		var g GlobalHelpers
		err := yaml.Unmarshal([]byte(`function helper() {}`), &g)
		require.NoError(t, err)
		assert.Equal(t, "function helper() {}", g.JS)
		assert.Empty(t, g.JSFiles)
	})

	t.Run("object with js", func(t *testing.T) {
		var g GlobalHelpers
		err := yaml.Unmarshal([]byte("js: function helper() {}"), &g)
		require.NoError(t, err)
		assert.Equal(t, "function helper() {}", g.JS)
	})

	t.Run("object with js_files", func(t *testing.T) {
		var g GlobalHelpers
		err := yaml.Unmarshal([]byte("js_files:\n  - helper.js"), &g)
		require.NoError(t, err)
		assert.Equal(t, []string{"helper.js"}, g.JSFiles)
	})
}
