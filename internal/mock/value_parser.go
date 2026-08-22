package mock

import (
	"fmt"

	"github.com/dop251/goja"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// ValueParser parses CSV cell values using custom JavaScript.
type ValueParser struct {
	vms *js.PooledVM
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

	return &ValueParser{vms: js.NewPooledVM("value_parser", program, helpers)}, nil
}

// Parse parses a single cell value using the custom JS.
// Execution is limited to js.JSTimeout to prevent infinite loops.
func (p *ValueParser) Parse(value string) (any, error) {
	return p.vms.Call(nil, goja.Value.Export, value)
}
