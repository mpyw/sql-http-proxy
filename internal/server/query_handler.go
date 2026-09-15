// The query handler is one unit with the base handler in handler.go: it
// assembles a baseHandler field by field. See the namespace note there.
//
//declscope:namespace handler
package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/executor"
	"github.com/mpyw/sql-http-proxy/internal/server/body"
)

// QueryHandler handles HTTP requests for queries.
type QueryHandler = baseHandler[executor.ExecuteResult]

// queryHandlerProcessor implements HandlerResultProcessor for ExecuteResult.
type queryHandlerProcessor struct{}

func (p *queryHandlerProcessor) BuildResponse(result *executor.ExecuteResult) (int, any, http.Header) {
	status := lo.CoalesceOrEmpty(result.Status, http.StatusOK)
	return status, result.Output, result.ResponseHeader
}

// NewQueryHandler creates a new QueryHandler.
// db can be nil if mock is configured.
func NewQueryHandler(db *sqlx.DB, query config.Query) (*QueryHandler, error) {
	return newQueryHandlerWithOptions(db, query, handlerOptions{})
}

// newQueryHandlerWithOptions creates a new QueryHandler with options.
// db can be nil if mock is configured.
// Shared on purpose: server.go (NewServeMux) builds every query route here.
//
//declscope:package
func newQueryHandlerWithOptions(db *sqlx.DB, query config.Query, opts handlerOptions) (*QueryHandler, error) {
	execOpts := executor.CompileTransformOptions{
		ConfigDir:   opts.ConfigDir,
		Helpers:     opts.Helpers,
		ValueParser: opts.ValueParser,
	}
	exec, err := executor.NewQueryExecutor(db, query, execOpts)
	if err != nil {
		return nil, err
	}

	// Parse delay if specified
	var delay time.Duration
	if query.Delay != "" {
		delay, err = time.ParseDuration(query.Delay)
		if err != nil {
			return nil, fmt.Errorf("query %s: invalid delay %q: %w", query.Path, query.Delay, err)
		}
	}

	return &QueryHandler{
		exec:          exec,
		method:        query.GetMethod(),
		pathParams:    handlerPathParams(query.Path),
		parser:        body.NewParser(query.GetAccepts()),
		recorder:      opts.Recorder,
		processor:     &queryHandlerProcessor{},
		checkNotFound: true,
		delay:         delay,
	}, nil
}
