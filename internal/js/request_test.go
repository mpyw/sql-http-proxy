package js

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test/path?foo=bar", nil)
	req.Header.Set("X-Custom", "value")
	req.Header.Set("Content-Type", "application/json")

	r := NewRequest(req)

	t.Run("Method", func(t *testing.T) {
		assert.Equal(t, "POST", r.Method())
	})

	t.Run("URL", func(t *testing.T) {
		assert.Equal(t, "/test/path?foo=bar", r.URL())
	})

	t.Run("Headers", func(t *testing.T) {
		h := r.Headers()
		require.NotNil(t, h)
		assert.Equal(t, "value", h.Get("X-Custom"))
		assert.Equal(t, "application/json", h.Get("Content-Type"))
	})
}

func TestRequest_ToJSObject(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users?id=123", nil)
	req.Header.Set("Authorization", "Bearer token")

	r := NewRequest(req)
	vm := goja.New()
	obj := r.ToJSObject(vm)

	err := vm.Set("request", obj)
	require.NoError(t, err)

	t.Run("method property", func(t *testing.T) {
		result, err := vm.RunString(`request.method`)
		require.NoError(t, err)
		assert.Equal(t, "GET", result.String())
	})

	t.Run("url property", func(t *testing.T) {
		result, err := vm.RunString(`request.url`)
		require.NoError(t, err)
		assert.Equal(t, "/api/users?id=123", result.String())
	})

	t.Run("headers.get()", func(t *testing.T) {
		result, err := vm.RunString(`request.headers.get("Authorization")`)
		require.NoError(t, err)
		assert.Equal(t, "Bearer token", result.String())
	})

	t.Run("headers.has()", func(t *testing.T) {
		result, err := vm.RunString(`request.headers.has("Authorization")`)
		require.NoError(t, err)
		assert.True(t, result.ToBoolean())

		result, err = vm.RunString(`request.headers.has("X-Missing")`)
		require.NoError(t, err)
		assert.False(t, result.ToBoolean())
	})

	t.Run("request is sealed - new property ignored", func(t *testing.T) {
		// Sealed objects silently ignore new property assignments in non-strict mode
		_, err := vm.RunString(`request.newProp = "test"`)
		require.NoError(t, err)

		result, err := vm.RunString(`request.newProp`)
		require.NoError(t, err)
		assert.True(t, goja.IsUndefined(result)) // Property was not added
	})

	t.Run("method is read-only - assignment ignored", func(t *testing.T) {
		// Read-only properties silently ignore assignments in non-strict mode
		_, err := vm.RunString(`request.method = "PUT"`)
		require.NoError(t, err)

		result, err := vm.RunString(`request.method`)
		require.NoError(t, err)
		assert.Equal(t, "GET", result.String()) // Value unchanged
	})
}
