// Package js provides JavaScript-based data transformation using goja runtime.
package js

import (
	"fmt"

	"github.com/dop251/goja"
)

// Transformer applies JavaScript transformations.
type Transformer struct {
	program *goja.Program
	helpers *CompiledHelpers
}

// SetHelpers sets the compiled helpers for this transformer.
func (t *Transformer) SetHelpers(h *CompiledHelpers) {
	t.helpers = h
}

// compile is an internal helper that compiles JavaScript code into a Transformer.
// wrapper is the function wrapper format, sourceName is for goja.Compile,
// errorName is for error wrapping (empty string means no wrapping).
func compile(code, wrapper, sourceName, errorName string) (*Transformer, error) {
	wrapped := fmt.Sprintf(wrapper, code)
	program, err := goja.Compile(sourceName, wrapped, true)
	if err != nil {
		if errorName != "" {
			return nil, fmt.Errorf("failed to compile %s: %w", errorName, err)
		}
		return nil, err
	}
	return &Transformer{program: program}, nil
}

// CompilePre creates a pre-transform with free variables (ctx, sql) and parameter (input).
// Free variables ctx and sql can be read/modified in the function body.
func CompilePre(code string) (*Transformer, error) {
	return compile(code, `(function(input) { %s })`, "transform", "pre-transform")
}

// CompilePost creates a post-transform with free variable (ctx) and parameters (input, output).
// Free variable ctx can be read/modified in the function body.
func CompilePost(code string) (*Transformer, error) {
	return compile(code, `(function(input, output) { %s })`, "transform", "post-transform")
}

// CompileFilter creates a filter function with free variable (ctx) and parameters (row, input).
// Returns boolean indicating whether the row should be included.
func CompileFilter(code string) (*Transformer, error) {
	return compile(code, `(function(row, input) { %s })`, "filter", "filter")
}

// CompileMockJS creates a mock JS function with free variables (ctx, sql) and parameter (input).
// Same signature as CompilePre but returns raw error (no wrapping).
func CompileMockJS(code string) (*Transformer, error) {
	return compile(code, `(function(input) { %s })`, "mock", "")
}
