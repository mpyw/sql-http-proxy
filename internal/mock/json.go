package mock

import (
	"bufio"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// jsonSource holds pre-parsed JSON data.
type jsonSource struct {
	data any
}

// newJSON creates a jsonSource from inline YAML/JSON data.
// If data is a string, it is parsed as JSON.
// Otherwise, it is used as-is (already parsed by YAML unmarshaler).
func newJSON(data any) (*jsonSource, error) {
	// If data is a string, parse it as JSON
	if s, ok := data.(string); ok {
		var parsed any
		if err := json.Unmarshal([]byte(s), &parsed); err != nil {
			return nil, err
		}
		return &jsonSource{data: parsed}, nil
	}
	return &jsonSource{data: data}, nil
}

// parseJSONString parses a JSON string into a jsonSource.
func parseJSONString(jsonStr string) (*jsonSource, error) {
	var parsed any
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, err
	}
	return &jsonSource{data: parsed}, nil
}

// parseJSONFile parses a JSON file into a jsonSource.
func parseJSONFile(path string) (*jsonSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Debug("Failed to close file", "error", err)
		}
	}()

	var data any
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}

	return &jsonSource{data: data}, nil
}

// parseJSONL parses inline JSONL (JSON Lines) data.
// Each line is a separate JSON object.
func parseJSONL(data string) (*jsonSource, error) {
	return parseJSONLReader(strings.NewReader(data))
}

// parseJSONLFile parses a JSONL file into a jsonSource.
func parseJSONLFile(path string) (*jsonSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Debug("Failed to close file", "error", err)
		}
	}()
	return parseJSONLReader(f)
}

func parseJSONLReader(r io.Reader) (*jsonSource, error) {
	var rows []any
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Skip empty lines
		}

		var obj any
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			return nil, err
		}
		rows = append(rows, obj)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &jsonSource{data: rows}, nil
}

// Data returns the parsed JSON data.
// ctx, sql, and tc are ignored for static data sources.
func (s *jsonSource) Data(_ map[string]any, _ string, _ map[string]any, _ *js.TransformContext) (any, map[string]any, error) {
	// Return a deep copy to prevent mutation
	return deepCopy(s.data), nil, nil
}

// deepCopy creates a deep copy of JSON-compatible data recursively.
// Supported types: nil, bool, float64, string, []any, map[string]any
func deepCopy(data any) any {
	switch v := data.(type) {
	case nil:
		return nil
	case bool, float64, string:
		// Primitives are immutable, return as-is
		return v
	case int:
		// YAML may parse integers as int
		return v
	case int64:
		// Some JSON parsers use int64
		return v
	case []any:
		if v == nil {
			return nil
		}
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = deepCopy(item)
		}
		return result
	case map[string]any:
		if v == nil {
			return nil
		}
		result := make(map[string]any, len(v))
		for key, val := range v {
			result[key] = deepCopy(val)
		}
		return result
	default:
		// Fallback for unexpected types: return as-is (shallow)
		// This shouldn't happen with JSON/YAML parsed data
		return v
	}
}
