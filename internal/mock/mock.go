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

// Compile compiles a MockTransform into a Source.
func Compile(mt *config.MockTransform, opts CompileOptions) (Source, error) {
	if mt == nil {
		return nil, nil
	}

	// Validate mutual exclusivity
	if err := mt.Validate(); err != nil {
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

	switch {
	case mt.JS != "":
		src, err := CompileJS(mt.JS)
		if err != nil {
			return nil, err
		}
		src.SetHelpers(opts.Helpers)
		return src, nil // JS handles its own filtering

	case mt.CSV != "":
		source, err = ParseCSVWithOptions(mt.CSV, csvOpts)

	case mt.CSVFile != "":
		source, err = ParseCSVFileWithOptions(resolvePath(mt.CSVFile), csvOpts)

	case mt.JSON != nil:
		source, err = NewJSON(mt.JSON)

	case mt.JSONFile != "":
		source, err = ParseJSONFile(resolvePath(mt.JSONFile))

	case mt.JSONL != "":
		source, err = ParseJSONL(mt.JSONL)

	case mt.JSONLFile != "":
		source, err = ParseJSONLFile(resolvePath(mt.JSONLFile))

	default:
		return nil, fmt.Errorf("mock: no source specified")
	}

	if err != nil {
		return nil, err
	}

	// Wrap with FilteredSource if filter_by is specified
	if mt.FilterBy == "input" {
		return NewFilteredSource(source), nil
	}

	return source, nil
}
