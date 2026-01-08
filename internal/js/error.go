package js

import (
	"errors"
	"fmt"

	"github.com/dop251/goja"
)

// TransformError represents an error thrown from JavaScript transform.
// Supports Lambda-style format: throw { status: 400, body: "message" } or throw { status: 400, body: { message: "error" } }
type TransformError struct {
	Status int
	Body   any // string or object (will be JSON serialized)
}

func (e *TransformError) Error() string {
	if msg, ok := e.Body.(string); ok {
		return msg
	}
	if obj, ok := e.Body.(map[string]any); ok {
		if msg, ok := obj["message"].(string); ok {
			return msg
		}
	}
	return "transform error"
}

// parseJSError extracts status and body from a thrown JS error.
// Supports Lambda-style format: throw { status: 400, body: "message" } or throw { status: 400, body: { message: "error" } }
// Native Error objects (throw new Error("...")) always return 500.
func parseJSError(err error) error {
	var jsErr *goja.Exception
	if !errors.As(err, &jsErr) {
		return fmt.Errorf("transform error: %w", err)
	}

	val := jsErr.Value().Export()

	// Handle thrown object
	if obj, ok := val.(map[string]any); ok {
		_, hasStatus := obj["status"]
		_, hasBody := obj["body"]

		// Native Error object: has "message" but no "status" or "body"
		// Always return 500 for native errors (don't trust user JS)
		if !hasStatus && !hasBody {
			return &TransformError{Status: 500, Body: jsErr.String()}
		}

		// Lambda-style: throw { status: 400, body: ... }
		status := 0
		if s, ok := obj["status"].(int64); ok {
			status = int(s)
		} else if s, ok := obj["status"].(float64); ok {
			status = int(s)
		}

		body, hasBody := obj["body"]
		if !hasBody {
			body = jsErr.String()
		}

		return &TransformError{Status: status, Body: body}
	}

	// Handle thrown string: throw "error message"
	if msg, ok := val.(string); ok {
		return &TransformError{Status: 0, Body: msg}
	}

	// Fallback
	return &TransformError{Status: 0, Body: jsErr.String()}
}
