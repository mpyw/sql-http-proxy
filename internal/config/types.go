package config

import (
	"errors"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// OpType represents the operation type for result handling.
type OpType string

const (
	// OpTypeOne returns a single row.
	OpTypeOne OpType = "one"
	// OpTypeMany returns multiple rows.
	OpTypeMany OpType = "many"
	// OpTypeNone returns no content (mutations only).
	OpTypeNone OpType = "none"
)

// EntityType represents the entity type (query or mutation).
type EntityType string

const (
	// EntityQuery represents a query operation.
	EntityQuery EntityType = "query"
	// EntityMutation represents a mutation operation.
	EntityMutation EntityType = "mutation"
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
type AcceptType string

const (
	AcceptJSON AcceptType = "json"
	AcceptForm AcceptType = "form"
)

// AcceptTypes represents a list of accepted content types.
// Can be unmarshaled from either a string or array of strings.
type AcceptTypes []AcceptType

// DefaultAcceptTypes is the default list when accepts is not specified.
var DefaultAcceptTypes = AcceptTypes{AcceptJSON, AcceptForm}

// UnmarshalYAML implements custom YAML unmarshaling for AcceptTypes.
// Accepts either a string or an array of strings.
func (a *AcceptTypes) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		var s AcceptType
		if err := value.Decode(&s); err != nil {
			return err
		}
		*a = AcceptTypes{s}
		return nil
	}

	var arr []AcceptType
	if err := value.Decode(&arr); err != nil {
		return err
	}
	*a = arr
	return nil
}

// Query represents a single query configuration.
type Query struct {
	Type           QueryType    `yaml:"type"`
	Method         string       `yaml:"method,omitempty"` // HTTP method (default: GET)
	Path           string       `yaml:"path"`
	SQL            string       `yaml:"sql,omitempty"`
	Mock           *Mock        `yaml:"mock,omitempty"`             // Mock data source (alternative to SQL)
	Accepts        *AcceptTypes `yaml:"accepts,omitempty"`          // Accepted content types for body (nil: default, empty: none)
	HandleNotFound bool         `yaml:"handle_not_found,omitempty"` // Pass null to post-transform instead of 404 (type: one only)
	Transform      *Transform   `yaml:"transform,omitempty"`
	Delay          string       `yaml:"delay,omitempty"` // Artificial delay before response (e.g., "100ms", "1s")
}

// GetMethod returns the HTTP method for the query (default: GET).
func (q Query) GetMethod() string {
	if q.Method == "" {
		return "GET"
	}
	return q.Method
}

// GetAccepts returns the accepted content types.
// Returns default [json, form] if not specified, or empty slice if explicitly set to [].
func (q Query) GetAccepts() AcceptTypes {
	if q.Accepts == nil {
		return DefaultAcceptTypes
	}
	return *q.Accepts
}

// HasMock returns true if this query uses mock data.
func (q Query) HasMock() bool {
	return q.Mock != nil && !q.Mock.IsEmpty()
}

// Mutation represents a single mutation configuration.
type Mutation struct {
	Type      MutationType `yaml:"type"`
	Method    string       `yaml:"method,omitempty"` // HTTP method (default: POST)
	Path      string       `yaml:"path"`
	SQL       string       `yaml:"sql,omitempty"`
	Mock      *Mock        `yaml:"mock,omitempty"`    // Mock data source (alternative to SQL)
	Accepts   *AcceptTypes `yaml:"accepts,omitempty"` // Accepted content types for body (nil: default, empty: none)
	Transform *Transform   `yaml:"transform,omitempty"`
	Delay     string       `yaml:"delay,omitempty"` // Artificial delay before response (e.g., "100ms", "1s")
}

// GetMethod returns the HTTP method for the mutation (default: POST).
func (m Mutation) GetMethod() string {
	if m.Method == "" {
		return "POST"
	}
	return m.Method
}

// GetAccepts returns the accepted content types.
// Returns default [json, form] if not specified, or empty slice if explicitly set to [].
func (m Mutation) GetAccepts() AcceptTypes {
	if m.Accepts == nil {
		return DefaultAcceptTypes
	}
	return *m.Accepts
}

// HasMock returns true if this mutation uses mock data.
func (m Mutation) HasMock() bool {
	return m.Mock != nil && !m.Mock.IsEmpty()
}

// Transform defines pre/post JavaScript transformations.
// Note: mock is now at query/mutation level, not inside transform.
type Transform struct {
	Pre  string        `yaml:"pre,omitempty"`
	Post PostTransform `yaml:"post,omitempty"`
}

