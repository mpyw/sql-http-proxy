package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
)

// formatValidationError converts a jsonschema ValidationError into a user-friendly message.
// It focuses on improving sql/mock exclusivity errors while preserving other error messages.
func formatValidationError(err *jsonschema.ValidationError) string {
	// Try to find sql/mock specific errors (our main improvement target)
	if msg := findSqlMockError(err); msg != "" {
		return msg
	}

	// Fall back to default error formatting for all other cases
	return err.Error()
}

// findSqlMockError looks for sql/mock exclusivity errors and formats them nicely.
func findSqlMockError(err *jsonschema.ValidationError) string {
	if err == nil {
		return ""
	}

	// Collect all errors with their paths
	errors := collectErrors(err)

	// Look for sql/mock both present (oneOf with 2 matches at query/mutation level)
	for _, e := range errors {
		if oneOf, ok := e.err.ErrorKind.(*kind.OneOf); ok {
			if len(oneOf.Subschemas) >= 2 && isQueryOrMutationPath(e.path) {
				return fmt.Sprintf("at '%s': cannot have both 'sql' and 'mock' - use one or the other", formatPath(e.path))
			}
		}
	}

	// Look for missing sql/mock (oneOf with 0 matches at query/mutation level)
	// But only if the error specifically mentions sql or mock
	for _, e := range errors {
		if oneOf, ok := e.err.ErrorKind.(*kind.OneOf); ok {
			if len(oneOf.Subschemas) == 0 && isQueryOrMutationPath(e.path) {
				// Check if this is really about sql/mock by looking at direct children
				if isSqlMockOneOfError(e.err) {
					return fmt.Sprintf("at '%s': must have either 'sql' or 'mock'", formatPath(e.path))
				}
			}
		}
	}

	return ""
}

// isSqlMockOneOfError checks if the oneOf error is specifically about sql/mock.
// This is determined by looking at the direct causes for missing sql/mock properties.
func isSqlMockOneOfError(err *jsonschema.ValidationError) bool {
	if err == nil {
		return false
	}

	// Check direct causes for missing sql or mock
	for _, cause := range err.Causes {
		if req, ok := cause.ErrorKind.(*kind.Required); ok {
			for _, missing := range req.Missing {
				if missing == "sql" || missing == "mock" {
					return true
				}
			}
		}
		// Recursively check nested causes (for oneOf within oneOf)
		if isSqlMockOneOfError(cause) {
			return true
		}
	}
	return false
}

// errorWithPath holds an error and its path for easier processing.
type errorWithPath struct {
	path string
	err  *jsonschema.ValidationError
}

// collectErrors flattens the error tree into a list.
func collectErrors(err *jsonschema.ValidationError) []errorWithPath {
	var result []errorWithPath
	collectErrorsRecursive(err, &result)
	return result
}

func collectErrorsRecursive(err *jsonschema.ValidationError, result *[]errorWithPath) {
	if err == nil {
		return
	}

	path := "/" + strings.Join(err.InstanceLocation, "/")
	*result = append(*result, errorWithPath{path: path, err: err})

	for _, cause := range err.Causes {
		collectErrorsRecursive(cause, result)
	}
}

// isQueryOrMutationPath checks if path matches /queries/N or /mutations/N.
func isQueryOrMutationPath(path string) bool {
	matched, _ := regexp.MatchString(`^/(queries|mutations)/\d+$`, path)
	return matched
}

// formatPath converts JSON pointer path to a more readable format.
// e.g., "/queries/0/mock" -> "queries[0].mock"
func formatPath(path string) string {
	if path == "" || path == "/" {
		return "(root)"
	}

	// Remove leading slash
	path = strings.TrimPrefix(path, "/")

	// Split by /
	parts := strings.Split(path, "/")
	var result strings.Builder

	for i, part := range parts {
		// Check if part is a number (array index)
		if isNumeric(part) {
			result.WriteString("[")
			result.WriteString(part)
			result.WriteString("]")
		} else {
			// Add dot separator if not first element
			if i > 0 {
				result.WriteString(".")
			}
			result.WriteString(part)
		}
	}

	return result.String()
}

// isNumeric checks if a string is a numeric value.
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
