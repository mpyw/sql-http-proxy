package mock

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"
)

// JSONSource holds pre-parsed JSON data.
type JSONSource struct {
	data any
}

// NewJSON creates a JSONSource from inline YAML/JSON data.
// If data is a string, it is parsed as JSON.
// Otherwise, it is used as-is (already parsed by YAML unmarshaler).
func NewJSON(data any) (*JSONSource, error) {
	// If data is a string, parse it as JSON
	if s, ok := data.(string); ok {
		var parsed any
		if err := json.Unmarshal([]byte(s), &parsed); err != nil {
			return nil, err
		}
		return &JSONSource{data: parsed}, nil
	}
	return &JSONSource{data: data}, nil
}

// ParseJSONString parses a JSON string into a JSONSource.
func ParseJSONString(jsonStr string) (*JSONSource, error) {
	var parsed any
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, err
	}
	return &JSONSource{data: parsed}, nil
}

// ParseJSONFile parses a JSON file into a JSONSource.
func ParseJSONFile(path string) (*JSONSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var data any
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}

	return &JSONSource{data: data}, nil
}

// ParseJSONL parses inline JSONL (JSON Lines) data.
// Each line is a separate JSON object.
func ParseJSONL(data string) (*JSONSource, error) {
	return parseJSONLReader(strings.NewReader(data))
}

// ParseJSONLFile parses a JSONL file into a JSONSource.
func ParseJSONLFile(path string) (*JSONSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return parseJSONLReader(f)
}

func parseJSONLReader(r io.Reader) (*JSONSource, error) {
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

	return &JSONSource{data: rows}, nil
}

// Data returns the parsed JSON data.
// ctx and sql are ignored for static data sources.
func (s *JSONSource) Data(_ map[string]any, _ string, _ map[string]any) (any, map[string]any, error) {
	// Return a deep copy to prevent mutation
	return deepCopy(s.data), nil, nil
}

// deepCopy creates a deep copy of the data using JSON marshaling.
func deepCopy(data any) any {
	if data == nil {
		return nil
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return data // Fall back to shallow reference if marshal fails
	}

	var result any
	if err := json.Unmarshal(bytes, &result); err != nil {
		return data // Fall back to shallow reference if unmarshal fails
	}

	return result
}
