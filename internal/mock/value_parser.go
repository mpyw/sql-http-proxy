package mock

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dop251/goja"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// ValueParser parses CSV cell values using custom JavaScript.
type ValueParser struct {
	program *goja.Program
	helpers *js.CompiledHelpers
}

// CompileValueParser compiles a custom value parser.
// jsCode is inline JS, jsFile is path to JS file (relative to configDir).
// One of jsCode or jsFile must be non-empty.
func CompileValueParser(jsCode, jsFile, configDir string, helpers *js.CompiledHelpers) (*ValueParser, error) {
	code := jsCode
	if jsFile != "" {
		path := jsFile
		if !filepath.IsAbs(path) {
			path = filepath.Join(configDir, path)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read value_parser file %s: %w", jsFile, err)
		}
		code = string(content)
	}

	if code == "" {
		return nil, fmt.Errorf("value_parser: no code provided")
	}

	// Wrap in function that takes value parameter
	wrapped := fmt.Sprintf(`(function(value) { %s })`, code)
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
