package server

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

// NewServeMux creates a new HTTP ServeMux with all configured handlers.
// db can be nil if all endpoints use mock.
func NewServeMux(db *sqlx.DB, cfg config.Config) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	for _, query := range cfg.Queries {
		handler, err := CreateHandler(db, query)
		if err != nil {
			return nil, fmt.Errorf("failed to create handler for %s: %w", query.Path, err)
		}
		mux.Handle(query.Path, handler)
	}

	for _, mutation := range cfg.Mutations {
		handler, err := CreateMutationHandler(db, mutation)
		if err != nil {
			return nil, fmt.Errorf("failed to create handler for %s: %w", mutation.Path, err)
		}
		mux.Handle(fmt.Sprintf("%s %s", mutation.GetMethod(), mutation.Path), handler)
	}

	mux.Handle("/", CreateNotFoundHandler())

	return mux, nil
}
