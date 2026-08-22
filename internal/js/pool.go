package js

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dop251/goja"
)

// PooledVM reuses goja runtimes across calls to a single compiled function.
//
// It exists for the hot paths - CSV value parsing and mock row filtering - that
// invoke the same short function once per row. Building a runtime and injecting
// helpers per row dominates the cost there, so runtimes are pooled instead.
type PooledVM struct {
	name    string
	program *goja.Program
	helpers *CompiledHelpers
	pool    sync.Pool
}

// pooledCallable is a runtime paired with the callable compiled into it.
// A callable is only ever invoked on the runtime it came from.
type pooledCallable struct {
	vm       *goja.Runtime
	callable goja.Callable
}

// NewPooledVM creates a PooledVM for a program that evaluates to a function.
// name is used in error messages.
func NewPooledVM(name string, program *goja.Program, helpers *CompiledHelpers) *PooledVM {
	p := &PooledVM{name: name, program: program, helpers: helpers}
	p.pool.New = func() any {
		// sync.Pool.New cannot report failure. Return nil and let Call surface
		// the real error by building a runtime directly.
		pc, err := p.newCallable()
		if err != nil {
			return nil
		}
		return pc
	}
	return p
}

// newCallable creates a runtime with helpers injected and the program's
// function resolved.
func (p *PooledVM) newCallable() (*pooledCallable, error) {
	vm := goja.New()

	if p.helpers != nil {
		if err := p.helpers.InjectInto(vm); err != nil {
			return nil, fmt.Errorf("failed to inject helpers into %s: %w", p.name, err)
		}
	}

	fn, err := vm.RunProgram(p.program)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", p.name, err)
	}

	callable, ok := goja.AssertFunction(fn)
	if !ok {
		return nil, fmt.Errorf("%s is not a function", p.name)
	}

	return &pooledCallable{vm: vm, callable: callable}, nil
}

// Call invokes the pooled function and converts its result with convert.
//
// convert runs while the runtime is still checked out, because a goja.Value
// belongs to the runtime that produced it and must not outlive its return to
// the pool. That is why the result type is a parameter of this method rather
// than the caller unwrapping a goja.Value afterwards: passing
// goja.Value.ToBoolean or goja.Value.Export pins the conversion to the right
// scope. globals are set on the runtime before the call; pass nil for none.
//
// Execution is bounded by JSTimeout; exceeding it returns ErrJSTimeout.
func (p *PooledVM) Call[T any](globals map[string]any, convert func(goja.Value) T, args ...any) (T, error) {
	var zero T

	pc, ok := p.pool.Get().(*pooledCallable)
	if !ok {
		// pool.New failed; rebuild directly so the error can be reported.
		var err error
		if pc, err = p.newCallable(); err != nil {
			return zero, err
		}
	}

	// A runtime is only pooled again once this call is known to have left it in
	// a good state: a timeout aborted the script mid-execution, and a panic
	// says nothing about what the runtime is holding.
	reusable := false
	defer func() {
		if reusable {
			p.pool.Put(pc)
		}
	}()

	for name, value := range globals {
		if err := pc.vm.Set(name, value); err != nil {
			return zero, fmt.Errorf("failed to set %s in %s: %w", name, p.name, err)
		}
	}

	jsArgs := make([]goja.Value, len(args))
	for i, arg := range args {
		jsArgs[i] = pc.vm.ToValue(arg)
	}

	result, err := runWithTimeout(pc.vm, func() (goja.Value, error) {
		return pc.callable(goja.Undefined(), jsArgs...)
	})
	if err != nil {
		reusable = !errors.Is(err, ErrJSTimeout)
		return zero, err
	}

	reusable = true
	return convert(result), nil
}

// runWithTimeout runs fn with a watchdog that interrupts vm after JSTimeout,
// translating the interrupt into ErrJSTimeout.
//
// The watchdog is disarmed under a mutex and the interrupt flag is cleared
// before returning. Both matter for reuse: goja.Runtime.Interrupt sets a flag
// that stays set if it lands after the script already finished, and a runtime
// carrying a stale flag aborts the very next call made on it. time.Timer.Stop
// alone does not close that window because it does not wait for an
// already-running timer function.
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
