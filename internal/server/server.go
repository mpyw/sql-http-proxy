package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/js"
	"github.com/mpyw/sql-http-proxy/internal/mock"
)

// NewServeMux creates a new chi Router with all configured handlers.
// db can be nil if all endpoints use mock.
// configDir is the directory of the config file for resolving relative paths.
func NewServeMux(db *sqlx.DB, cfg config.Config, configDir string) (http.Handler, error) {
	// Compile global helpers once at startup
	var helpers *js.CompiledHelpers
	if cfg.GlobalHelpers != nil {
		var err error
		helpers, err = js.CompileHelpers(cfg.GlobalHelpers.JS, cfg.GlobalHelpers.JSFiles, configDir)
		if err != nil {
			return nil, fmt.Errorf("global_helpers: %w", err)
		}
	}

	// Compile CSV value parser once at startup
	var valueParser *mock.ValueParser
	if cfg.CSV != nil && cfg.CSV.ValueParser != "" {
		var err error
		valueParser, err = mock.CompileValueParser(cfg.CSV.ValueParser, helpers)
		if err != nil {
			return nil, fmt.Errorf("csv.value_parser: %w", err)
		}
	}

	r := chi.NewRouter()
	opts := HandlerOptions{
		ConfigDir:   configDir,
		Helpers:     helpers,
		ValueParser: valueParser,
	}

	for _, query := range cfg.Queries {
		handler, err := NewQueryHandlerWithOptions(db, query, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to create handler for %s: %w", query.Path, err)
		}
		r.Method(query.GetMethod(), convertPathPattern(query.Path), handler)
	}

	for _, mutation := range cfg.Mutations {
		handler, err := NewMutationHandlerWithOptions(db, mutation, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to create handler for %s: %w", mutation.Path, err)
		}
		r.Method(mutation.GetMethod(), convertPathPattern(mutation.Path), handler)
	}

	r.NotFound(CreateNotFoundHandler().ServeHTTP)

	return r, nil
}

// convertPathPattern converts path patterns from :param syntax to {param} syntax for chi.
// Example: /users/:id -> /users/{id}
func convertPathPattern(path string) string {
	result := make([]byte, 0, len(path))
	i := 0
	for i < len(path) {
		if path[i] == ':' {
			// Start of a path parameter
			result = append(result, '{')
			i++
			// Read parameter name (alphanumeric and underscore)
			for i < len(path) && (path[i] == '_' || isAlphanumeric(path[i])) {
				result = append(result, path[i])
				i++
			}
			result = append(result, '}')
		} else {
			result = append(result, path[i])
			i++
		}
	}
	return string(result)
}

func isAlphanumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
