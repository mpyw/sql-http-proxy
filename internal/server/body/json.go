package body

import (
	"encoding/json"
	"fmt"
	"io"
)

// parseJSON parses JSON body with optional charset conversion.
// Returns empty map if body is empty (common for DELETE requests).
func parseJSON(body io.Reader, charsetName string) (map[string]any, error) {
	data, err := readBody(body, charsetName)
	if err != nil {
		return nil, err
	}
	return parseJSONBytes(data)
}

// parseJSONBytes parses JSON from bytes.
// Returns empty map if data is empty (common for DELETE requests).
func parseJSONBytes(data []byte) (map[string]any, error) {
	// Empty body is valid - return empty map (common for DELETE/GET requests)
	if len(data) == 0 {
		return map[string]any{}, nil
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON: %v", ErrBadRequest, err)
	}

	return result, nil
}
