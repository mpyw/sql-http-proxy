package body

import (
	"encoding/json"
	"fmt"
	"io"
)

// parseJSON parses JSON body with optional charset conversion.
func parseJSON(body io.Reader, charsetName string) (map[string]any, error) {
	data, err := readBody(body, charsetName)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON: %v", ErrBadRequest, err)
	}

	return result, nil
}
