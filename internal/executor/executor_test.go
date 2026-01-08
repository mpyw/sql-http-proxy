package executor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapPreError(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		result := WrapPreError(nil)
		require.Nil(t, result)
	})

	t.Run("non-nil error wraps with phase", func(t *testing.T) {
		err := errors.New("test error")
		result := WrapPreError(err)
		require.NotNil(t, result)

		var phaseErr *PhaseError
		require.ErrorAs(t, result, &phaseErr)
		require.Equal(t, "pre", phaseErr.Phase)
		require.Equal(t, "test error", phaseErr.Error())
		require.ErrorIs(t, phaseErr.Unwrap(), err)
	})
}

func TestWrapMockError(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		result := WrapMockError(nil)
		require.Nil(t, result)
	})

	t.Run("non-nil error wraps with phase", func(t *testing.T) {
		err := errors.New("test error")
		result := WrapMockError(err)
		require.NotNil(t, result)

		var phaseErr *PhaseError
		require.ErrorAs(t, result, &phaseErr)
		require.Equal(t, "mock", phaseErr.Phase)
	})
}

func TestWrapPostError(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		result := WrapPostError(nil)
		require.Nil(t, result)
	})

	t.Run("non-nil error wraps with phase", func(t *testing.T) {
		err := errors.New("test error")
		result := WrapPostError(err)
		require.NotNil(t, result)

		var phaseErr *PhaseError
		require.ErrorAs(t, result, &phaseErr)
		require.Equal(t, "post", phaseErr.Phase)
	})
}

func TestPhaseError(t *testing.T) {
	err := errors.New("underlying error")
	phaseErr := &PhaseError{Phase: "test", Err: err}

	t.Run("Error returns underlying error message", func(t *testing.T) {
		require.Equal(t, "underlying error", phaseErr.Error())
	})

	t.Run("Unwrap returns underlying error", func(t *testing.T) {
		require.Equal(t, err, phaseErr.Unwrap())
	})
}
