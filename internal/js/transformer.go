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

// CompilePost creates a post-transform with free variable (ctx) and parameters (input, output).
// Free variable ctx can be read/modified in the function body.
func CompilePost(code string) (*Transformer, error) {
	// Wrap code in a function that takes input and output as parameters
	// ctx is set as a global variable before execution
	wrapped := fmt.Sprintf(`(function(input, output) { %s })`, code)

	program, err := goja.Compile("transform", wrapped, true)
	if err != nil {
		return nil, fmt.Errorf("failed to compile post-transform: %w", err)
	}

	return &Transformer{program: program}, nil
}
