package server

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

// NewServeMux creates a new HTTP ServeMux with all configured handlers.
// db can be nil if all endpoints use mock.
// configDir is the directory of the config file for resolving relative paths.
func NewServeMux(db *sqlx.DB, cfg config.Config, configDir string) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	opts := HandlerOptions{ConfigDir: configDir}

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
