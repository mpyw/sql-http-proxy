package js

import (
	"net/http"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResponse(t *testing.T) {
	r := NewResponse()

	assert.Equal(t, http.StatusOK, r.Status())
	assert.Equal(t, "OK", r.StatusText())
	assert.NotNil(t, r.Headers())
	assert.True(t, r.Ok())
}

func TestResponse_Status(t *testing.T) {
	r := NewResponse()

	t.Run("default status is 200", func(t *testing.T) {
		assert.Equal(t, 200, r.Status())
	})

	t.Run("SetStatus with valid code", func(t *testing.T) {
		r.SetStatus(201)
		assert.Equal(t, 201, r.Status())
		assert.Equal(t, "Created", r.StatusText())
	})

	t.Run("SetStatus with 404", func(t *testing.T) {
		r.SetStatus(404)
		assert.Equal(t, 404, r.Status())
		assert.Equal(t, "Not Found", r.StatusText())
	})

	t.Run("SetStatus with unknown code sets Unknown text", func(t *testing.T) {
		r.SetStatus(299) // Valid but no standard text
		assert.Equal(t, 299, r.Status())
		assert.Equal(t, "Unknown", r.StatusText())
	})

	t.Run("SetStatus ignores code below 100", func(t *testing.T) {
		r.SetStatus(200)
		r.SetStatus(99)
		assert.Equal(t, 200, r.Status()) // Unchanged
	})

	t.Run("SetStatus ignores code above 599", func(t *testing.T) {
		r.SetStatus(200)
		r.SetStatus(600)
		assert.Equal(t, 200, r.Status()) // Unchanged
	})
}

func TestResponse_StatusText(t *testing.T) {
	r := NewResponse()

	t.Run("default statusText is OK", func(t *testing.T) {
		assert.Equal(t, "OK", r.StatusText())
	})

	t.Run("SetStatusText", func(t *testing.T) {
		r.SetStatusText("Custom Text")
		assert.Equal(t, "Custom Text", r.StatusText())
	})
}

func TestResponse_Ok(t *testing.T) {
	r := NewResponse()

	t.Run("Ok returns true for 2xx", func(t *testing.T) {
		r.SetStatus(200)
		assert.True(t, r.Ok())

		r.SetStatus(201)
		assert.True(t, r.Ok())

		r.SetStatus(299)
		assert.True(t, r.Ok())
	})

	t.Run("Ok returns false for non-2xx", func(t *testing.T) {
		r.SetStatus(199)
		assert.False(t, r.Ok())

		r.SetStatus(300)
		assert.False(t, r.Ok())

		r.SetStatus(404)
		assert.False(t, r.Ok())

		r.SetStatus(500)
		assert.False(t, r.Ok())
	})
}

func TestResponse_Headers(t *testing.T) {
	r := NewResponse()

	t.Run("Headers returns writable headers", func(t *testing.T) {
		h := r.Headers()
		require.NotNil(t, h)

		h.Set("X-Custom", "value")
		assert.Equal(t, "value", h.Get("X-Custom"))
	})
}

func TestResponse_ToHTTPHeader(t *testing.T) {
	r := NewResponse()
	r.Headers().Set("X-Test", "test-value")

	httpHeader := r.ToHTTPHeader()
	assert.Equal(t, "test-value", httpHeader.Get("X-Test"))
}

func TestResponse_ToJSObject(t *testing.T) {
	r := NewResponse()
	vm := goja.New()
	obj := r.ToJSObject(vm)

	err := vm.Set("response", obj)
	require.NoError(t, err)

	t.Run("status get", func(t *testing.T) {
		result, err := vm.RunString(`response.status`)
		require.NoError(t, err)
		assert.Equal(t, int64(200), result.ToInteger())
	})

	t.Run("status set with int", func(t *testing.T) {
		_, err := vm.RunString(`response.status = 201`)
		require.NoError(t, err)
		assert.Equal(t, 201, r.Status())
	})

	t.Run("status set with float64", func(t *testing.T) {
		_, err := vm.RunString(`response.status = 202.5`)
		require.NoError(t, err)
		assert.Equal(t, 202, r.Status()) // Converted via ToInteger()
	})

	t.Run("statusText get", func(t *testing.T) {
		r.SetStatus(200)
		result, err := vm.RunString(`response.statusText`)
		require.NoError(t, err)
		assert.Equal(t, "OK", result.String())
	})

	t.Run("statusText set", func(t *testing.T) {
		_, err := vm.RunString(`response.statusText = "Custom Status"`)
		require.NoError(t, err)
		assert.Equal(t, "Custom Status", r.StatusText())
	})

	t.Run("ok property", func(t *testing.T) {
		r.SetStatus(200)
		result, err := vm.RunString(`response.ok`)
		require.NoError(t, err)
		assert.True(t, result.ToBoolean())

		r.SetStatus(404)
		result, err = vm.RunString(`response.ok`)
		require.NoError(t, err)
		assert.False(t, result.ToBoolean())
	})

	t.Run("headers.set()", func(t *testing.T) {
		_, err := vm.RunString(`response.headers.set("X-JS-Header", "js-value")`)
		require.NoError(t, err)
		assert.Equal(t, "js-value", r.Headers().Get("X-JS-Header"))
	})

	t.Run("headers.get()", func(t *testing.T) {
		r.Headers().Set("X-Go-Header", "go-value")
		result, err := vm.RunString(`response.headers.get("X-Go-Header")`)
		require.NoError(t, err)
		assert.Equal(t, "go-value", result.String())
	})

	t.Run("headers.append()", func(t *testing.T) {
		_, err := vm.RunString(`
			response.headers.set("X-Multi", "first");
			response.headers.append("X-Multi", "second");
		`)
		require.NoError(t, err)
		values := r.ToHTTPHeader().Values("X-Multi")
		assert.Contains(t, values, "first")
		assert.Contains(t, values, "second")
	})

	t.Run("headers.delete()", func(t *testing.T) {
		r.Headers().Set("X-Delete", "to-delete")
		_, err := vm.RunString(`response.headers.delete("X-Delete")`)
		require.NoError(t, err)
		assert.Nil(t, r.Headers().Get("X-Delete"))
	})

	t.Run("response is sealed - new property ignored", func(t *testing.T) {
		// Sealed objects silently ignore new property assignments in non-strict mode
		_, err := vm.RunString(`response.newProp = "test"`)
		require.NoError(t, err)

		result, err := vm.RunString(`response.newProp`)
		require.NoError(t, err)
		assert.True(t, goja.IsUndefined(result)) // Property was not added
	})

	t.Run("ok is read-only - assignment ignored", func(t *testing.T) {
		// Read-only properties silently ignore assignments in non-strict mode
		r.SetStatus(200)
		_, err := vm.RunString(`response.ok = false`)
		require.NoError(t, err)

		result, err := vm.RunString(`response.ok`)
		require.NoError(t, err)
		assert.True(t, result.ToBoolean()) // Value unchanged (still true for 200)
	})
}
