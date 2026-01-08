// Package config provides configuration parsing for sql-http-proxy.
package config

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/a8m/envsubst"
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

// Config represents the application configuration.
type Config struct {
	GlobalHelpers *GlobalHelpers `yaml:"global_helpers,omitempty"`
	CSV           *CSVConfig     `yaml:"csv,omitempty"`
	DSN           string         `yaml:"dsn"`
	Queries       []Query        `yaml:"queries"`
	Mutations     []Mutation     `yaml:"mutations"`
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

// RequiresDB returns true if any query or mutation requires a database connection.
func (cfg *Config) RequiresDB() bool {
	return lo.SomeBy(cfg.Queries, func(q Query) bool {
		return !q.HasMock()
	}) || lo.SomeBy(cfg.Mutations, func(m Mutation) bool {
		return !m.HasMock()
	})
}

// validate checks semantic constraints that JSON Schema cannot express.
// Note: type-specific constraints (each for many, handle_not_found for one)
// are now enforced by JSON Schema via separate queryOne/queryMany definitions.
func (cfg *Config) validate() error {
	var errs []error

	// Validate global_helpers JS compilation
	if cfg.GlobalHelpers != nil && cfg.GlobalHelpers.JS != "" {
		if _, err := js.CompilePre(cfg.GlobalHelpers.JS); err != nil {
			errs = append(errs, fmt.Errorf("global_helpers.js: %w", err))
		}
	}

	// Validate csv.value_parser
	if cfg.CSV != nil && cfg.CSV.ValueParser != "" {
		if _, err := js.CompilePre(cfg.CSV.ValueParser); err != nil {
			errs = append(errs, fmt.Errorf("csv.value_parser: %w", err))
		}
	}

	return errors.Join(errs...)
}

// ValidateTransforms validates all transforms and mocks are valid.
// Returns all errors joined, not just the first one.
func (cfg *Config) ValidateTransforms() error {
	var errs []error

	validateTransform := func(t *Transform, location string) {
		if t == nil {
			return
		}
		// Validate pre-transform JS
		if t.Pre != "" {
			if _, err := js.CompilePre(t.Pre); err != nil {
				errs = append(errs, fmt.Errorf("%s pre: %w", location, err))
			}
		}
		// Validate post.each JS
		if t.Post.Each != "" {
			if _, err := js.CompilePost(t.Post.Each); err != nil {
				errs = append(errs, fmt.Errorf("%s post.each: %w", location, err))
			}
		}
		// Validate post.all JS
		if t.Post.All != "" {
			if _, err := js.CompilePost(t.Post.All); err != nil {
				errs = append(errs, fmt.Errorf("%s post.all: %w", location, err))
			}
		}
	}

	validateMock := func(m *Mock, isTypeOne bool, location string) {
		if m == nil || m.IsEmpty() {
			return
		}
		// Validate type-specific constraints
		if isTypeOne {
			if err := m.ValidateForTypeOne(); err != nil {
				errs = append(errs, fmt.Errorf("%s %w", location, err))
			}
		} else {
			if err := m.ValidateForTypeMany(); err != nil {
				errs = append(errs, fmt.Errorf("%s %w", location, err))
			}
		}
		// Compile JS mock if specified
		if m.HasJS() {
			if _, err := js.CompilePre(m.GetJS()); err != nil {
				errs = append(errs, fmt.Errorf("%s mock js: %w", location, err))
			}
		}
		// Compile filter JS if specified
		if m.Filter != "" {
			if _, err := js.CompileFilter(m.Filter); err != nil {
				errs = append(errs, fmt.Errorf("%s mock.filter: %w", location, err))
			}
		}
	}

	for i, q := range cfg.Queries {
		location := fmt.Sprintf("queries[%d] (%s)", i, q.Path)
		validateTransform(q.Transform, location)
		validateMock(q.Mock, q.Type == QueryTypeOne, location)
	}
	for i, m := range cfg.Mutations {
		location := fmt.Sprintf("mutations[%d] (%s)", i, m.Path)
		validateTransform(m.Transform, location)
		isTypeOne := m.Type == MutationTypeOne || m.Type == MutationTypeNone
		validateMock(m.Mock, isTypeOne, location)
	}
	return errors.Join(errs...)
}

// Parse parses configuration from YAML bytes.
func Parse(data []byte) (Config, error) {
	// Parse YAML to generic interface for schema validation
	var raw any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return Config{}, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate against JSON Schema
	if err := configSchema.Validate(raw); err != nil {
		if validationErr, ok := err.(*jsonschema.ValidationError); ok {
			return Config{}, fmt.Errorf("config validation failed: %s", formatValidationError(validationErr))
		}
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

	// Expand environment variables in DSN (e.g., ${DB_PASSWORD}, $DB_HOST, ${DB_PORT:-5432})
	if cfg.DSN != "" {
		expanded, err := envsubst.String(cfg.DSN)
		if err != nil {
			return Config{}, fmt.Errorf("failed to expand DSN: %w", err)
		}
		cfg.DSN = expanded
	}

	// DSN is required only if any query needs a database connection
	if cfg.DSN == "" && cfg.RequiresDB() {
		return Config{}, errors.New("missing dsn: set dsn in YAML config (supports ${VAR} and ${VAR:-default} env expansion)")
	}
	return cfg, nil
}

// ParseFile parses configuration from a YAML file.
func ParseFile(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open file: %w", err)
	}
	return Parse(data)
}
