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
// Accepts either a string (shorthand for js:) or an object with source type.
func (m *MockTransform) UnmarshalYAML(value *yaml.Node) error {
	// Try as string first (backward compatibility: mock: "code" → mock.js: "code")
	if value.Kind == yaml.ScalarNode {
		var s string
		if err := value.Decode(&s); err != nil {
			return err
		}
		m.JS = s
		return nil
	}

	// Try as object
	type mockTransformRaw MockTransform
	var raw mockTransformRaw
	if err := value.Decode(&raw); err != nil {
		return err
	}
	*m = MockTransform(raw)
	return nil
}

// IsEmpty returns true if no mock source is configured.
func (m *MockTransform) IsEmpty() bool {
	return m.JS == "" && m.JSFile == "" &&
		m.CSV == "" && m.CSVFile == "" &&
		m.JSON == nil && m.JSONFile == "" &&
		m.JSONL == "" && m.JSONLFile == ""
}

// Type returns the mock source type ("js", "csv", "json", "jsonl", or "" if empty).
func (m *MockTransform) Type() string {
	switch {
	case m.JS != "" || m.JSFile != "":
		return "js"
	case m.CSV != "" || m.CSVFile != "":
		return "csv"
	case m.JSON != nil || m.JSONFile != "":
		return "json"
	case m.JSONL != "" || m.JSONLFile != "":
		return "jsonl"
	default:
		return ""
	}
}

// IsFile returns true if this mock uses an external file.
func (m *MockTransform) IsFile() bool {
	return m.JSFile != "" || m.CSVFile != "" || m.JSONFile != "" || m.JSONLFile != ""
}

// FilePath returns the file path if IsFile() is true, otherwise empty string.
func (m *MockTransform) FilePath() string {
	switch {
	case m.JSFile != "":
		return m.JSFile
	case m.CSVFile != "":
		return m.CSVFile
	case m.JSONFile != "":
		return m.JSONFile
	case m.JSONLFile != "":
		return m.JSONLFile
	default:
		return ""
	}
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

// IsEmpty returns true if no post-transform is configured.
func (p PostTransform) IsEmpty() bool {
	return p.Each == "" && p.All == ""
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
