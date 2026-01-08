package server

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/executor"
	"github.com/mpyw/sql-http-proxy/internal/server/body"
)

// MutationHandler handles HTTP requests for mutations.
type MutationHandler = baseHandler[executor.MutationResult]

// mutationResultProcessor implements ResultProcessor for MutationResult.
type mutationResultProcessor struct{}

func (p *mutationResultProcessor) BuildResponse(result *executor.MutationResult) (int, any, http.Header) {
	// For NoContent, always use 204 (original behavior)
	// result.Status from Response defaults to 200, which should not override 204
	if result.NoContent {
		return http.StatusNoContent, nil, result.ResponseHeader
	}
	status := lo.CoalesceOrEmpty(result.Status, http.StatusOK)
	return status, result.Data, result.ResponseHeader
}

// NewMutationHandler creates a new MutationHandler.
// db can be nil if mock is configured.
func NewMutationHandler(db *sqlx.DB, mutation config.Mutation) (*MutationHandler, error) {
	return newMutationHandlerWithOptions(db, mutation, handlerOptions{})
}

// newMutationHandlerWithOptions creates a new MutationHandler with options.
// db can be nil if mock is configured.
func newMutationHandlerWithOptions(db *sqlx.DB, mutation config.Mutation, opts handlerOptions) (*MutationHandler, error) {
	execOpts := executor.CompileTransformOptions{
		ConfigDir:   opts.ConfigDir,
		Helpers:     opts.Helpers,
		ValueParser: opts.ValueParser,
	}
	exec, err := executor.NewMutationExecutor(db, mutation, execOpts)
	if err != nil {
		return nil, err
	}

	return &MutationHandler{
		exec:          exec,
		method:        mutation.GetMethod(),
		pathParams:    extractPathParams(mutation.Path),
		parser:        body.NewParser(mutation.GetAccepts()),
		recorder:      opts.Recorder,
		processor:     &mutationResultProcessor{},
		checkNotFound: false,
	}, nil
}
