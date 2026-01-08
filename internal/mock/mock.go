package mock

import (
	"fmt"
	"path/filepath"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

// Source represents a mock data source.
type Source interface {
	// Data returns the mock data.
	// For JS: executes the function with ctx, sql, input
	// For static sources: returns the pre-parsed data (ctx unchanged)
	Data(ctx map[string]any, sql string, input map[string]any) (any, map[string]any, error)
}

// Compile compiles a MockTransform into a Source.
// configDir is the directory of the YAML config file (for resolving relative paths).
func Compile(mt *config.MockTransform, configDir string) (Source, error) {
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
		return filepath.Join(configDir, path)
	}

	switch {
	case mt.JS != "":
		return CompileJS(mt.JS)

	case mt.JSFile != "":
		return CompileJSFile(resolvePath(mt.JSFile))

	case mt.CSV != "":
		return ParseCSV(mt.CSV)

	case mt.CSVFile != "":
		return ParseCSVFile(resolvePath(mt.CSVFile))

	case mt.JSON != nil:
		return NewJSON(mt.JSON), nil

	case mt.JSONFile != "":
		return ParseJSONFile(resolvePath(mt.JSONFile))

	case mt.JSONL != "":
		return ParseJSONL(mt.JSONL)

	case mt.JSONLFile != "":
		return ParseJSONLFile(resolvePath(mt.JSONLFile))

	default:
		return nil, fmt.Errorf("mock: no source specified")
	}
}
