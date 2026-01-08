package js

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransforms_IsMock(t *testing.T) {
	t.Run("returns true when Mock is set", func(t *testing.T) {
		transformer, err := CompilePre(`return input`)
		require.NoError(t, err)

		transforms := &Transforms{Mock: transformer}
		assert.True(t, transforms.IsMock())
	})

	t.Run("returns false when Mock is nil", func(t *testing.T) {
		transforms := &Transforms{}
		assert.False(t, transforms.IsMock())
	})
}

func TestCompileTransforms(t *testing.T) {
	t.Run("compiles all transforms", func(t *testing.T) {
		transforms, err := CompileTransforms(
			`return { ...input, preProcessed: true }`,
			`return { id: 1, mocked: true }`,
			`return { ...output, eachProcessed: true }`,
			`return { total: output.length, items: output }`,
		)
		require.NoError(t, err)
		assert.NotNil(t, transforms.Pre)
		assert.NotNil(t, transforms.Mock)
		assert.NotNil(t, transforms.PostEach)
		assert.NotNil(t, transforms.PostAll)
	})

	t.Run("compiles with empty sources", func(t *testing.T) {
		transforms, err := CompileTransforms("", "", "", "")
		require.NoError(t, err)
		assert.Nil(t, transforms.Pre)
		assert.Nil(t, transforms.Mock)
		assert.Nil(t, transforms.PostEach)
		assert.Nil(t, transforms.PostAll)
	})

	t.Run("compiles partial transforms", func(t *testing.T) {
		transforms, err := CompileTransforms(
			`return input`,
			"",
			`return output`,
			"",
		)
		require.NoError(t, err)
		assert.NotNil(t, transforms.Pre)
		assert.Nil(t, transforms.Mock)
		assert.NotNil(t, transforms.PostEach)
		assert.Nil(t, transforms.PostAll)
	})

	t.Run("returns error for invalid pre-transform", func(t *testing.T) {
		_, err := CompileTransforms(`invalid {`, "", "", "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pre-transform")
	})

	t.Run("returns error for invalid mock", func(t *testing.T) {
		_, err := CompileTransforms("", `invalid {`, "", "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mock")
	})

	t.Run("returns error for invalid post.each", func(t *testing.T) {
		_, err := CompileTransforms("", "", `invalid {`, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "post.each")
	})

	t.Run("returns error for invalid post.all", func(t *testing.T) {
		_, err := CompileTransforms("", "", "", `invalid {`)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "post.all")
	})

	t.Run("returns all errors joined", func(t *testing.T) {
		_, err := CompileTransforms(`invalid {`, `invalid {`, `invalid {`, `invalid {`)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pre-transform")
		assert.Contains(t, err.Error(), "mock")
		assert.Contains(t, err.Error(), "post.each")
		assert.Contains(t, err.Error(), "post.all")
	})
}
