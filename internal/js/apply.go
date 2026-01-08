package js

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/samber/lo"
)

// PreTransformResult holds the result of a pre-transform execution.
type PreTransformResult struct {
	Output map[string]any // Transformed input parameters
	SQL    string         // Modified SQL (or original if unchanged)
	Ctx    map[string]any // Modified context
}

// MockResult holds the result of a mock execution.
type MockResult struct {
	Output any            // Mock response data
	Ctx    map[string]any // Modified context
}

// ApplyPre applies pre-transform with free variables ctx and sql.
// Returns the transformed input, modified SQL, and updated context.
func (t *Transformer) ApplyPre(ctx map[string]any, sql string, input map[string]any) (*PreTransformResult, error) {
	vm := goja.New()

	// Set free variables on globalThis
	globalThis := vm.GlobalObject()
	if err := globalThis.Set("ctx", ctx); err != nil {
		return nil, fmt.Errorf("failed to set ctx: %w", err)
	}
	if err := globalThis.Set("sql", sql); err != nil {
		return nil, fmt.Errorf("failed to set sql: %w", err)
	}

	fn, err := vm.RunProgram(t.program)
	if err != nil {
		return nil, fmt.Errorf("failed to load pre-transform function: %w", err)
	}

	callable, ok := goja.AssertFunction(fn)
	if !ok {
		return nil, fmt.Errorf("pre-transform is not a function")
	}

	result, err := callable(goja.Undefined(), vm.ToValue(input))
	if err != nil {
		return nil, parseJSError(err)
	}

	// Read back free variables from globalThis
	newSQL := sql
	if sqlVal := globalThis.Get("sql"); sqlVal != nil && !goja.IsUndefined(sqlVal) {
		if s, ok := sqlVal.Export().(string); ok {
			newSQL = s
		}
	}

	newCtx := ctx
	if ctxVal := globalThis.Get("ctx"); ctxVal != nil && !goja.IsUndefined(ctxVal) {
		if c, ok := ctxVal.Export().(map[string]any); ok {
			newCtx = c
		}
	}

	// Parse output
	output, ok := result.Export().(map[string]any)
	if !ok {
		return nil, fmt.Errorf("pre-transform must return an object")
	}

	return &PreTransformResult{
		Output: output,
		SQL:    newSQL,
		Ctx:    newCtx,
	}, nil
}

// ApplyMock applies mock transform with free variables ctx and sql.
// Returns the mock response data and updated context.
// Mock can return any data type (object for "one", array for "many").
func (t *Transformer) ApplyMock(ctx map[string]any, sql string, input map[string]any) (*MockResult, error) {
	vm := goja.New()

	// Set free variables on globalThis
	globalThis := vm.GlobalObject()
	if err := globalThis.Set("ctx", ctx); err != nil {
		return nil, fmt.Errorf("failed to set ctx: %w", err)
	}
	if err := globalThis.Set("sql", sql); err != nil {
		return nil, fmt.Errorf("failed to set sql: %w", err)
	}

	fn, err := vm.RunProgram(t.program)
	if err != nil {
		return nil, fmt.Errorf("failed to load mock function: %w", err)
	}

	callable, ok := goja.AssertFunction(fn)
	if !ok {
		return nil, fmt.Errorf("mock is not a function")
	}

	result, err := callable(goja.Undefined(), vm.ToValue(input))
	if err != nil {
		return nil, parseJSError(err)
	}

	// Read back ctx from globalThis
	newCtx := ctx
	if ctxVal := globalThis.Get("ctx"); ctxVal != nil && !goja.IsUndefined(ctxVal) {
		if c, ok := ctxVal.Export().(map[string]any); ok {
			newCtx = c
		}
	}

	return &MockResult{
		Output: result.Export(),
		Ctx:    newCtx,
	}, nil
}

// ApplyToEachRow applies the transformation to each row individually.
// For post-transform with "each" mode: ctx, input (original params), output (each row)
func (t *Transformer) ApplyToEachRow(ctx map[string]any, input map[string]any, rows []map[string]any) ([]any, error) {
	results := make([]any, 0, len(rows))
	for _, row := range rows {
		result, err := t.Apply(ctx, input, row)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// ApplyToAllRows applies the transformation to the entire result array.
// For post-transform with "all" mode: ctx, input (original params), output (entire array)
func (t *Transformer) ApplyToAllRows(ctx map[string]any, input map[string]any, rows []map[string]any) (any, error) {
	return t.Apply(ctx, input, lo.ToAnySlice(rows))
}
