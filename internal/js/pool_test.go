package js

import (
	"errors"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestPool(t *testing.T, name, src string) *PooledVM {
	t.Helper()

	program, err := goja.Compile(name, src, true)
	require.NoError(t, err)
	return NewPooledVM(name, program, nil)
}

func TestPooledVM_Call(t *testing.T) {
	t.Run("converts to the type the caller asked for", func(t *testing.T) {
		p := newTestPool(t, "double", `(function(n) { return n * 2 })`)

		n, err := p.Call(nil, goja.Value.ToInteger, 21)
		require.NoError(t, err)
		assert.Equal(t, int64(42), n)

		s, err := p.Call(nil, goja.Value.String, 21)
		require.NoError(t, err)
		assert.Equal(t, "42", s)

		v, err := p.Call(nil, goja.Value.Export, 21)
		require.NoError(t, err)
		assert.Equal(t, int64(42), v)
	})

	t.Run("reads globals set per call", func(t *testing.T) {
		p := newTestPool(t, "readCtx", `(function() { return ctx.n })`)

		// The second call reuses the pooled runtime, so ctx must be refreshed
		// rather than left over from the first call.
		for _, want := range []int64{1, 2, 3} {
			got, err := p.Call(map[string]any{"ctx": map[string]any{"n": want}}, goja.Value.ToInteger)
			require.NoError(t, err)
			assert.Equal(t, want, got)
		}
	})

	t.Run("reports a program that is not a function", func(t *testing.T) {
		p := newTestPool(t, "notFn", `42`)

		_, err := p.Call(nil, goja.Value.Export)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a function")
	})

	t.Run("propagates a thrown error", func(t *testing.T) {
		p := newTestPool(t, "thrower", `(function() { throw new Error("boom") })`)

		_, err := p.Call(nil, goja.Value.Export)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "boom")
	})
}

func TestRunWithTimeout_ClearsLateInterrupt(t *testing.T) {
	// goja.Runtime.Interrupt leaves its flag set if it lands when no script is
	// running, and the flag then aborts the next call made on that runtime.
	// The watchdog can lose that race with a script that finishes right at the
	// deadline, so runWithTimeout must hand back a runtime with a clear flag -
	// otherwise a single near-timeout call would poison a pooled runtime for
	// every request that followed.
	vm := goja.New()
	vm.Interrupt(ErrJSTimeout)

	_, err := runWithTimeout(vm, func() (goja.Value, error) {
		return goja.Undefined(), nil
	})
	require.NoError(t, err)

	v, err := vm.RunString(`1 + 1`)
	require.NoError(t, err, "runtime still carries an interrupt flag")
	assert.Equal(t, int64(2), v.ToInteger())
}

func TestPooledVM_Call_Timeout(t *testing.T) {
	t.Parallel() // JSTimeout is 5s; overlap with the other slow timeout test.

	p := newTestPool(t, "maybeSpin", `(function(spin) { if (spin) { while (true) {} } return 42 })`)

	_, err := p.Call(nil, goja.Value.ToInteger, true)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrJSTimeout), "expected timeout error, got: %v", err)

	// The timed-out runtime is discarded rather than pooled, so later calls
	// must still succeed instead of inheriting the interrupt.
	n, err := p.Call(nil, goja.Value.ToInteger, false)
	require.NoError(t, err)
	assert.Equal(t, int64(42), n)
}
