package mock

import (
	"fmt"

	"github.com/samber/lo"
)

// FilteredSource wraps a Source and filters data based on input parameters.
type FilteredSource struct {
	source Source
}

// NewFilteredSource creates a FilteredSource that filters by input parameters.
func NewFilteredSource(source Source) *FilteredSource {
	return &FilteredSource{source: source}
}

// Data returns filtered data based on input parameters.
// For arrays: returns rows where all input keys match.
// For objects: returns the object if all input keys match, nil otherwise.
func (f *FilteredSource) Data(ctx map[string]any, sql string, input map[string]any) (any, map[string]any, error) {
	data, newCtx, err := f.source.Data(ctx, sql, input)
	if err != nil {
		return nil, nil, err
	}

	filtered := filterByInput(data, input)
	return filtered, newCtx, nil
}

// filterByInput filters data based on input parameters.
func filterByInput(data any, input map[string]any) any {
	if data == nil || len(input) == 0 {
		return data
	}

	switch v := data.(type) {
	case []any:
		return filterArray(v, input)
	case []map[string]any:
		arr := lo.Map(v, func(m map[string]any, _ int) any { return m })
		return filterArray(arr, input)
	case map[string]any:
		if matchesInput(v, input) {
			return v
		}
		return nil
	default:
		return data
	}
}

// filterArray filters an array of objects based on input parameters.
func filterArray(arr []any, input map[string]any) []any {
	result := make([]any, 0)
	for _, item := range arr {
		if m, ok := item.(map[string]any); ok {
			if matchesInput(m, input) {
				result = append(result, item)
			}
		}
	}
	return result
}

// matchesInput checks if a row matches all input parameters.
func matchesInput(row map[string]any, input map[string]any) bool {
	for key, inputVal := range input {
		rowVal, exists := row[key]
		if !exists {
			continue // Skip keys that don't exist in the row
		}
		if !valuesEqual(rowVal, inputVal) {
			return false
		}
	}
	return true
}

// valuesEqual compares two values for equality.
// Handles type coercion for common cases (e.g., string "1" == int 1).
func valuesEqual(a, b any) bool {
	// Direct equality
	if a == b {
		return true
	}

	// Try string comparison (input values are often strings)
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	return aStr == bStr
}
