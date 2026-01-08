package mock

import (
	"errors"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/samber/lo"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// JSFilteredSource wraps a Source and filters data using JavaScript.
type JSFilteredSource struct {
	source  Source
	program *goja.Program
	helpers *js.CompiledHelpers
	pool    sync.Pool
}

// filterVM holds a cached VM and callable for reuse.
type filterVM struct {
	vm       *goja.Runtime
	callable goja.Callable
}

// NewJSFilteredSource creates a JSFilteredSource that filters by JavaScript code.
// The filter JS receives: row (param), input (param), ctx (free var).
// It should return a boolean (true to include the row).
func NewJSFilteredSource(source Source, filterJS string, helpers *js.CompiledHelpers) (*JSFilteredSource, error) {
	// Compile filter as a function that takes row and input as parameters
	wrapped := "(function(row, input) { " + filterJS + " })"
	program, err := goja.Compile("filter", wrapped, true)
	if err != nil {
		return nil, err
	}

	f := &JSFilteredSource{
		source:  source,
		program: program,
		helpers: helpers,
	}
	f.pool.New = func() any {
		vm, callable, err := f.createVM()
		if err != nil {
			return nil
		}
		return &filterVM{vm: vm, callable: callable}
	}
	return f, nil
}

// createVM creates a new VM with helpers and callable.
func (f *JSFilteredSource) createVM() (*goja.Runtime, goja.Callable, error) {
	vm := goja.New()

	if f.helpers != nil {
		if err := f.helpers.InjectInto(vm); err != nil {
			return nil, nil, err
		}
	}

	val, err := vm.RunProgram(f.program)
	if err != nil {
		return nil, nil, err
	}

	callable, ok := goja.AssertFunction(val)
	if !ok {
		return nil, nil, nil
	}

	return vm, callable, nil
}

// Data returns filtered data based on JavaScript filter function.
func (f *JSFilteredSource) Data(ctx map[string]any, sql string, input map[string]any, tc *js.TransformContext) (any, map[string]any, error) {
	data, newCtx, err := f.source.Data(ctx, sql, input, tc)
	if err != nil {
		return nil, nil, err
	}

	// If ctx wasn't modified by source, use the original
	if newCtx == nil {
		newCtx = ctx
	}

	filtered, err := f.filterData(data, input, newCtx)
	if err != nil {
		return nil, nil, err
	}

	return filtered, newCtx, nil
}

// filterData filters the data using the JavaScript filter function.
func (f *JSFilteredSource) filterData(data any, input map[string]any, ctx map[string]any) (any, error) {
	if data == nil {
		return data, nil
	}

	switch v := data.(type) {
	case []any:
		return f.filterArray(v, input, ctx)
	case []map[string]any:
		arr := lo.Map(v, func(m map[string]any, _ int) any { return m })
		return f.filterArray(arr, input, ctx)
	default:
		return data, nil
	}
}

// filterArray filters an array of objects using the JavaScript filter function.
func (f *JSFilteredSource) filterArray(arr []any, input map[string]any, ctx map[string]any) ([]any, error) {
	result := make([]any, 0)

	for _, item := range arr {
		include, err := f.evaluateFilter(item, input, ctx)
		if err != nil {
			return nil, err
		}
		if include {
			result = append(result, item)
		}
	}

	return result, nil
}

// evaluateFilter evaluates the filter function for a single row.
// Execution is limited to JSTimeout to prevent infinite loops.
func (f *JSFilteredSource) evaluateFilter(row any, input map[string]any, ctx map[string]any) (bool, error) {
	// Get or create a VM from pool
	fvm := f.pool.Get()
	if fvm == nil {
		// Pool.New failed, create directly
		vm, callable, err := f.createVM()
		if err != nil {
			return false, err
		}
		if callable == nil {
			return false, nil
		}
		fvm = &filterVM{vm: vm, callable: callable}
	}
	cached := fvm.(*filterVM)
	defer f.pool.Put(cached)

	// Set ctx as a free variable (can be updated on reused VM)
	if ctx == nil {
		ctx = make(map[string]any)
	}
	if err := cached.vm.Set("ctx", ctx); err != nil {
		return false, err
	}

	// Set up timeout to prevent infinite loops
	timer := time.AfterFunc(js.JSTimeout, func() {
		cached.vm.Interrupt(js.ErrJSTimeout)
	})
	defer timer.Stop()

	// Call the function with row and input
	result, err := cached.callable(goja.Undefined(), cached.vm.ToValue(row), cached.vm.ToValue(input))
	if err != nil {
		// Check if error was due to timeout interrupt
		var interrupted *goja.InterruptedError
		if errors.As(err, &interrupted) {
			if timeoutErr, ok := interrupted.Value().(error); ok && errors.Is(timeoutErr, js.ErrJSTimeout) {
				return false, js.ErrJSTimeout
			}
		}
		return false, err
	}

	// Convert result to boolean
	return result.ToBoolean(), nil
}
