package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		// Null
		{"null lowercase", "null", nil},
		{"null uppercase", "NULL", nil},
		{"null mixed case", "Null", nil},

		// Boolean true
		{"true", "true", true},
		{"TRUE", "TRUE", true},
		{"1 as bool", "1", true},
		{"yes", "yes", true},
		{"YES", "YES", true},
		{"on", "on", true},
		{"ON", "ON", true},

		// Boolean false
		{"false", "false", false},
		{"FALSE", "FALSE", false},
		{"0 as bool", "0", false},
		{"no", "no", false},
		{"NO", "NO", false},
		{"off", "off", false},
		{"OFF", "OFF", false},
		{"empty string", "", false},

		// Numeric (integers)
		{"positive integer", "42", float64(42)},
		{"negative integer", "-42", float64(-42)},
		{"explicit positive", "+42", float64(42)},

		// Numeric (floats)
		{"float", "3.14", float64(3.14)},
		{"negative float", "-3.14", float64(-3.14)},
		{"float with leading dot", ".5", float64(0.5)},
		{"negative float with leading dot", "-.5", float64(-0.5)},

		// Numeric (exponential)
		{"exponential lowercase", "1e10", float64(1e10)},
		{"exponential uppercase", "1E10", float64(1e10)},
		{"negative exponential", "1e-5", float64(1e-5)},
		{"float exponential", "1.5e3", float64(1500)},

		// String fallback
		{"plain string", "hello", "hello"},
		{"string with spaces", "hello world", "hello world"},
		{"string starting with number", "123abc", "123abc"},
		{"invalid number", "1.2.3", "1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseValue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Valid numeric
		{"0", true},
		{"1", true},
		{"42", true},
		{"-1", true},
		{"+1", true},
		{"3.14", true},
		{"-3.14", true},
		{"+3.14", true},
		{".5", true},
		{"-.5", true},
		{"+.5", true},
		{"1e10", true},
		{"1E10", true},
		{"1e-5", true},
		{"1.5e3", true},

		// Invalid numeric
		{"", false},
		{"+", false},
		{"-", false},
		{"abc", false},
		{"1.2.3", false},
		{"1a", false},
		{"a1", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isNumeric(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input string
		value bool
		ok    bool
	}{
		// True values
		{"true", true, true},
		{"True", true, true},
		{"TRUE", true, true},
		{"1", true, true},
		{"yes", true, true},
		{"Yes", true, true},
		{"YES", true, true},
		{"on", true, true},
		{"On", true, true},
		{"ON", true, true},

		// False values
		{"false", false, true},
		{"False", false, true},
		{"FALSE", false, true},
		{"0", false, true},
		{"no", false, true},
		{"No", false, true},
		{"NO", false, true},
		{"off", false, true},
		{"Off", false, true},
		{"OFF", false, true},
		{"", false, true},

		// Not boolean
		{"maybe", false, false},
		{"2", false, false},
		{"y", false, false},
		{"n", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			value, ok := parseBool(tt.input)
			assert.Equal(t, tt.ok, ok)
			if ok {
				assert.Equal(t, tt.value, value)
			}
		})
	}
}
