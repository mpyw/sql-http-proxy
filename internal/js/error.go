package js

import (
	"errors"
	"fmt"

	"github.com/dop251/goja"
	"github.com/samber/lo"
)

// TransformError represents an error thrown from JavaScript transform.
// Supports Lambda-style format: throw { status: 400, body: "message" } or throw { status: 400, body: { message: "error" } }
type TransformError struct {
	Status int
	Body   any // string or object (will be JSON serialized)
}

func (e *TransformError) Error() string {
	switch body := e.Body.(type) {
	case string:
		return body
	case map[string]any:
		if msg, ok := body["message"].(string); ok {
			return msg
		}
		if msg, ok := body["error"].(string); ok {
			return msg
		}
	}
	return "transform error"
}

// toHTTPStatus converts goja numeric types to an HTTP status code.
// Returns 0 if the value is not int64 or float64.
func toHTTPStatus(v any) int {
	switch n := v.(type) {
	case int64:
		return int(n)
	case float64:
		return int(n)
	default:
		return 0
	}
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
		status, hasStatus := obj["status"]
		body, hasBody := obj["body"]

		// Native Error object: has "message" but no "status" or "body"
		// Always return 500 for native errors (don't trust user JS)
		if !hasStatus && !hasBody {
			return &TransformError{Status: 500, Body: jsErr.String()}
		}

		return &TransformError{
			// Lambda-style: throw { status: 400, body: ... }
			Status: toHTTPStatus(status),
			Body:   lo.Ternary(hasBody, body, any(jsErr.String())),
		}
	}

	// Handle thrown string: throw "error message"
	if msg, ok := val.(string); ok {
		return &TransformError{Status: 0, Body: msg}
	}

	// Fallback
	return &TransformError{Status: 0, Body: jsErr.String()}
}
