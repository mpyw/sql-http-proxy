package mock

import (
	"strconv"
	"strings"
)

// parseCSVValue converts a string value to an appropriate Go type.
// Conversion rules (similar to PHP):
//   - Numeric strings → float64 (PHP is_numeric style)
//   - Boolean strings → bool (PHP filter_var FILTER_VALIDATE_BOOLEAN style)
//   - "null" (case-insensitive) → nil
//   - Everything else → string
//
// Shared on purpose: csv.go applies it to every cell when no custom
// value_parser is configured.
//
//declscope:package
func parseCSVValue(s string) any {
	// Check for null
	if strings.EqualFold(s, "null") {
		return nil
	}

	// Check for boolean (PHP filter_var FILTER_VALIDATE_BOOLEAN style)
	if b, ok := boolFromCSVValue(s); ok {
		return b
	}

	// Check for numeric (PHP is_numeric style)
	if isNumericCSVValue(s) {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}

	// Default to string
	return s
}

// isNumericCSVValue checks if a string is numeric (PHP is_numeric style).
// Accepts: integers, floats, exponential notation, leading +/-
// Examples: "1", "1.5", "-1", "+1.5", "1e10", "1E-5", ".5", "-.5"
func isNumericCSVValue(s string) bool {
	if s == "" {
		return false
	}

	// Remove leading + or -
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}

	if s == "" {
		return false
	}

	// Check for valid float format
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// boolFromCSVValue attempts to parse a string as a boolean.
// PHP filter_var FILTER_VALIDATE_BOOLEAN style:
//   - true: "true", "1", "yes", "on" (case-insensitive)
//   - false: "false", "0", "no", "off", "" (case-insensitive)
//
// Returns (value, ok) where ok is false if not a recognized boolean string.
func boolFromCSVValue(s string) (bool, bool) {
	lower := strings.ToLower(s)
	switch lower {
	case "true", "1", "yes", "on":
		return true, true
	case "false", "0", "no", "off", "":
		return false, true
	default:
		return false, false
	}
}
