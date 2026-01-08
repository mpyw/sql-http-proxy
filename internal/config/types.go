package config

import (
	"errors"

	"gopkg.in/yaml.v3"
)

// QueryType represents the type of query result.
type QueryType string

const (
	// QueryTypeOne returns a single row.
	QueryTypeOne QueryType = "one"
	// QueryTypeMany returns multiple rows.
	QueryTypeMany QueryType = "many"
)

// MutationType represents the type of mutation result.
type MutationType string

const (
	// MutationTypeOne returns a single row (via RETURNING).
	MutationTypeOne MutationType = "one"
	// MutationTypeMany returns multiple rows (via RETURNING).
	MutationTypeMany MutationType = "many"
	// MutationTypeNone returns no content (204).
	MutationTypeNone MutationType = "none"
)

// AcceptType represents accepted content types.
const (
	AcceptJSON = "json"
	AcceptForm = "form"
)

// AcceptTypes represents a list of accepted content types.
// Can be unmarshaled from either a string or array of strings.
type AcceptTypes []string

// DefaultAcceptTypes is the default list when accepts is not specified.
var DefaultAcceptTypes = AcceptTypes{AcceptJSON, AcceptForm}

// UnmarshalYAML implements custom YAML unmarshaling for AcceptTypes.
// Accepts either a string or an array of strings.
func (a *AcceptTypes) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		var s string
		if err := value.Decode(&s); err != nil {
			return err
		}
		*a = AcceptTypes{s}
		return nil
	}

	var arr []string
	if err := value.Decode(&arr); err != nil {
		return err
	}
	*a = arr
	return nil
}

// Query represents a single query configuration.
type Query struct {
	Type           QueryType   `yaml:"type"`
	Method         string      `yaml:"method,omitempty"` // HTTP method (default: GET)
	Path           string      `yaml:"path"`
	SQL            string      `yaml:"sql"`
	Accepts        AcceptTypes `yaml:"accepts,omitempty"`          // Accepted content types for body (default: [json, form])
	HandleNotFound bool        `yaml:"handle_not_found,omitempty"` // Pass null to post-transform instead of 404 (type: one only)
	Transform      *Transform  `yaml:"transform,omitempty"`
}

// GetMethod returns the HTTP method for the query (default: GET).
func (q Query) GetMethod() string {
	if q.Method == "" {
		return "GET"
	}
	return q.Method
}

// GetAccepts returns the accepted content types (default: [json, form]).
func (q Query) GetAccepts() AcceptTypes {
	if len(q.Accepts) == 0 {
		return DefaultAcceptTypes
	}
	return q.Accepts
}

// Mutation represents a single mutation configuration.
type Mutation struct {
	Type      MutationType `yaml:"type"`
	Method    string       `yaml:"method,omitempty"` // HTTP method (default: POST)
	Path      string       `yaml:"path"`
	SQL       string       `yaml:"sql"`
	Accepts   AcceptTypes  `yaml:"accepts,omitempty"` // Accepted content types for body (default: [json, form])
	Transform *Transform   `yaml:"transform,omitempty"`
}

// GetMethod returns the HTTP method for the mutation (default: POST).
func (m Mutation) GetMethod() string {
	if m.Method == "" {
		return "POST"
	}
	return m.Method
}

// GetAccepts returns the accepted content types (default: [json, form]).
func (m Mutation) GetAccepts() AcceptTypes {
	if len(m.Accepts) == 0 {
		return DefaultAcceptTypes
	}
	return m.Accepts
}

// Transform defines pre/post JavaScript transformations.
type Transform struct {
	Pre  string         `yaml:"pre,omitempty"`
	Mock *MockTransform `yaml:"mock,omitempty"` // Mock DB response (replaces actual query)
	Post PostTransform  `yaml:"post,omitempty"`
}

// MockTransform represents mock data source configuration.
// Only one field should be set (mutually exclusive).
type MockTransform struct {
	JS        string `yaml:"js,omitempty"`
	JSFile    string `yaml:"js_file,omitempty"`
	CSV       string `yaml:"csv,omitempty"`
	CSVFile   string `yaml:"csv_file,omitempty"`
	JSON      any    `yaml:"json,omitempty"`
	JSONFile  string `yaml:"json_file,omitempty"`
	JSONL     string `yaml:"jsonl,omitempty"`
	JSONLFile string `yaml:"jsonl_file,omitempty"`
}

