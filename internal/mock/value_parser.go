package mock

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dop251/goja"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// ValueParser parses CSV cell values using custom JavaScript.
type ValueParser struct {
	program *goja.Program
	helpers *js.CompiledHelpers
	pool    sync.Pool
}

// valueParserVM holds a cached VM and callable for reuse.
type valueParserVM struct {
	vm       *goja.Runtime
	callable goja.Callable
}

// CompileValueParser compiles a custom value parser from inline JavaScript code.
func CompileValueParser(jsCode string, helpers *js.CompiledHelpers) (*ValueParser, error) {
	if jsCode == "" {
		return nil, fmt.Errorf("value_parser: no code provided")
	}

	// Wrap in function that takes value parameter
	wrapped := fmt.Sprintf(`(function(value) { %s })`, jsCode)
	program, err := goja.Compile("value_parser", wrapped, true)
	if err != nil {
		return nil, fmt.Errorf("failed to compile value_parser: %w", err)
	}

	p := &ValueParser{program: program, helpers: helpers}
	p.pool.New = func() any {
		vm, callable, err := p.createVM()
		if err != nil {
			return nil
		}
		return &valueParserVM{vm: vm, callable: callable}
	}
	return p, nil
}

// createVM creates a new VM with helpers and callable.
func (p *ValueParser) createVM() (*goja.Runtime, goja.Callable, error) {
	vm := goja.New()

	// Inject helpers first
	if p.helpers != nil {
		if err := p.helpers.InjectInto(vm); err != nil {
			return nil, nil, fmt.Errorf("failed to inject helpers: %w", err)
		}
	}

	fn, err := vm.RunProgram(p.program)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load value_parser: %w", err)
	}

	callable, ok := goja.AssertFunction(fn)
	if !ok {
		return nil, nil, fmt.Errorf("value_parser is not a function")
	}

	return vm, callable, nil
}

// Parse parses a single cell value using the custom JS.
// Execution is limited to JSTimeout to prevent infinite loops.
func (p *ValueParser) Parse(value string) (any, error) {
	// Get or create a VM from pool
	vpvm := p.pool.Get()
	if vpvm == nil {
		// Pool.New failed, create directly
		vm, callable, err := p.createVM()
		if err != nil {
			return nil, err
		}
		vpvm = &valueParserVM{vm: vm, callable: callable}
	}
	cached := vpvm.(*valueParserVM)
	defer p.pool.Put(cached)

	// Set up timeout to prevent infinite loops
	timer := time.AfterFunc(js.JSTimeout, func() {
		cached.vm.Interrupt(js.ErrJSTimeout)
	})
	defer timer.Stop()

	result, err := cached.callable(goja.Undefined(), cached.vm.ToValue(value))
	if err != nil {
		// Check if error was due to timeout interrupt
		var interrupted *goja.InterruptedError
		if errors.As(err, &interrupted) {
			if timeoutErr, ok := interrupted.Value().(error); ok && errors.Is(timeoutErr, js.ErrJSTimeout) {
				return nil, js.ErrJSTimeout
			}
		}
		return nil, err
	}

	return result.Export(), nil
}
