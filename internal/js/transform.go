// Package js provides JavaScript-based data transformation using goja runtime.
package js

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

// TransformError represents an error thrown from JavaScript transform.
// Supports Lambda-style format: throw { status: 400, body: "message" } or throw { status: 400, body: { message: "error" } }
type TransformError struct {
	Status int
	Body   any // string or object (will be JSON serialized)
}

func (e *TransformError) Error() string {
	if msg, ok := e.Body.(string); ok {
		return msg
	}
	if obj, ok := e.Body.(map[string]any); ok {
		if msg, ok := obj["message"].(string); ok {
			return msg
		}
	}
	return "transform error"
}

// Transformer applies JavaScript transformations.
type Transformer struct {
	program *goja.Program
}

// Compile creates a new Transformer with the given JavaScript code and argument names.
// Post-transform uses: ctx, input, output as parameters
func Compile(code string, argNames []string) (*Transformer, error) {
	args := strings.Join(argNames, ", ")
	wrapped := fmt.Sprintf(`(function(%s) { %s })`, args, code)

	program, err := goja.Compile("transform", wrapped, true)
	if err != nil {
		return nil, fmt.Errorf("failed to compile transform: %w", err)
	}

	return &Transformer{program: program}, nil
}

// CompilePre creates a pre-transform with free variables (ctx, sql) and parameter (input).
// Free variables ctx and sql can be read/modified in the function body.
func CompilePre(code string) (*Transformer, error) {
	// Wrap code in a function that takes input as parameter
	// ctx and sql are set as global variables before execution
	wrapped := fmt.Sprintf(`(function(input) { %s })`, code)

	program, err := goja.Compile("transform", wrapped, true)
	if err != nil {
		return nil, fmt.Errorf("failed to compile pre-transform: %w", err)
	}

	return &Transformer{program: program}, nil
}

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

// Apply applies the transformation with the given arguments.
// If the JS code throws an error, it returns a TransformError with status and message.
func (t *Transformer) Apply(args ...any) (any, error) {
	vm := goja.New()

	fn, err := vm.RunProgram(t.program)
	if err != nil {
		return nil, fmt.Errorf("failed to load transform function: %w", err)
	}

	callable, ok := goja.AssertFunction(fn)
	if !ok {
		return nil, fmt.Errorf("transform is not a function")
	}

	// Convert args to goja values
	jsArgs := make([]goja.Value, len(args))
	for i, arg := range args {
		jsArgs[i] = vm.ToValue(arg)
	}

	result, err := callable(goja.Undefined(), jsArgs...)
	if err != nil {
		return nil, parseJSError(err)
	}

	return result.Export(), nil
}

// parseJSError extracts status and body from a thrown JS error.
// Supports Lambda-style format: throw { status: 400, body: "message" } or throw { status: 400, body: { message: "error" } }
// Native Error objects (throw new Error("...")) always return 500.
func parseJSError(err error) error {
	var jsErr *goja.Exception
	if !errors.As(err, &jsErr) {
		return fmt.Errorf("transform error: %w", err)
	}

	val := jsErr.Value().Export()

	// Handle thrown object
	if obj, ok := val.(map[string]any); ok {
		_, hasStatus := obj["status"]
		_, hasBody := obj["body"]

		// Native Error object: has "message" but no "status" or "body"
		// Always return 500 for native errors (don't trust user JS)
		if !hasStatus && !hasBody {
			return &TransformError{Status: 500, Body: jsErr.String()}
		}

		// Lambda-style: throw { status: 400, body: ... }
		status := 0
		if s, ok := obj["status"].(int64); ok {
			status = int(s)
		} else if s, ok := obj["status"].(float64); ok {
			status = int(s)
		}

		body, hasBody := obj["body"]
		if !hasBody {
			body = jsErr.String()
		}

		return &TransformError{Status: status, Body: body}
	}

	// Handle thrown string: throw "error message"
	if msg, ok := val.(string); ok {
		return &TransformError{Status: 0, Body: msg}
	}

	// Fallback
	return &TransformError{Status: 0, Body: jsErr.String()}
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
	// Convert []map[string]any to []any for JS compatibility
	rowsAny := make([]any, len(rows))
	for i, row := range rows {
		rowsAny[i] = row
	}
	return t.Apply(ctx, input, rowsAny)
}
