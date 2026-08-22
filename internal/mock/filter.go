package mock

import (
	"github.com/dop251/goja"
	"github.com/samber/lo"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// jsFilteredSource wraps a Source and filters data using JavaScript.
type jsFilteredSource struct {
	source Source
	vms    *js.PooledVM
}

// newJSFilteredSource creates a jsFilteredSource that filters by JavaScript code.
// The filter JS receives: row (param), input (param), ctx (free var).
// It should return a boolean (true to include the row).
func newJSFilteredSource(source Source, filterJS string, helpers *js.CompiledHelpers) (*jsFilteredSource, error) {
	// Compile filter as a function that takes row and input as parameters
	wrapped := "(function(row, input) { " + filterJS + " })"
	program, err := goja.Compile("filter", wrapped, true)
	if err != nil {
		return nil, err
	}

	return &jsFilteredSource{
		source: source,
		vms:    js.NewPooledVM("filter", program, helpers),
	}, nil
}

// Data returns filtered data based on JavaScript filter function.
func (f *jsFilteredSource) Data(ctx map[string]any, sql string, input map[string]any, tc *js.TransformContext) (any, map[string]any, error) {
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
func (f *jsFilteredSource) filterData(data any, input map[string]any, ctx map[string]any) (any, error) {
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
func (f *jsFilteredSource) filterArray(arr []any, input map[string]any, ctx map[string]any) ([]any, error) {
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
// Execution is limited to js.JSTimeout to prevent infinite loops.
func (f *jsFilteredSource) evaluateFilter(row any, input map[string]any, ctx map[string]any) (bool, error) {
	// ctx is a free variable, refreshed per call because the pooled runtime may
	// still hold the previous caller's value.
	if ctx == nil {
		ctx = make(map[string]any)
	}
	return f.vms.Call(map[string]any{"ctx": ctx}, goja.Value.ToBoolean, row, input)
}
