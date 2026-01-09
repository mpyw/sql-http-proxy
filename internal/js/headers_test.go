package js

import (
	"net/http"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHeaders(t *testing.T) {
	t.Run("nil header", func(t *testing.T) {
		h := NewHeaders(nil, false)
		assert.NotNil(t, h)
		assert.NotNil(t, h.headers)
	})

	t.Run("with header", func(t *testing.T) {
		header := http.Header{"Content-Type": []string{"application/json"}}
		h := NewHeaders(header, false)
		assert.Equal(t, "application/json", h.Get("Content-Type"))
	})

	t.Run("readonly", func(t *testing.T) {
		h := NewReadonlyHeaders(nil)
		assert.True(t, h.readonly)
	})

	t.Run("writable", func(t *testing.T) {
		h := NewWritableHeaders(nil)
		assert.False(t, h.readonly)
	})
}

func TestHeaders_Get(t *testing.T) {
	h := NewHeaders(http.Header{
		"Content-Type": []string{"application/json"},
	}, false)

	t.Run("existing header", func(t *testing.T) {
		result := h.Get("Content-Type")
		assert.Equal(t, "application/json", result)
	})

	t.Run("case insensitive", func(t *testing.T) {
		result := h.Get("content-type")
		assert.Equal(t, "application/json", result)
	})

	t.Run("non-existing header", func(t *testing.T) {
		result := h.Get("X-Custom")
		assert.Nil(t, result)
	})
}

func TestHeaders_Has(t *testing.T) {
	h := NewHeaders(http.Header{
		"Content-Type": []string{"application/json"},
	}, false)

	t.Run("existing header", func(t *testing.T) {
		assert.True(t, h.Has("Content-Type"))
	})

	t.Run("case insensitive", func(t *testing.T) {
		assert.True(t, h.Has("content-type"))
	})

	t.Run("non-existing header", func(t *testing.T) {
		assert.False(t, h.Has("X-Custom"))
	})
}

func TestHeaders_Set(t *testing.T) {
	t.Run("writable", func(t *testing.T) {
		h := NewWritableHeaders(nil)
		h.Set("Content-Type", "text/plain")
		assert.Equal(t, "text/plain", h.Get("Content-Type"))
	})

	t.Run("replaces existing", func(t *testing.T) {
		h := NewWritableHeaders(http.Header{"Content-Type": []string{"application/json"}})
		h.Set("Content-Type", "text/plain")
		assert.Equal(t, "text/plain", h.Get("Content-Type"))
	})

	t.Run("readonly ignored", func(t *testing.T) {
		h := NewReadonlyHeaders(nil)
		h.Set("Content-Type", "text/plain")
		assert.Nil(t, h.Get("Content-Type"))
	})
}

func TestHeaders_Append(t *testing.T) {
	t.Run("writable", func(t *testing.T) {
		h := NewWritableHeaders(nil)
		h.Append("Accept", "text/plain")
		h.Append("Accept", "application/json")
		values := h.headers["Accept"]
		assert.Len(t, values, 2)
		assert.Contains(t, values, "text/plain")
		assert.Contains(t, values, "application/json")
	})

	t.Run("readonly ignored", func(t *testing.T) {
		h := NewReadonlyHeaders(nil)
		h.Append("Accept", "text/plain")
		assert.Nil(t, h.Get("Accept"))
	})
}

func TestHeaders_Delete(t *testing.T) {
	t.Run("writable", func(t *testing.T) {
		h := NewWritableHeaders(http.Header{"Content-Type": []string{"application/json"}})
		h.Delete("Content-Type")
		assert.Nil(t, h.Get("Content-Type"))
	})

	t.Run("readonly ignored", func(t *testing.T) {
		h := NewReadonlyHeaders(http.Header{"Content-Type": []string{"application/json"}})
		h.Delete("Content-Type")
		assert.Equal(t, "application/json", h.Get("Content-Type"))
	})
}

func TestHeaders_Entries(t *testing.T) {
	h := NewHeaders(http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"text/plain", "text/html"},
	}, false)

	entries := h.Entries()
	assert.Len(t, entries, 3)
	// Entries are sorted by lowercase name
	assert.Equal(t, "accept", entries[0][0])
	assert.Equal(t, "accept", entries[1][0])
	assert.Equal(t, "content-type", entries[2][0])
}

