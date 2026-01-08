package server

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/executor"
	"github.com/mpyw/sql-http-proxy/internal/server/body"
)

// QueryHandler handles HTTP requests for queries.
type QueryHandler = BaseHandler[executor.ExecuteResult]

// queryResultProcessor implements ResultProcessor for ExecuteResult.
type queryResultProcessor struct{}

func (p *queryResultProcessor) BuildResponse(result *executor.ExecuteResult) (int, any, http.Header) {
	status := lo.CoalesceOrEmpty(result.Status, http.StatusOK)
	return status, result.Output, result.ResponseHeader
}

// NewQueryHandler creates a new QueryHandler.
// db can be nil if mock is configured.
func NewQueryHandler(db *sqlx.DB, query config.Query) (*QueryHandler, error) {
	return NewQueryHandlerWithOptions(db, query, HandlerOptions{})
}

// NewQueryHandlerWithOptions creates a new QueryHandler with options.
// db can be nil if mock is configured.
func NewQueryHandlerWithOptions(db *sqlx.DB, query config.Query, opts HandlerOptions) (*QueryHandler, error) {
	execOpts := executor.CompileTransformOptions{
		ConfigDir:   opts.ConfigDir,
		Helpers:     opts.Helpers,
		ValueParser: opts.ValueParser,
	}
	exec, err := executor.NewQueryExecutor(db, query, execOpts)
	if err != nil {
		return nil, err
	}

	return &QueryHandler{
		exec:          exec,
		method:        query.GetMethod(),
		pathParams:    extractPathParams(query.Path),
		parser:        body.NewParser(query.GetAccepts()),
		recorder:      opts.Recorder,
		processor:     &queryResultProcessor{},
		checkNotFound: true,
	}, nil
}
