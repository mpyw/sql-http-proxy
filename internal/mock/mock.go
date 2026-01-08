package mock

import (
	"fmt"
	"path/filepath"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/js"
)

// Source represents a mock data source.
type Source interface {
	// Data returns the mock data.
	// For JS: executes the function with ctx, sql, input
	// For static sources: returns the pre-parsed data (ctx unchanged)
	Data(ctx map[string]any, sql string, input map[string]any) (any, map[string]any, error)
}

// CompileOptions contains options for compiling mock sources.
type CompileOptions struct {
	ConfigDir   string
	Helpers     *js.CompiledHelpers
	ValueParser *ValueParser // Global CSV value parser
}

// Compile compiles a Mock into a Source.
func Compile(m *config.Mock, opts CompileOptions) (Source, error) {
	if m == nil {
		return nil, nil
	}

	// Validate mutual exclusivity
	if err := m.Validate(); err != nil {
		return nil, err
	}

	// Helper to resolve file paths
	resolvePath := func(path string) string {
		if filepath.IsAbs(path) {
			return path
		}
		return filepath.Join(opts.ConfigDir, path)
	}

	// CSV parsing options
	csvOpts := ParseCSVOptions{ValueParser: opts.ValueParser}

	var source Source
	var err error
	var filterJS string

	switch {
	// Object sources (type: one only)
	case m.Object != nil:
		source, err = NewJSON(m.Object)

	case m.ObjectJSON != "":
		source, err = ParseJSONString(m.ObjectJSON)

	case m.ObjectJSONFile != "":
		source, err = ParseJSONFile(resolvePath(m.ObjectJSONFile))

	case m.ObjectJS != "":
		src, err := CompileJS(m.ObjectJS)
		if err != nil {
			return nil, err
		}
		src.SetHelpers(opts.Helpers)
		return src, nil // JS handles its own execution

	// Array sources
	case m.Array != nil:
		source, err = NewJSON(m.Array)
		filterJS = m.Filter

	case m.ArrayJSON != "":
		source, err = ParseJSONString(m.ArrayJSON)
		filterJS = m.Filter

	case m.ArrayJSONFile != "":
		source, err = ParseJSONFile(resolvePath(m.ArrayJSONFile))
		filterJS = m.Filter

	case m.ArrayJS != "":
		src, err := CompileJS(m.ArrayJS)
		if err != nil {
			return nil, err
		}
		src.SetHelpers(opts.Helpers)
		// For array_js with filter, wrap with JSFilteredSource
		if m.Filter != "" {
			return NewJSFilteredSource(src, m.Filter, opts.Helpers)
		}
		return src, nil

	case m.CSV != "":
		source, err = ParseCSVWithOptions(m.CSV, csvOpts)
		filterJS = m.Filter

	case m.CSVFile != "":
		source, err = ParseCSVFileWithOptions(resolvePath(m.CSVFile), csvOpts)
		filterJS = m.Filter

	case m.JSONL != "":
		source, err = ParseJSONL(m.JSONL)
		filterJS = m.Filter

	case m.JSONLFile != "":
		source, err = ParseJSONLFile(resolvePath(m.JSONLFile))
		filterJS = m.Filter

	default:
		return nil, fmt.Errorf("mock: no source specified")
	}

	if err != nil {
		return nil, err
	}

	// Wrap with JSFilteredSource if filter is specified
	if filterJS != "" {
		return NewJSFilteredSource(source, filterJS, opts.Helpers)
	}

	return source, nil
}
