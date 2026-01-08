// Package js provides JavaScript-based data transformation using goja runtime.
package js

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

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
