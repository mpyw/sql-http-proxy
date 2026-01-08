package server

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/js"
	"github.com/mpyw/sql-http-proxy/internal/mock"
)

// NewServeMux creates a new HTTP ServeMux with all configured handlers.
// db can be nil if all endpoints use mock.
// configDir is the directory of the config file for resolving relative paths.
func NewServeMux(db *sqlx.DB, cfg config.Config, configDir string) (*http.ServeMux, error) {
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
	if cfg.CSV != nil && cfg.CSV.ValueParser != nil {
		vp := cfg.CSV.ValueParser
		var err error
		valueParser, err = mock.CompileValueParser(vp.JS, vp.JSFile, configDir, helpers)
		if err != nil {
			return nil, fmt.Errorf("csv.value_parser: %w", err)
		}
	}

	mux := http.NewServeMux()
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
		mux.Handle(query.Path, handler)
	}

	for _, mutation := range cfg.Mutations {
		handler, err := NewMutationHandlerWithOptions(db, mutation, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to create handler for %s: %w", mutation.Path, err)
		}
		mux.Handle(fmt.Sprintf("%s %s", mutation.GetMethod(), mutation.Path), handler)
	}

	mux.Handle("/", CreateNotFoundHandler())

	return mux, nil
}
