package config

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"

	"github.com/mpyw/sql-http-proxy/internal"
)

// Mock source categories
var (
	objectSources = []string{"object", "object_json", "object_json_file", "object_js"}
	arraySources  = []string{"array", "array_json", "array_json_file", "array_js", "csv", "csv_file", "jsonl", "jsonl_file"}
)

// Compiled regexes for path matching
var (
	queryOrMutationPathRe = regexp.MustCompile(`^/(queries|mutations)/\d+$`)
	mockPathRe            = regexp.MustCompile(`^/(queries|mutations)/\d+/mock$`)
)

// formatValidationError converts a jsonschema ValidationError into a user-friendly message.
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
			if len(oneOf.Subschemas) >= 2 && isQueryOrMutationPath(e.path) {
				return fmt.Sprintf("at '%s': cannot have both 'sql' and 'mock' - use one or the other", formatPath(e.path))
			}
		}
	}
	return ""
}

// findMissingSqlMockError detects when neither sql nor mock is present.
func findMissingSqlMockError(errors []errorWithPath) string {
	for _, e := range errors {
		if oneOf, ok := e.err.ErrorKind.(*kind.OneOf); ok {
			if len(oneOf.Subschemas) == 0 && isQueryOrMutationPath(e.path) {
				if isSqlMockOneOfError(e.err) {
					return fmt.Sprintf("at '%s': must have either 'sql' or 'mock'", formatPath(e.path))
				}
			}
		}
	}
	return ""
}

// findMockSourceError detects mock-related errors and provides helpful messages.
func findMockSourceError(errors []errorWithPath) string {
	// First check for type: none with mock (detected at mutation level, not mock level)
	for _, e := range errors {
		if !isQueryOrMutationPath(e.path) {
			continue
		}

		if addl, ok := e.err.ErrorKind.(*kind.AdditionalProperties); ok {
			// Check if mock is in the additional properties and type is none
			if slices.Contains(addl.Properties, "mock") {
				// Check if this is type: none by looking for const errors
				if isTypeNone(errors, e.path) {
					return fmt.Sprintf("at '%s': type 'none' does not support mock - use 'sql' instead",
						formatPath(e.path))
				}
			}
		}
	}

	// Look for mock path errors
	for _, e := range errors {
		if !isMockPath(e.path) {
			continue
		}

		// Get parent query/mutation path to understand context
		parentPath := getParentPath(e.path)
		typeName := findTypeValue(errors, parentPath)

		// Check for additional properties error (wrong source type)
		if addl, ok := e.err.ErrorKind.(*kind.AdditionalProperties); ok {
			return formatMockSourceMismatchError(parentPath, typeName, addl.Properties)
		}
	}

	// Look for missing filter error (array source with type: one without filter)
	for _, e := range errors {
		if !isMockPath(e.path) {
			continue
		}

		if req, ok := e.err.ErrorKind.(*kind.Required); ok {
			if slices.Contains(req.Missing, "filter") {
				if usedSource := findUsedArraySource(errors, e.path); usedSource != "" {
					parentPath := getParentPath(e.path)
					return fmt.Sprintf("at '%s': type 'one' with '%s' source requires 'filter' to select a single row",
						formatPath(parentPath), usedSource)
				}
			}
		}
	}

	// Look for multiple sources error
	for _, e := range errors {
		if !isMockPath(e.path) {
			continue
		}

		if addl, ok := e.err.ErrorKind.(*kind.AdditionalProperties); ok {
			if len(addl.Properties) > 1 {
				parentPath := getParentPath(e.path)
				return fmt.Sprintf("at '%s': mock must have exactly one source, found multiple: %s",
					formatPath(parentPath), strings.Join(addl.Properties, ", "))
			}
		}
	}

	return ""
}

// formatMockSourceMismatchError creates a helpful message for source/type mismatches.
func formatMockSourceMismatchError(parentPath string, typeName string, invalidSources []string) string {
	if len(invalidSources) == 0 {
		return ""
	}

	source := invalidSources[0]

	// Check if multiple sources
	if len(invalidSources) > 1 {
		return fmt.Sprintf("at '%s': mock must have exactly one source, found: %s",
			formatPath(parentPath), strings.Join(invalidSources, ", "))
	}

	// type: one with array source (without filter - but we should have caught this above)
	if typeName == "one" && isArraySource(source) {
		return fmt.Sprintf("at '%s': type 'one' with '%s' requires 'filter' to select a single row, or use object/object_js for a single object",
			formatPath(parentPath), source)
	}

	// type: many with object source
	if typeName == "many" && isObjectSource(source) {
		return fmt.Sprintf("at '%s': type 'many' does not support '%s' - use array, csv, or jsonl sources",
			formatPath(parentPath), source)
	}

	// type: none with mock
	if typeName == "none" {
		return fmt.Sprintf("at '%s': type 'none' does not support mock - use 'sql' instead",
			formatPath(parentPath))
	}

	// Generic message
	return fmt.Sprintf("at '%s': invalid mock source '%s' for type '%s'",
		formatPath(parentPath), source, typeName)
}

// findTypeValue extracts the type value from errors at the given query/mutation path.
func findTypeValue(errors []errorWithPath, parentPath string) string {
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

// isTypeNone checks if the type at the given path is 'none'.
// This is determined by looking for const errors expecting 'one' or 'many' (meaning the actual is 'none').
func isTypeNone(errors []errorWithPath, parentPath string) bool {
	typePath := parentPath + "/type"

	// Count how many different const errors we have for the type
	// If we see both 'one' and 'many' as expected values, the actual must be 'none'
	foundOne := false
	foundMany := false

	for _, e := range errors {
		if e.path == typePath {
			if constErr, ok := e.err.ErrorKind.(*kind.Const); ok {
				if s, ok := constErr.Want.(string); ok {
					if s == "one" {
						foundOne = true
					}
					if s == "many" {
						foundMany = true
					}
				}
			}
		}
	}

	return foundOne && foundMany
}

// findUsedArraySource finds which array source was used from error messages.
func findUsedArraySource(errors []errorWithPath, mockPath string) string {
	for _, e := range errors {
		if e.path == mockPath {
			if addl, ok := e.err.ErrorKind.(*kind.AdditionalProperties); ok {
				for _, prop := range addl.Properties {
					if isArraySource(prop) {
						return prop
					}
				}
			}
		}
	}
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

func isQueryOrMutationPath(path string) bool {
	return queryOrMutationPathRe.MatchString(path)
}

func isMockPath(path string) bool {
	return mockPathRe.MatchString(path)
}

func getParentPath(path string) string {
	if lastSlash := strings.LastIndex(path, "/"); lastSlash > 0 {
		return path[:lastSlash]
	}
	return path
}

func isObjectSource(source string) bool {
	return slices.Contains(objectSources, source)
}

func isArraySource(source string) bool {
	return slices.Contains(arraySources, source)
}

func formatPath(path string) string {
	if path == "" || path == "/" {
		return "(root)"
	}

	path = strings.TrimPrefix(path, "/")
	parts := strings.Split(path, "/")
	var result strings.Builder

	for i, part := range parts {
		if isNumeric(part) {
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

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
