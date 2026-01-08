// Package config provides configuration parsing for sql-http-proxy.
package config

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/samber/lo"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

//go:embed schema.json
var schemaJSON []byte

var configSchema = func() *jsonschema.Schema {
	schemaDoc := lo.Must(jsonschema.UnmarshalJSON(bytes.NewReader(schemaJSON)))
	compiler := jsonschema.NewCompiler()
	lo.Must0(compiler.AddResource("schema.json", schemaDoc))
	return lo.Must(compiler.Compile("schema.json"))
}()

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

// Config represents the application configuration.
type Config struct {
	DSN       string     `yaml:"dsn"`
	Queries   []Query    `yaml:"queries"`
	Mutations []Mutation `yaml:"mutations"`
}

// Query represents a single query configuration.
type Query struct {
	Type           QueryType  `yaml:"type"`
	Path           string     `yaml:"path"`
	SQL            string     `yaml:"sql"`
	HandleNotFound bool       `yaml:"handle_not_found,omitempty"` // Pass null to post-transform instead of 404 (type: one only)
	Transform      *Transform `yaml:"transform,omitempty"`
}

// Mutation represents a single mutation configuration.
type Mutation struct {
	Type      MutationType `yaml:"type"`
	Method    string       `yaml:"method,omitempty"` // HTTP method (default: POST)
	Path      string       `yaml:"path"`
	SQL       string       `yaml:"sql"`
	Transform *Transform   `yaml:"transform,omitempty"`
}

// GetMethod returns the HTTP method for the mutation (default: POST).
func (m Mutation) GetMethod() string {
	if m.Method == "" {
		return "POST"
	}
	return m.Method
}

// Transform defines pre/post JavaScript transformations.
type Transform struct {
	Pre  string        `yaml:"pre,omitempty"`
	Mock string        `yaml:"mock,omitempty"` // Mock DB response (replaces actual query)
	Post PostTransform `yaml:"post,omitempty"`
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

// Driver returns the database driver name from the DSN.
func (cfg *Config) Driver() (string, error) {
	u, err := url.Parse(cfg.DSN)
	if err != nil {
		return "", fmt.Errorf("failed to parse DSN: %w", err)
	}
	switch u.Scheme {
	case "postgres", "postgresql":
		return "pgx", nil
	case "file", "sqlite":
		return "sqlite", nil
	default:
		return u.Scheme, nil
	}
}

// Getenv is a variable for testing purposes.
var Getenv = os.Getenv

// Parse parses configuration from YAML bytes.
func Parse(data []byte) (Config, error) {
	// Parse YAML to generic interface for schema validation
	var raw any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return Config{}, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate against JSON Schema
	if err := configSchema.Validate(raw); err != nil {
		return Config{}, fmt.Errorf("config validation failed: %w", err)
	}

	// Unmarshal into struct
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate semantic constraints
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	if cfg.DSN == "" {
		cfg.DSN = Getenv("SQL_PROXY_DSN")
	}

	// DSN is required only if any query needs a database connection
	if cfg.DSN == "" && cfg.RequiresDB() {
		return Config{}, errors.New("missing DSN: YAML config \"dsn\" or environment variable \"SQL_PROXY_DSN\" must be passed")
	}
	return cfg, nil
}

// validate checks semantic constraints that JSON Schema cannot express.
// Note: type-specific constraints (each for many, handle_not_found for one)
// are now enforced by JSON Schema via separate queryOne/queryMany definitions.
func (cfg *Config) validate() error {
	// Reserved for future semantic validations
	return nil
}

// RequiresDB returns true if any query or mutation requires a database connection.
func (cfg *Config) RequiresDB() bool {
	for _, q := range cfg.Queries {
		if q.Transform == nil || q.Transform.Mock == "" {
			return true
		}
	}
	for _, m := range cfg.Mutations {
		if m.Transform == nil || m.Transform.Mock == "" {
			return true
		}
	}
	return false
}

// ParseFile parses configuration from a YAML file.
func ParseFile(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open file: %w", err)
	}
	return Parse(data)
}

// ValidateTransforms validates all JavaScript transforms compile correctly.
func (cfg *Config) ValidateTransforms() error {
	for i, q := range cfg.Queries {
		if q.Transform == nil {
			continue
		}
		if q.Transform.Pre != "" {
			if _, err := js.CompilePre(q.Transform.Pre); err != nil {
				return fmt.Errorf("queries[%d] (%s): pre-transform: %w", i, q.Path, err)
			}
		}
		if q.Transform.Mock != "" {
			if _, err := js.CompilePre(q.Transform.Mock); err != nil {
				return fmt.Errorf("queries[%d] (%s): mock: %w", i, q.Path, err)
			}
		}
		if q.Transform.Post.Each != "" {
			if _, err := js.Compile(q.Transform.Post.Each, []string{"ctx", "input", "output"}); err != nil {
				return fmt.Errorf("queries[%d] (%s): post.each: %w", i, q.Path, err)
			}
		}
		if q.Transform.Post.All != "" {
			if _, err := js.Compile(q.Transform.Post.All, []string{"ctx", "input", "output"}); err != nil {
				return fmt.Errorf("queries[%d] (%s): post.all: %w", i, q.Path, err)
			}
		}
	}

	for i, m := range cfg.Mutations {
		if m.Transform == nil {
			continue
		}
		if m.Transform.Pre != "" {
			if _, err := js.CompilePre(m.Transform.Pre); err != nil {
				return fmt.Errorf("mutations[%d] (%s): pre-transform: %w", i, m.Path, err)
			}
		}
		if m.Transform.Mock != "" {
			if _, err := js.CompilePre(m.Transform.Mock); err != nil {
				return fmt.Errorf("mutations[%d] (%s): mock: %w", i, m.Path, err)
			}
		}
		if m.Transform.Post.Each != "" {
			if _, err := js.Compile(m.Transform.Post.Each, []string{"ctx", "input", "output"}); err != nil {
				return fmt.Errorf("mutations[%d] (%s): post.each: %w", i, m.Path, err)
			}
		}
		if m.Transform.Post.All != "" {
			if _, err := js.Compile(m.Transform.Post.All, []string{"ctx", "input", "output"}); err != nil {
				return fmt.Errorf("mutations[%d] (%s): post.all: %w", i, m.Path, err)
			}
		}
	}
	return nil
}
