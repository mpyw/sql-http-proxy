package js

import (
	"fmt"
	"net/http"

	"github.com/dop251/goja"
	"github.com/samber/lo"
)

// PreTransformResult holds the result of a pre-transform execution.
type PreTransformResult struct {
	Output map[string]any // Transformed input parameters
	SQL    string         // Modified SQL (or original if unchanged)
	Ctx    map[string]any // Modified context
}

// TransformContext holds HTTP context for transforms.
type TransformContext struct {
	Request  *Request  // HTTP request (read-only in JS)
	Response *Response // HTTP response (writable in post-transform)
}

// NewTransformContext creates a new TransformContext from an HTTP request.
func NewTransformContext(r *http.Request) *TransformContext {
	var req *Request
	if r != nil {
		req = NewRequest(r)
	}
	return &TransformContext{
		Request:  req,
		Response: NewResponse(),
	}
}

// vmSetupOptions configures what globals to set on the VM.
type vmSetupOptions struct {
	sql             string
	includeResponse bool
}

// setupVM creates a new goja VM, injects helpers, and sets common globals.
func (t *Transformer) setupVM(ctx map[string]any, tc *TransformContext, opts vmSetupOptions) (*goja.Runtime, *goja.Object, error) {
	vm := goja.New()

	if err := t.helpers.InjectInto(vm); err != nil {
		return nil, nil, fmt.Errorf("failed to inject helpers: %w", err)
	}

	globalThis := vm.GlobalObject()

	for _, entry := range []lo.Tuple3[string, bool, func() any]{
		lo.T3("ctx", true, func() any {
			return ctx
		}),
		lo.T3("sql", opts.sql != "", func() any {
			return opts.sql
		}),
		lo.T3("request", tc != nil && tc.Request != nil, func() any {
			return tc.Request.ToJSObject(vm)
		}),
		lo.T3("response", opts.includeResponse && tc != nil && tc.Response != nil, func() any {
			return tc.Response.ToJSObject(vm)
		}),
	} {
		name, shouldSet, getValue := entry.Unpack()
		if !shouldSet {
			continue
		}
		if err := globalThis.Set(name, getValue()); err != nil {
			return nil, nil, fmt.Errorf("failed to set %s: %w", name, err)
		}
	}

	return vm, globalThis, nil
}

// runCallable runs the compiled program and calls it with the given arguments.
func (t *Transformer) runCallable(vm *goja.Runtime, fnName string, args ...goja.Value) (goja.Value, error) {
	fn, err := vm.RunProgram(t.program)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s function: %w", fnName, err)
	}

	callable, ok := goja.AssertFunction(fn)
	if !ok {
		return nil, fmt.Errorf("%s is not a function", fnName)
	}

	result, err := callable(goja.Undefined(), args...)
	if err != nil {
		return nil, parseJSError(err)
	}

	return result, nil
}

// readBackCtx reads ctx from globalThis, returning fallback if not found.
func readBackCtx(globalThis *goja.Object, fallback map[string]any) map[string]any {
	if ctxVal := globalThis.Get("ctx"); ctxVal != nil && !goja.IsUndefined(ctxVal) {
		if c, ok := ctxVal.Export().(map[string]any); ok {
			return c
		}
	}
	return fallback
}

// readBackSQL reads sql from globalThis, returning fallback if not found.
func readBackSQL(globalThis *goja.Object, fallback string) string {
	if sqlVal := globalThis.Get("sql"); sqlVal != nil && !goja.IsUndefined(sqlVal) {
		if s, ok := sqlVal.Export().(string); ok {
			return s
		}
	}
	return fallback
}

// MockResult holds the result of a mock execution.
type MockResult struct {
	Output any            // Mock response data
	Ctx    map[string]any // Modified context
}

// ApplyPre applies pre-transform with free variables ctx and sql.
// Returns the transformed input, modified SQL, and updated context.
// tc is optional and provides HTTP request context.
func (t *Transformer) ApplyPre(ctx map[string]any, sql string, input map[string]any, tc *TransformContext) (*PreTransformResult, error) {
	vm, globalThis, err := t.setupVM(ctx, tc, vmSetupOptions{sql: sql})
	if err != nil {
		return nil, err
	}

	result, err := t.runCallable(vm, "pre-transform", vm.ToValue(input))
	if err != nil {
		return nil, err
	}

	output, ok := result.Export().(map[string]any)
	if !ok {
		return nil, fmt.Errorf("pre-transform must return an object")
	}

	return &PreTransformResult{
		Output: output,
		SQL:    readBackSQL(globalThis, sql),
		Ctx:    readBackCtx(globalThis, ctx),
	}, nil
}

// ApplyMock applies mock transform with free variables ctx and sql.
// Returns the mock response data and updated context.
// Mock can return any data type (object for "one", array for "many").
// tc is optional and provides HTTP request context.
func (t *Transformer) ApplyMock(ctx map[string]any, sql string, input map[string]any, tc *TransformContext) (*MockResult, error) {
	vm, globalThis, err := t.setupVM(ctx, tc, vmSetupOptions{sql: sql})
	if err != nil {
		return nil, err
	}

	result, err := t.runCallable(vm, "mock", vm.ToValue(input))
	if err != nil {
		return nil, err
	}

	return &MockResult{
		Output: result.Export(),
		Ctx:    readBackCtx(globalThis, ctx),
	}, nil
}

// PostTransformResult holds the result of a post-transform execution.
type PostTransformResult struct {
	Output any            // Transformed output
	Ctx    map[string]any // Modified context
}

// ApplyPost applies post-transform with free variable ctx.
// Returns the transformed output and updated context.
// tc is optional and provides HTTP request/response context.
func (t *Transformer) ApplyPost(ctx map[string]any, input map[string]any, output any, tc *TransformContext) (*PostTransformResult, error) {
	vm, globalThis, err := t.setupVM(ctx, tc, vmSetupOptions{includeResponse: true})
	if err != nil {
		return nil, err
	}

	result, err := t.runCallable(vm, "post-transform", vm.ToValue(input), vm.ToValue(output))
	if err != nil {
		return nil, err
	}

	return &PostTransformResult{
		Output: result.Export(),
		Ctx:    readBackCtx(globalThis, ctx),
	}, nil
}

// ApplyPostToEachRow applies post-transform to each row individually.
// For post-transform with "each" mode: ctx as free variable, input and output (each row) as parameters.
// tc is optional and provides HTTP request/response context.
func (t *Transformer) ApplyPostToEachRow(ctx map[string]any, input map[string]any, rows []map[string]any, tc *TransformContext) ([]any, map[string]any, error) {
	results := make([]any, 0, len(rows))
	currentCtx := ctx
	for _, row := range rows {
		result, err := t.ApplyPost(currentCtx, input, row, tc)
		if err != nil {
			return nil, nil, err
		}
		results = append(results, result.Output)
		currentCtx = result.Ctx
	}
	return results, currentCtx, nil
}

// ApplyPostToAllRows applies post-transform to the entire result array.
// For post-transform with "all" mode: ctx as free variable, input and output (entire array) as parameters.
// tc is optional and provides HTTP request/response context.
func (t *Transformer) ApplyPostToAllRows(ctx map[string]any, input map[string]any, rows any, tc *TransformContext) (*PostTransformResult, error) {
	return t.ApplyPost(ctx, input, rows, tc)
}
