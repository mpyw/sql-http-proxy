package js

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformError(t *testing.T) {
	t.Run("string body returns message", func(t *testing.T) {
		err := &TransformError{Status: 400, Body: "bad request"}
		assert.Equal(t, "bad request", err.Error())
	})

	t.Run("object body with message", func(t *testing.T) {
		err := &TransformError{
			Status: 400,
			Body:   map[string]any{"message": "validation error", "code": "INVALID"},
		}
		assert.Equal(t, "validation error", err.Error())
	})

	t.Run("object body without message", func(t *testing.T) {
		err := &TransformError{
			Status: 400,
			Body:   map[string]any{"code": "INVALID"},
		}
		assert.Equal(t, "transform error", err.Error())
	})

	t.Run("nil body", func(t *testing.T) {
		err := &TransformError{Status: 204, Body: nil}
		assert.Equal(t, "transform error", err.Error())
	})

	t.Run("numeric body", func(t *testing.T) {
		err := &TransformError{Status: 400, Body: 123}
		assert.Equal(t, "transform error", err.Error())
	})
}

func TestParseJSError_LambdaStyle(t *testing.T) {
	t.Run("status and body object", func(t *testing.T) {
		transformer, err := Compile(`throw { status: 400, body: { message: "bad request" } }`, nil)
		require.NoError(t, err)

		_, err = transformer.Apply()
		require.Error(t, err)

		var transformErr *TransformError
		require.ErrorAs(t, err, &transformErr)
		assert.Equal(t, 400, transformErr.Status)

		body := transformErr.Body.(map[string]any)
		assert.Equal(t, "bad request", body["message"])
	})

	t.Run("status and body string", func(t *testing.T) {
		transformer, err := Compile(`throw { status: 422, body: "validation error" }`, nil)
		require.NoError(t, err)

		_, err = transformer.Apply()
		require.Error(t, err)

		var transformErr *TransformError
		require.ErrorAs(t, err, &transformErr)
		assert.Equal(t, 422, transformErr.Status)
		assert.Equal(t, "validation error", transformErr.Body)
	})

	t.Run("status only (no body)", func(t *testing.T) {
		transformer, err := Compile(`throw { status: 500 }`, nil)
		require.NoError(t, err)

		_, err = transformer.Apply()
		require.Error(t, err)

		var transformErr *TransformError
		require.ErrorAs(t, err, &transformErr)
		assert.Equal(t, 500, transformErr.Status)
	})

	t.Run("body only (no status)", func(t *testing.T) {
		transformer, err := Compile(`throw { body: "error message" }`, nil)
		require.NoError(t, err)

		_, err = transformer.Apply()
		require.Error(t, err)

		var transformErr *TransformError
		require.ErrorAs(t, err, &transformErr)
		assert.Equal(t, 0, transformErr.Status) // Default status
		assert.Equal(t, "error message", transformErr.Body)
	})

	t.Run("status as float64", func(t *testing.T) {
		transformer, err := Compile(`throw { status: 400.0, body: "error" }`, nil)
		require.NoError(t, err)

		_, err = transformer.Apply()
		require.Error(t, err)

		var transformErr *TransformError
		require.ErrorAs(t, err, &transformErr)
		assert.Equal(t, 400, transformErr.Status)
	})
}

func TestParseJSError_NativeError(t *testing.T) {
	t.Run("new Error returns 500", func(t *testing.T) {
		transformer, err := Compile(`throw new Error("something went wrong")`, nil)
		require.NoError(t, err)

		_, err = transformer.Apply()
		require.Error(t, err)

		var transformErr *TransformError
		require.ErrorAs(t, err, &transformErr)
		assert.Equal(t, 500, transformErr.Status)
	})

	t.Run("object without status or body returns 500", func(t *testing.T) {
		transformer, err := Compile(`throw { message: "error", code: "ERR_001" }`, nil)
		require.NoError(t, err)

		_, err = transformer.Apply()
		require.Error(t, err)

		var transformErr *TransformError
		require.ErrorAs(t, err, &transformErr)
		assert.Equal(t, 500, transformErr.Status)
	})
}

func TestParseJSError_ThrownString(t *testing.T) {
	t.Run("thrown string", func(t *testing.T) {
		transformer, err := Compile(`throw "simple error message"`, nil)
		require.NoError(t, err)

		_, err = transformer.Apply()
		require.Error(t, err)

		var transformErr *TransformError
		require.ErrorAs(t, err, &transformErr)
		assert.Equal(t, 0, transformErr.Status)
		assert.Equal(t, "simple error message", transformErr.Body)
	})
}
