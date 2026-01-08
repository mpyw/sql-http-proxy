package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Driver(t *testing.T) {
	t.Run("postgres", func(t *testing.T) {
		cfg := &Config{DSN: "postgres://user:pass@localhost:5432/db"}
		driver, err := cfg.Driver()
		require.NoError(t, err)
		assert.Equal(t, "pgx", driver)
	})

	t.Run("postgresql", func(t *testing.T) {
		cfg := &Config{DSN: "postgresql://user:pass@localhost:5432/db"}
		driver, err := cfg.Driver()
		require.NoError(t, err)
		assert.Equal(t, "pgx", driver)
	})

	t.Run("mysql", func(t *testing.T) {
		cfg := &Config{DSN: "mysql://user:pass@localhost:3306/db"}
		driver, err := cfg.Driver()
		require.NoError(t, err)
		assert.Equal(t, "mysql", driver)
	})

	t.Run("file (sqlite)", func(t *testing.T) {
		cfg := &Config{DSN: "file:./test.db"}
		driver, err := cfg.Driver()
		require.NoError(t, err)
		assert.Equal(t, "sqlite", driver)
	})

	t.Run("sqlite", func(t *testing.T) {
		cfg := &Config{DSN: "sqlite:./test.db"}
		driver, err := cfg.Driver()
		require.NoError(t, err)
		assert.Equal(t, "sqlite", driver)
	})

	t.Run("sqlserver", func(t *testing.T) {
		cfg := &Config{DSN: "sqlserver://user:pass@localhost:1433"}
		driver, err := cfg.Driver()
		require.NoError(t, err)
		assert.Equal(t, "sqlserver", driver)
	})

	t.Run("invalid DSN", func(t *testing.T) {
		cfg := &Config{DSN: "://invalid"}
		_, err := cfg.Driver()
		require.Error(t, err)
	})
}

func TestConfig_RequiresDB(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		cfg := &Config{}
		assert.False(t, cfg.RequiresDB())
	})

	t.Run("all mock queries", func(t *testing.T) {
		cfg := &Config{
			Queries: []Query{
				{Mock: &Mock{Object: map[string]any{"id": 1}}},
			},
		}
		assert.False(t, cfg.RequiresDB())
	})

	t.Run("query with SQL", func(t *testing.T) {
		cfg := &Config{
			Queries: []Query{
				{SQL: "SELECT * FROM users"},
			},
		}
		assert.True(t, cfg.RequiresDB())
	})

	t.Run("mixed queries", func(t *testing.T) {
		cfg := &Config{
			Queries: []Query{
				{Mock: &Mock{Object: map[string]any{"id": 1}}},
				{SQL: "SELECT * FROM users"},
			},
		}
		assert.True(t, cfg.RequiresDB())
	})

	t.Run("all mock mutations", func(t *testing.T) {
		cfg := &Config{
			Mutations: []Mutation{
				{Mock: &Mock{Object: map[string]any{"id": 1}}},
			},
		}
		assert.False(t, cfg.RequiresDB())
	})

	t.Run("mutation with SQL", func(t *testing.T) {
		cfg := &Config{
			Mutations: []Mutation{
				{SQL: "INSERT INTO users (name) VALUES (:name)"},
			},
		}
		assert.True(t, cfg.RequiresDB())
	})
}