func TestHeaders_Keys(t *testing.T) {
	h := NewHeaders(http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"text/plain"},
	}, false)

	keys := h.Keys()
	assert.Len(t, keys, 2)
	// Keys are sorted lowercase
	assert.Equal(t, "accept", keys[0])
	assert.Equal(t, "content-type", keys[1])
}

func TestHeaders_Values(t *testing.T) {
	h := NewHeaders(http.Header{
		"Accept": []string{"text/plain", "application/json"},
	}, false)

	values := h.Values()
	assert.Len(t, values, 2)
	assert.Contains(t, values, "text/plain")
	assert.Contains(t, values, "application/json")
}

func TestHeaders_ForEach(t *testing.T) {
	h := NewHeaders(http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"text/plain"},
	}, false)

	vm := goja.New()
	var calls []struct {
		key   string
		value string
	}

	callback := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		value := call.Argument(0).String()
		key := call.Argument(1).String()
		calls = append(calls, struct{ key, value string }{key, value})
		return goja.Undefined()
	})
	callable, _ := goja.AssertFunction(callback)

	err := h.ForEach(vm, callable)
	require.NoError(t, err)
	assert.Len(t, calls, 2)
}

func TestHeaders_ToHTTPHeader(t *testing.T) {
	original := http.Header{"Content-Type": []string{"application/json"}}
	h := NewHeaders(original, false)
	result := h.ToHTTPHeader()
	assert.Equal(t, original, result)
}

func TestHeaders_ToJSObject(t *testing.T) {
	vm := goja.New()
	h := NewWritableHeaders(http.Header{
		"Content-Type": []string{"application/json"},
	})

	obj := h.ToJSObject(vm)
	require.NotNil(t, obj)

	t.Run("get method", func(t *testing.T) {
		result, err := vm.RunString(`
			var headers = this.headers;
			headers.get("Content-Type");
		`)
		_ = vm.Set("headers", obj)
		result, err = vm.RunString(`headers.get("Content-Type")`)
		require.NoError(t, err)
		assert.Equal(t, "application/json", result.String())
	})

	t.Run("has method", func(t *testing.T) {
		result, err := vm.RunString(`headers.has("Content-Type")`)
		require.NoError(t, err)
		assert.True(t, result.ToBoolean())

		result, err = vm.RunString(`headers.has("X-Missing")`)
		require.NoError(t, err)
		assert.False(t, result.ToBoolean())
	})

	t.Run("set method", func(t *testing.T) {
		_, err := vm.RunString(`headers.set("X-Custom", "value")`)
		require.NoError(t, err)
		assert.Equal(t, "value", h.Get("X-Custom"))
	})

	t.Run("append method", func(t *testing.T) {
		_, err := vm.RunString(`headers.append("X-Multi", "val1")`)
		require.NoError(t, err)
		_, err = vm.RunString(`headers.append("X-Multi", "val2")`)
		require.NoError(t, err)
		values := h.headers["X-Multi"]
		assert.Len(t, values, 2)
	})

	t.Run("delete method", func(t *testing.T) {
		h.Set("X-Delete", "test")
		_, err := vm.RunString(`headers.delete("X-Delete")`)
		require.NoError(t, err)
		assert.Nil(t, h.Get("X-Delete"))
	})

	t.Run("keys method", func(t *testing.T) {
		result, err := vm.RunString(`headers.keys()`)
		require.NoError(t, err)
		assert.NotNil(t, result.Export())
	})

	t.Run("values method", func(t *testing.T) {
		result, err := vm.RunString(`headers.values()`)
		require.NoError(t, err)
		assert.NotNil(t, result.Export())
	})

	t.Run("entries method", func(t *testing.T) {
		result, err := vm.RunString(`headers.entries()`)
		require.NoError(t, err)
		assert.NotNil(t, result.Export())
	})

	t.Run("forEach method", func(t *testing.T) {
		vm.Set("callCount", 0)
		_, err := vm.RunString(`
			headers.forEach(function(value, key) {
				callCount++;
			})
		`)
		require.NoError(t, err)
		callCount := vm.Get("callCount").ToInteger()
		assert.True(t, callCount > 0)
	})
}
