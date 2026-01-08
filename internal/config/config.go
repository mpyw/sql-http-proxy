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

	"github.com/mpyw/sql-http-proxy/internal"
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

// ValidateTransforms validates all transforms and mocks are valid.
// Returns all errors joined, not just the first one.
func (cfg *Config) ValidateTransforms() error {
	var err error

	// Validate global helpers and CSV value parser
	for _, entry := range internal.SliceOf(
		lo.T2("global_helpers.js", lo.FromPtr(cfg.GlobalHelpers).JS),
		lo.T2("csv.value_parser", lo.FromPtr(cfg.CSV).ValueParser),
	) {
		name, script := entry.Unpack()
		if script != "" {
			if _, e := js.CompilePre(script); e != nil {
				err = errors.Join(err, fmt.Errorf("%s: %w", name, e))
			}
		}
	}

	validateJS := func(
		entries []lo.Tuple3[
			string,                                // name
			func(string) (*js.Transformer, error), // compileFunc
			string,                                // script
		],
		location string,
	) {
		for _, entry := range entries {
			name, compileFunc, script := entry.Unpack()
			if script != "" {
				if _, e := compileFunc(script); e != nil {
					err = errors.Join(err, fmt.Errorf("%s %s: %w", location, name, e))
				}
			}
		}
	}

	validateTransform := func(t *Transform, location string) {
		if t == nil {
			return
		}
		validateJS(internal.SliceOf(
			lo.T3("pre", js.CompilePre, t.Pre),
			lo.T3("post.each", js.CompilePost, t.Post.Each),
			lo.T3("post.all", js.CompilePost, t.Post.All),
		), location)
	}

	validateMock := func(m *Mock, isTypeOne bool, location string) {
		if m == nil || m.IsEmpty() {
			return
		}
		// Validate type-specific constraints
		if e := lo.Ternary(isTypeOne, m.ValidateForTypeOne, m.ValidateForTypeMany)(); e != nil {
			err = errors.Join(err, fmt.Errorf("%s %w", location, e))
		}
		validateJS(internal.SliceOf(
			lo.T3("mock js", js.CompilePre, m.GetJS()),
			lo.T3("mock.filter", js.CompileFilter, m.Filter),
		), location)
	}

	for i, q := range cfg.Queries {
		location := fmt.Sprintf("queries[%d] (%s)", i, q.Path)
		validateTransform(q.Transform, location)
		validateMock(
			q.Mock,
			q.Type == QueryTypeOne,
			location,
		)
	}
	for i, m := range cfg.Mutations {
		location := fmt.Sprintf("mutations[%d] (%s)", i, m.Path)
		validateTransform(m.Transform, location)
		validateMock(
			m.Mock,
			m.Type == MutationTypeOne || m.Type == MutationTypeNone,
			location,
		)
	}

	return err
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
		if validationErr := (*jsonschema.ValidationError)(nil); errors.As(err, &validationErr) {
			return Config{}, fmt.Errorf("config validation failed: %s", formatValidationError(validationErr))
		}
		return Config{}, fmt.Errorf("config validation failed: %w", err)
	}

	// Unmarshal into struct
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
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
