// The mutation handler is one unit with the base handler in handler.go: it
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

// MutationHandler handles HTTP requests for mutations.
type MutationHandler = baseHandler[executor.MutationResult]

// mutationHandlerProcessor implements HandlerResultProcessor for MutationResult.
type mutationHandlerProcessor struct{}

func (p *mutationHandlerProcessor) BuildResponse(result *executor.MutationResult) (int, any, http.Header) {
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
// Shared on purpose: server.go (NewServeMux) builds every mutation route here.
//
//declscope:package
func newMutationHandlerWithOptions(db *sqlx.DB, mutation config.Mutation, opts handlerOptions) (*MutationHandler, error) {
	execOpts := executor.CompileTransformOptions{
		ConfigDir:   opts.configDir,
		Helpers:     opts.helpers,
		ValueParser: opts.valueParser,
	}
	exec, err := executor.NewMutationExecutor(db, mutation, execOpts)
	if err != nil {
		return nil, err
	}

	// Parse delay if specified
	var delay time.Duration
	if mutation.Delay != "" {
		delay, err = time.ParseDuration(mutation.Delay)
		if err != nil {
			return nil, fmt.Errorf("mutation %s: invalid delay %q: %w", mutation.Path, mutation.Delay, err)
		}
	}

	return &MutationHandler{
		exec:          exec,
		method:        mutation.GetMethod(),
		pathParams:    handlerPathParams(mutation.Path),
		parser:        body.NewParser(mutation.GetAccepts()),
		recorder:      opts.recorder,
		processor:     &mutationHandlerProcessor{},
		checkNotFound: false,
		delay:         delay,
	}, nil
}
