package mock

import (
	"github.com/dop251/goja"
	"github.com/samber/lo"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// JSFilteredSource wraps a Source and filters data using JavaScript.
type JSFilteredSource struct {
	source  Source
	program *goja.Program
	helpers *js.CompiledHelpers
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

	return &JSFilteredSource{
		source:  source,
		program: program,
		helpers: helpers,
	}, nil
}

// Data returns filtered data based on JavaScript filter function.
func (f *JSFilteredSource) Data(ctx map[string]any, sql string, input map[string]any) (any, map[string]any, error) {
	data, newCtx, err := f.source.Data(ctx, sql, input)
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
func (f *JSFilteredSource) evaluateFilter(row any, input map[string]any, ctx map[string]any) (bool, error) {
	vm := goja.New()

	// Set up helpers if available
	if f.helpers != nil {
		if err := f.helpers.InjectInto(vm); err != nil {
			return false, err
		}
	}

	// Set ctx as a free variable
	if ctx == nil {
		ctx = make(map[string]any)
	}
	if err := vm.Set("ctx", ctx); err != nil {
		return false, err
	}

	// Run the program to get the filter function
	val, err := vm.RunProgram(f.program)
	if err != nil {
		return false, err
	}

	// Get the function
	fn, ok := goja.AssertFunction(val)
	if !ok {
		return false, nil
	}

	// Call the function with row and input
	result, err := fn(goja.Undefined(), vm.ToValue(row), vm.ToValue(input))
	if err != nil {
		return false, err
	}

	// Convert result to boolean
	return result.ToBoolean(), nil
}
