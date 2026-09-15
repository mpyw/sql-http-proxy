package js

import (
	"errors"
	"sync"
	"time"

	"github.com/dop251/goja"
)

// JSTimeout is the maximum execution time for JavaScript transforms.
const JSTimeout = 5 * time.Second

// ErrJSTimeout is returned when JS execution exceeds JSTimeout.
var ErrJSTimeout = errors.New("JavaScript execution timeout")

// runWithTimeout runs fn with a watchdog that interrupts vm after JSTimeout,
// translating the interrupt into ErrJSTimeout.
//
// The watchdog is disarmed under a mutex and the interrupt flag is cleared
// before returning. Both matter for reuse: goja.Runtime.Interrupt sets a flag
// that stays set if it lands after the script already finished, and a runtime
// carrying a stale flag aborts the very next call made on it. time.Timer.Stop
// alone does not close that window because it does not wait for an
// already-running timer function.
//
// Shared on purpose: both the pooled hot path (pool.go) and one-shot
// transformer runs (apply.go) bound their scripts with the same watchdog.
//
//declscope:package
func runWithTimeout(vm *goja.Runtime, fn func() (goja.Value, error)) (goja.Value, error) {
	var mu sync.Mutex
	done := false

	timer := time.AfterFunc(JSTimeout, func() {
		mu.Lock()
		defer mu.Unlock()
		if done {
			return
		}
		vm.Interrupt(ErrJSTimeout)
	})

	// Deferred so the runtime is left usable even if fn panics.
	defer func() {
		timer.Stop()
		mu.Lock()
		done = true
		mu.Unlock()
		vm.ClearInterrupt()
	}()

	result, err := fn()
	if err != nil {
		if interrupted, ok := errors.AsType[*goja.InterruptedError](err); ok {
			if timeoutErr, ok := interrupted.Value().(error); ok && errors.Is(timeoutErr, ErrJSTimeout) {
				return nil, ErrJSTimeout
			}
		}
		return nil, err
	}

	return result, nil
}