// UnmarshalYAML implements custom YAML unmarshaling for MockTransform.
// Accepts:
// - string → mock.js (JS code)
// - array → mock.json (JSON array shorthand)
// - object with known key (js, json, csv, etc.) → use as-is
// - object without known keys → mock.json (JSON object shorthand)
func (m *MockTransform) UnmarshalYAML(value *yaml.Node) error {
	// String → JS code
	if value.Kind == yaml.ScalarNode {
		var s string
		if err := value.Decode(&s); err != nil {
			return err
		}
		m.JS = s
		return nil
	}

	// Array → JSON array shorthand
	if value.Kind == yaml.SequenceNode {
		var arr []any
		if err := value.Decode(&arr); err != nil {
			return err
		}
		m.JSON = arr
		return nil
	}

	// Object: check if it has known keys
	if value.Kind == yaml.MappingNode {
		// First, try to decode as MockTransform with known keys
		type mockTransformRaw MockTransform
		var raw mockTransformRaw
		if err := value.Decode(&raw); err != nil {
			return err
		}

		// If any known field is set, use the decoded result
		if raw.JS != "" || raw.JSFile != "" ||
			raw.CSV != "" || raw.CSVFile != "" ||
			raw.JSON != nil || raw.JSONFile != "" ||
			raw.JSONL != "" || raw.JSONLFile != "" {
			*m = MockTransform(raw)
			return nil
		}

		// No known keys found → treat as JSON object shorthand
		var obj map[string]any
		if err := value.Decode(&obj); err != nil {
			return err
		}
		m.JSON = obj
		return nil
	}

	return nil
}

// IsEmpty returns true if no mock source is configured.
func (m *MockTransform) IsEmpty() bool {
	return m.JS == "" && m.JSFile == "" &&
		m.CSV == "" && m.CSVFile == "" &&
		m.JSON == nil && m.JSONFile == "" &&
		m.JSONL == "" && m.JSONLFile == ""
}

// Validate checks that only one mock source is specified.
func (m *MockTransform) Validate() error {
	count := 0
	if m.JS != "" {
		count++
	}
	if m.JSFile != "" {
		count++
	}
	if m.CSV != "" {
		count++
	}
	if m.CSVFile != "" {
		count++
	}
	if m.JSON != nil {
		count++
	}
	if m.JSONFile != "" {
		count++
	}
	if m.JSONL != "" {
		count++
	}
	if m.JSONLFile != "" {
		count++
	}
	if count > 1 {
		return errors.New("mock: only one source type can be specified (js, js_file, csv, csv_file, json, json_file, jsonl, jsonl_file)")
	}
	return nil
}

// PostTransform represents post-transform configuration.
// Can be a simple string or an object with each/all.
type PostTransform struct {
	Each string `yaml:"each,omitempty"` // Transform each row individually
	All  string `yaml:"all,omitempty"`  // Transform entire result array (default for many)
}

// UnmarshalYAML implements custom YAML unmarshaling for PostTransform.
// Accepts either a string (defaults to All) or an object with each/all.
func (p *PostTransform) UnmarshalYAML(value *yaml.Node) error {
	// Try as string first
	if value.Kind == yaml.ScalarNode {
		var s string
		if err := value.Decode(&s); err != nil {
			return err
		}
		p.All = s
		return nil
	}

	// Try as object
	type postTransformRaw PostTransform
	var raw postTransformRaw
	if err := value.Decode(&raw); err != nil {
		return err
	}
	*p = PostTransform(raw)
	return nil
}

// GlobalHelpers defines global JavaScript helpers available in all JS contexts.
// Can be a simple string (js code) or an object with js and/or js_files.
type GlobalHelpers struct {
	JS      string   `yaml:"js,omitempty"`
	JSFiles []string `yaml:"js_files,omitempty"`
}

// UnmarshalYAML implements custom YAML unmarshaling for GlobalHelpers.
// Accepts either a string (shorthand for js:) or an object.
func (g *GlobalHelpers) UnmarshalYAML(value *yaml.Node) error {
	// Try as string first (shorthand: string → js)
	if value.Kind == yaml.ScalarNode {
		var s string
		if err := value.Decode(&s); err != nil {
			return err
		}
		g.JS = s
		return nil
	}

	// Try as object
	type globalHelpersRaw GlobalHelpers
	var raw globalHelpersRaw
	if err := value.Decode(&raw); err != nil {
		return err
	}
	*g = GlobalHelpers(raw)
	return nil
}

// CSVConfig defines global CSV parsing options.
type CSVConfig struct {
	ValueParser *JSCode `yaml:"value_parser,omitempty"`
}

// JSCode represents JavaScript code from inline or file.
// Only one of JS or JSFile should be set.
// Can be unmarshaled from a string (shorthand for js:) or an object.
type JSCode struct {
	JS     string `yaml:"js,omitempty"`
	JSFile string `yaml:"js_file,omitempty"`
}

// UnmarshalYAML implements custom YAML unmarshaling for JSCode.
// Accepts either a string (shorthand for js:) or an object.
func (j *JSCode) UnmarshalYAML(value *yaml.Node) error {
	// Try as string first (shorthand: string → js)
	if value.Kind == yaml.ScalarNode {
		var s string
		if err := value.Decode(&s); err != nil {
			return err
		}
		j.JS = s
		return nil
	}

	// Try as object
	type jsCodeRaw JSCode
	var raw jsCodeRaw
	if err := value.Decode(&raw); err != nil {
		return err
	}
	*j = JSCode(raw)
	return nil
}

// Validate checks that only one of js or js_file is specified.
func (j *JSCode) Validate() error {
	if j.JS != "" && j.JSFile != "" {
		return errors.New("only one of 'js' or 'js_file' can be specified")
	}
	if j.JS == "" && j.JSFile == "" {
		return errors.New("one of 'js' or 'js_file' must be specified")
	}
	return nil
}
