package body

import (
	json "encoding/json/v2"
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
// Returns empty map if data is empty (common for DELETE/GET requests).
//
// Decoding uses encoding/json/v2 defaults, which follow RFC 7493 (I-JSON):
// duplicate object member names and invalid UTF-8 are rejected rather than
// silently accepted. Rejecting duplicates matters here because the result is
// used directly as SQL bind parameters - under v1, `{"id":1,"id":2}` would
// silently bind the last occurrence, letting a later member override an
// earlier one that a caller (or an intermediary) believed it had set.
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