func TestConfig_ValidateTransforms(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := &Config{
			Queries: []Query{
				{
					Transform: &Transform{
						Pre:  "return input",
						Post: PostTransform{All: "return output"},
					},
				},
			},
		}
		assert.NoError(t, cfg.ValidateTransforms())
	})

	t.Run("invalid pre transform", func(t *testing.T) {
		cfg := &Config{
			Queries: []Query{
				{
					Path: "/test",
					Transform: &Transform{
						Pre: "invalid syntax {{{",
					},
				},
			},
		}
		err := cfg.ValidateTransforms()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "queries[0]")
		assert.Contains(t, err.Error(), "pre")
	})

	t.Run("invalid post transform", func(t *testing.T) {
		cfg := &Config{
			Queries: []Query{
				{
					Path: "/test",
					Transform: &Transform{
						Post: PostTransform{All: "invalid syntax {{{"},
					},
				},
			},
		}
		err := cfg.ValidateTransforms()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "queries[0]")
		assert.Contains(t, err.Error(), "post.all")
	})

	t.Run("invalid mock JS", func(t *testing.T) {
		cfg := &Config{
			Queries: []Query{
				{
					Path: "/test",
					Mock: &Mock{
						ObjectJS: "invalid syntax {{{",
					},
				},
			},
		}
		err := cfg.ValidateTransforms()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "queries[0]")
		assert.Contains(t, err.Error(), "mock")
	})

	t.Run("invalid global helpers", func(t *testing.T) {
		cfg := &Config{
			GlobalHelpers: &GlobalHelpers{
				JS: "invalid syntax {{{",
			},
		}
		err := cfg.ValidateTransforms()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global_helpers")
	})

	t.Run("invalid csv value parser", func(t *testing.T) {
		cfg := &Config{
			CSV: &CSVConfig{
				ValueParser: "invalid syntax {{{",
			},
		}
		err := cfg.ValidateTransforms()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "csv.value_parser")
	})

	t.Run("invalid mutation transform", func(t *testing.T) {
		cfg := &Config{
			Mutations: []Mutation{
				{
					Path: "/test",
					Transform: &Transform{
						Pre: "invalid syntax {{{",
					},
				},
			},
		}
		err := cfg.ValidateTransforms()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mutations[0]")
	})

	t.Run("invalid filter", func(t *testing.T) {
		cfg := &Config{
			Queries: []Query{
				{
					Path: "/test",
					Mock: &Mock{
						Array:  []any{map[string]any{"id": 1}},
						Filter: "invalid syntax {{{",
					},
				},
			},
		}
		err := cfg.ValidateTransforms()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "filter")
	})
}

func TestParse(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		yaml := `
dsn: postgres://localhost:5432/db
queries:
  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
`
		cfg, err := Parse([]byte(yaml))
		require.NoError(t, err)
		assert.Equal(t, "postgres://localhost:5432/db", cfg.DSN)
		assert.Len(t, cfg.Queries, 1)
		assert.Equal(t, "/user", cfg.Queries[0].Path)
	})

	t.Run("valid mock config", func(t *testing.T) {
		yaml := `
queries:
  - type: one
    path: /user
    mock:
      object:
        id: 1
        name: Alice
`
		cfg, err := Parse([]byte(yaml))
		require.NoError(t, err)
		assert.Len(t, cfg.Queries, 1)
		assert.NotNil(t, cfg.Queries[0].Mock)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		yaml := `invalid: yaml: syntax`
		_, err := Parse([]byte(yaml))
		require.Error(t, err)
	})

	t.Run("schema validation error - missing type", func(t *testing.T) {
		yaml := `
queries:
  - path: /user
    sql: SELECT * FROM users
`
		_, err := Parse([]byte(yaml))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "type")
	})

	t.Run("schema validation error - invalid type", func(t *testing.T) {
		yaml := `
queries:
  - type: invalid
    path: /user
    sql: SELECT * FROM users
`
		_, err := Parse([]byte(yaml))
		require.Error(t, err)
	})

	t.Run("schema validation error - both sql and mock", func(t *testing.T) {
		yaml := `
queries:
  - type: one
    path: /user
    sql: SELECT * FROM users
    mock:
      object:
        id: 1
`
		_, err := Parse([]byte(yaml))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sql")
		assert.Contains(t, err.Error(), "mock")
	})

	t.Run("schema validation error - neither sql nor mock", func(t *testing.T) {
		yaml := `
dsn: postgres://localhost:5432/db
queries:
  - type: one
    path: /user
`
		_, err := Parse([]byte(yaml))
		require.Error(t, err)
	})

	t.Run("env substitution", func(t *testing.T) {
		t.Setenv("TEST_DB_HOST", "myhost")
		t.Setenv("TEST_DB_PORT", "5433")
		yaml := `
dsn: postgres://${TEST_DB_HOST}:${TEST_DB_PORT}/db
queries:
  - type: one
    path: /user
    sql: SELECT 1
`
		cfg, err := Parse([]byte(yaml))
		require.NoError(t, err)
		assert.Equal(t, "postgres://myhost:5433/db", cfg.DSN)
	})

	t.Run("env substitution with default", func(t *testing.T) {
		yaml := `
dsn: postgres://${UNDEFINED_HOST:-localhost}:5432/db
queries:
  - type: one
    path: /user
    sql: SELECT 1
`
		cfg, err := Parse([]byte(yaml))
		require.NoError(t, err)
		assert.Equal(t, "postgres://localhost:5432/db", cfg.DSN)
	})

	// Note: Transform JS validation happens at server startup via ValidateTransforms,
	// not during Parse. Parse only validates the YAML schema structure.
}
