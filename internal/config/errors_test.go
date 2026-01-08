package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty path",
			input:    "",
			expected: "(root)",
		},
		{
			name:     "root path",
			input:    "/",
			expected: "(root)",
		},
		{
			name:     "simple property",
			input:    "/dsn",
			expected: "dsn",
		},
		{
			name:     "array index",
			input:    "/queries/0",
			expected: "queries[0]",
		},
		{
			name:     "nested array property",
			input:    "/queries/0/path",
			expected: "queries[0].path",
		},
		{
			name:     "deeply nested",
			input:    "/queries/0/mock/array/1",
			expected: "queries[0].mock.array[1]",
		},
		{
			name:     "mutations",
			input:    "/mutations/2/sql",
			expected: "mutations[2].sql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"0", true},
		{"1", true},
		{"123", true},
		{"", false},
		{"a", false},
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

func TestIsQueryOrMutationPath(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"/queries/0", true},
		{"/queries/10", true},
		{"/mutations/0", true},
		{"/mutations/5", true},
		{"/queries/0/sql", false},
		{"/queries", false},
		{"/dsn", false},
		{"/queries/abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isQueryOrMutationPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
