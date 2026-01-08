package mock

import (
	"fmt"

	"github.com/dop251/goja"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// ValueParser parses CSV cell values using custom JavaScript.
type ValueParser struct {
	program *goja.Program
	helpers *js.CompiledHelpers
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

	return &ValueParser{program: program, helpers: helpers}, nil
}

// Parse parses a single cell value using the custom JS.
func (p *ValueParser) Parse(value string) (any, error) {
	vm := goja.New()

	// Inject helpers first
	if err := p.helpers.InjectInto(vm); err != nil {
		return nil, fmt.Errorf("failed to inject helpers: %w", err)
	}

	fn, err := vm.RunProgram(p.program)
	if err != nil {
		return nil, fmt.Errorf("failed to load value_parser: %w", err)
	}

	callable, ok := goja.AssertFunction(fn)
	if !ok {
		return nil, fmt.Errorf("value_parser is not a function")
	}

	result, err := callable(goja.Undefined(), vm.ToValue(value))
	if err != nil {
		return nil, err
	}

	return result.Export(), nil
}