// Mock represents mock data source configuration.
// Contains all possible data source fields for both type: one and type: many.
// Only one data source field should be set (mutually exclusive).
// For type: none, mock can be just `true` (marker for mock mode).
type Mock struct {
	// Enabled is set to true when mock is specified as `mock: true` (for type: none)
	Enabled bool `yaml:"-"`

	// Object sources (for type: one - returns single object directly)
	Object         map[string]any `yaml:"object,omitempty"`
	ObjectJSON     string         `yaml:"object_json,omitempty"`
	ObjectJSONFile string         `yaml:"object_json_file,omitempty"`
	ObjectJS       string         `yaml:"object_js,omitempty"`

	// Array sources (for type: many, or type: one with filter)
	Array         []any  `yaml:"array,omitempty"`
	ArrayJSON     string `yaml:"array_json,omitempty"`
	ArrayJSONFile string `yaml:"array_json_file,omitempty"`
	ArrayJS       string `yaml:"array_js,omitempty"`
	CSV           string `yaml:"csv,omitempty"`
	CSVFile       string `yaml:"csv_file,omitempty"`
	JSONL         string `yaml:"jsonl,omitempty"`
	JSONLFile     string `yaml:"jsonl_file,omitempty"`

	// Filter for array sources (JS code: receives row, input params, ctx free var; returns boolean)
	// Required for type: one with array sources, optional for type: many
	Filter string `yaml:"filter,omitempty"`
}

// UnmarshalYAML implements custom YAML unmarshaling for Mock.
// Accepts either a boolean (for type: none) or an object (for type: one/many).
func (m *Mock) UnmarshalYAML(value *yaml.Node) error {
	// Try as boolean first (for type: none)
	if value.Kind == yaml.ScalarNode && value.Tag == "!!bool" {
		var b bool
		if err := value.Decode(&b); err != nil {
			return err
		}
		m.Enabled = b
		return nil
	}

	// Try as object
	type mockRaw Mock
	var raw mockRaw
	if err := value.Decode(&raw); err != nil {
		return err
	}
	*m = Mock(raw)
	return nil
}

// IsEmpty returns true if no mock source is configured.
func (m *Mock) IsEmpty() bool {
	if m.Enabled {
		return false
	}
	return !lo.Contains(m.sourceFlags(), true)
}

// IsArraySource returns true if this mock uses an array source.
func (m *Mock) IsArraySource() bool {
	return lo.Contains(m.arraySourceFlags(), true)
}

// IsObjectSource returns true if this mock uses an object source.
func (m *Mock) IsObjectSource() bool {
	return lo.Contains(m.objectSourceFlags(), true)
}

func (m *Mock) objectSourceFlags() []bool {
	return []bool{m.Object != nil, m.ObjectJSON != "", m.ObjectJSONFile != "", m.ObjectJS != ""}
}

func (m *Mock) arraySourceFlags() []bool {
	return []bool{
		m.Array != nil, m.ArrayJSON != "", m.ArrayJSONFile != "", m.ArrayJS != "",
		m.CSV != "", m.CSVFile != "",
		m.JSONL != "", m.JSONLFile != "",
	}
}

func (m *Mock) sourceFlags() []bool {
	return append(m.objectSourceFlags(), m.arraySourceFlags()...)
}

// GetJS returns the JavaScript code if present.
func (m *Mock) GetJS() string {
	if m.ObjectJS != "" {
		return m.ObjectJS
	}
	return m.ArrayJS
}

// Validate checks that only one mock source is specified.
// Does NOT check filter requirement - that depends on query type and is validated elsewhere.
func (m *Mock) Validate() error {
	if lo.Count(m.sourceFlags(), true) > 1 {
		return errors.New("mock: only one source type can be specified")
	}
	return nil
}

// ValidateForTypeOne validates mock configuration for type: one.
// Array sources require filter.
func (m *Mock) ValidateForTypeOne() error {
	if err := m.Validate(); err != nil {
		return err
	}

	// Array sources require filter for type: one
	if m.IsArraySource() && m.Filter == "" {
		return errors.New("mock: filter is required for array sources with type: one")
	}

	return nil
}

// ValidateForTypeMany validates mock configuration for type: many.
// Object sources are not allowed.
func (m *Mock) ValidateForTypeMany() error {
	if err := m.Validate(); err != nil {
		return err
	}

	// Object sources are not allowed for type: many
	if m.IsObjectSource() {
		return errors.New("mock: object sources are not allowed for type: many (use array sources)")
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
	ValueParser string `yaml:"value_parser,omitempty"`
}

// DatabaseInit defines SQL to execute on database connection.
// Can be a simple string (sql code) or an object with sql and/or sql_files.
type DatabaseInit struct {
	SQL      string   `yaml:"sql,omitempty"`
	SQLFiles []string `yaml:"sql_files,omitempty"`
}

// UnmarshalYAML implements custom YAML unmarshaling for DatabaseInit.
// Accepts either a string (shorthand for sql:) or an object.
func (d *DatabaseInit) UnmarshalYAML(value *yaml.Node) error {
	// Try as string first (shorthand: string → sql)
	if value.Kind == yaml.ScalarNode {
		var s string
		if err := value.Decode(&s); err != nil {
			return err
		}
		d.SQL = s
		return nil
	}

	// Try as object
	type databaseInitRaw DatabaseInit
	var raw databaseInitRaw
	if err := value.Decode(&raw); err != nil {
		return err
	}
	*d = DatabaseInit(raw)
	return nil
}
