package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/js"
	"github.com/mpyw/sql-http-proxy/internal/mock"
)

// Regex pattern shorthands for path parameters.
// These are expanded when the pattern starts and ends with *.
var regexShorthands = map[string]string{
	"uuid":    "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
	"uuid_v4": "[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}",
	"uuid_v7": "[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}",
}

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
	opts := handlerOptions{
		ConfigDir:   configDir,
		Helpers:     helpers,
		ValueParser: valueParser,
	}

	for _, query := range cfg.Queries {
		handler, err := newQueryHandlerWithOptions(db, query, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to create handler for %s: %w", query.Path, err)
		}
		r.Method(query.GetMethod(), expandPathShorthands(query.Path), handler)
	}

	for _, mutation := range cfg.Mutations {
		handler, err := newMutationHandlerWithOptions(db, mutation, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to create handler for %s: %w", mutation.Path, err)
		}
		r.Method(mutation.GetMethod(), expandPathShorthands(mutation.Path), handler)
	}

	r.NotFound(createNotFoundHandler().ServeHTTP)

	return r, nil
}

// expandPathShorthands expands regex shorthand notations in path patterns.
// Patterns like {id:*uuid*} are expanded to {id:[0-9a-f]{8}-...}.
// This is unambiguous because * at the start of a regex is invalid.
func expandPathShorthands(path string) string {
	result := strings.Builder{}
	i := 0
	for i < len(path) {
		if path[i] == '{' {
			// Find the closing brace
			start := i
			for i < len(path) && path[i] != '}' {
				i++
			}
			if i < len(path) {
				i++ // include closing brace
			}
			segment := path[start:i]

			// Check if segment has a shorthand pattern (e.g., {id:*uuid*})
			if colonIdx := strings.Index(segment, ":*"); colonIdx != -1 {
				// Extract the part after :* and before *}
				afterColon := segment[colonIdx+2:] // skip ":*"
				if strings.HasSuffix(afterColon, "*}") {
					shorthand := afterColon[:len(afterColon)-2] // remove "*}"
					if regex, ok := regexShorthands[shorthand]; ok {
						// Expand: {id:*uuid*} -> {id:regex}
						paramName := segment[1:colonIdx] // extract "id" from "{id:*uuid*}"
						result.WriteString("{")
						result.WriteString(paramName)
						result.WriteString(":")
						result.WriteString(regex)
						result.WriteString("}")
						continue
					}
				}
			}
			// No expansion needed, write as-is
			result.WriteString(segment)
		} else {
			result.WriteByte(path[i])
			i++
		}
	}
	return result.String()
}
