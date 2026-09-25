package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"

	"github.com/mpyw/sql-http-proxy/internal"
)

// Compiled regexes for path matching
var (
	queryOrMutationErrorPathRe = regexp.MustCompile(`^/(queries|mutations)/\d+$`)
	mockErrorPathRe            = regexp.MustCompile(`^/(queries|mutations)/\d+/mock$`)
)

// formatValidationError converts a jsonschema ValidationError into a user-friendly message.
// It is the one entry point config.go (parse) takes into this unit.
//
//declscope:package
func formatValidationError(err *jsonschema.ValidationError) string {
	// Collect all errors
	errors := collectErrors(err)

	for _, finder := range internal.SliceOf(
		// Priority 1: sql/mock both present
		findBothSqlMockError,
		// Priority 2: mock source errors (type mismatch, missing filter, multiple sources)
		findMockSourceError,
		// Priority 3: missing sql/mock
		findMissingSqlMockError,
	) {
		if msg := finder(errors); msg != "" {
			return msg
		}
	}

	// Fall back to default
	return err.Error()
}

// findBothSqlMockError detects when both sql and mock are present.
func findBothSqlMockError(errors []errorWithPath) string {
	for _, e := range errors {
		if oneOf, ok := e.err.ErrorKind.(*kind.OneOf); ok {
			if len(oneOf.Subschemas) >= 2 && isQueryOrMutationErrorPath(e.path) {
				return fmt.Sprintf("at '%s': cannot have both 'sql' and 'mock' - use one or the other", formatErrorPath(e.path))
			}
		}
	}
	return ""
}

// findMissingSqlMockError detects when neither sql nor mock is present.
func findMissingSqlMockError(errors []errorWithPath) string {
	for _, e := range errors {
		if oneOf, ok := e.err.ErrorKind.(*kind.OneOf); ok {
			if len(oneOf.Subschemas) == 0 && isQueryOrMutationErrorPath(e.path) {
				if isSqlMockOneOfError(e.err) {
					return fmt.Sprintf("at '%s': must have either 'sql' or 'mock'", formatErrorPath(e.path))
				}
			}
		}
	}
	return ""
}

// findMockSourceError detects mock-related errors and provides helpful messages.
func findMockSourceError(errors []errorWithPath) string {
	// Look for mock path errors
	for _, e := range errors {
		if !isMockErrorPath(e.path) {
			continue
		}

		// Get parent query/mutation path to understand context
		parentPath := parentErrorPath(e.path)
		typeName := findTypeValueInErrors(errors, parentPath)

		// Check for additional properties error (wrong source type)
		if addl, ok := e.err.ErrorKind.(*kind.AdditionalProperties); ok {
			return formatMockSourceMismatchError(parentPath, typeName, addl.Properties)
		}
	}

	return ""
}

// formatMockSourceMismatchError creates a helpful message for source/type mismatches.
func formatMockSourceMismatchError(parentPath string, typeName string, invalidSources []string) string {
	source := invalidSources[0]

	// Check if multiple sources
	if len(invalidSources) > 1 {
		return fmt.Sprintf("at '%s': mock must have exactly one source, found: %s",
			formatErrorPath(parentPath), strings.Join(invalidSources, ", "))
	}

	// type: one with array source without filter
	if typeName == "one" && isArraySourceKey(source) {
		return fmt.Sprintf("at '%s': type 'one' with '%s' requires 'filter' to select a single row, or use object/object_js for a single object",
			formatErrorPath(parentPath), source)
	}

	// type: many with object source
	if typeName == "many" && isObjectSourceKey(source) {
		return fmt.Sprintf("at '%s': type 'many' does not support '%s' - use array, csv, or jsonl sources",
			formatErrorPath(parentPath), source)
	}

	// Generic message for other cases
	return fmt.Sprintf("at '%s': invalid mock source '%s' for type '%s'",
		formatErrorPath(parentPath), source, typeName)
}

// findTypeValueInErrors extracts the type value from errors at the given query/mutation path.
func findTypeValueInErrors(errors []errorWithPath, parentPath string) string {
	typePath := parentPath + "/type"

	// Look for const errors at the type path to infer what type was expected
	for _, e := range errors {
		if e.path == typePath {
			if constErr, ok := e.err.ErrorKind.(*kind.Const); ok {
				// The expected value tells us what type is NOT the current one
				// If "value must be 'many'" then current type is 'one' (or invalid)
				if s, ok := constErr.Want.(string); ok {
					// Return the opposite
					if s == "many" {
						return "one"
					}
					if s == "one" {
						return "many"
					}
				}
			}
		}
	}

	// Default - couldn't determine
	return ""
}

// isSqlMockOneOfError checks if the oneOf error is specifically about sql/mock.
func isSqlMockOneOfError(err *jsonschema.ValidationError) bool {
	if err == nil {
		return false
	}

	for _, cause := range err.Causes {
		if req, ok := cause.ErrorKind.(*kind.Required); ok {
			for _, missing := range req.Missing {
				if missing == "sql" || missing == "mock" {
					return true
				}
			}
		}
		if isSqlMockOneOfError(cause) {
			return true
		}
	}
	return false
}

// Helper types and functions

type errorWithPath struct {
	path string
	err  *jsonschema.ValidationError
}

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

func isQueryOrMutationErrorPath(path string) bool {
	return queryOrMutationErrorPathRe.MatchString(path)
}

func isMockErrorPath(path string) bool {
	return mockErrorPathRe.MatchString(path)
}

func parentErrorPath(path string) string {
	if parent, _, found := strings.CutLast(path, "/"); found && parent != "" {
		return parent
	}
	return path
}

func formatErrorPath(path string) string {
	if path == "" || path == "/" {
		return "(root)"
	}

	path = strings.TrimPrefix(path, "/")
	parts := strings.Split(path, "/")
	var result strings.Builder

	for i, part := range parts {
		if isErrorPathIndex(part) {
			result.WriteString("[")
			result.WriteString(part)
			result.WriteString("]")
		} else {
			if i > 0 {
				result.WriteString(".")
			}
			result.WriteString(part)
		}
	}

	return result.String()
}

func isErrorPathIndex(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
