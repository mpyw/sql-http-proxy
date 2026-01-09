// Package executor provides query and mutation execution logic.
package executor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/js"
	"github.com/mpyw/sql-http-proxy/internal/mock"
)

// ErrNotFound is returned when a "one" type query finds no rows.
var ErrNotFound = errors.New("not found")

// PhaseError wraps an error with the phase (pre/mock/post) it occurred in.
type PhaseError struct {
	Phase string // "pre", "mock", or "post"
	Err   error
}

func (e *PhaseError) Error() string {
	return e.Err.Error()
}

func (e *PhaseError) Unwrap() error {
	return e.Err
}

// WrapPreError wraps an error as a pre-transform error.
func WrapPreError(err error) error {
	if err == nil {
		return nil
	}
	return &PhaseError{Phase: "pre", Err: err}
}

// WrapMockError wraps an error as a mock-transform error.
func WrapMockError(err error) error {
	if err == nil {
		return nil
	}
	return &PhaseError{Phase: "mock", Err: err}
}

// WrapPostError wraps an error as a post-transform error.
func WrapPostError(err error) error {
	if err == nil {
		return nil
	}
	return &PhaseError{Phase: "post", Err: err}
}

// Record represents an executed query/mutation record for testing/logging.
type Record struct {
	Type   config.OpType  // OpTypeOne, OpTypeMany, or OpTypeNone
	IsMock bool           // true if mock execution
	SQL    string         // executed SQL (bound for DB, original for mock)
	Params map[string]any // named parameters (for mock)
	Args   []any          // positional arguments (for DB)
}

// Recorder is called when a query/mutation is executed.
type Recorder func(record Record)

// Options contains optional settings for executors.
type Options struct {
	Recorder    Recorder
	HTTPRequest *http.Request // HTTP request for transform context
}

// ExecuteResult holds the result of an execution including response metadata.
type ExecuteResult struct {
	Output         any         // The response data
	Status         int         // HTTP status code (0 means default)
	ResponseHeader http.Header // Response headers to set
}

// ResultBuilder builds the final result from processed output.
type ResultBuilder[R any] interface {
	// BuildResult creates the result from output and transform context.
	BuildResult(output any, tc *js.TransformContext) *R
	// OnEmptyOne handles empty result for "one" type. Returns (result, error).
	// If error is non-nil, it's returned. If result is non-nil, it's returned.
	OnEmptyOne(tc *js.TransformContext) (*R, error)
}

// BaseExecutor contains common executor logic.
type BaseExecutor[R any] struct {
	DB         *sqlx.DB
	SQL        string
	OpType     config.OpType
	Entity     config.EntityType
	Transforms *Transforms
	MockSource mock.Source
	Builder    ResultBuilder[R]
}

// ExecContext holds the execution context for DB operations.
type ExecContext[R any] struct {
	Base           *BaseExecutor[R]
	Ctx            map[string]any
	SQL            string
	Params         map[string]any
	OriginalParams map[string]any
	Opts           Options
	TC             *js.TransformContext
}

// DBExecutor executes database operations.
type DBExecutor[R any] interface {
	ExecuteDB(reqCtx context.Context, ec *ExecContext[R]) (*R, error)
}

// ExecuteBase runs the common execution flow.
func (e *BaseExecutor[R]) ExecuteBase(reqCtx context.Context, params map[string]any, opts Options, dbExec DBExecutor[R]) (*R, error) {
	ec := &ExecContext[R]{
		Base:           e,
		Ctx:            make(map[string]any),
		SQL:            e.SQL,
		Params:         params,
		OriginalParams: maps.Clone(params),
		Opts:           opts,
		TC:             js.NewTransformContext(opts.HTTPRequest),
	}

	// Apply pre-transform
	if e.Transforms.Pre != nil {
		result, err := e.Transforms.Pre.ApplyPre(ec.Ctx, ec.SQL, ec.Params, ec.TC)
		if err != nil {
			return nil, WrapPreError(err)
		}
		ec.Params = result.Output
		ec.SQL = result.SQL
		ec.Ctx = result.Ctx
	}

	// Execute mock or DB
	if e.MockSource != nil {
		return e.executeMock(ec)
	}
	return dbExec.ExecuteDB(reqCtx, ec)
}

func (e *BaseExecutor[R]) executeMock(ec *ExecContext[R]) (*R, error) {
	slog.Debug(fmt.Sprintf("Executing mock %s", e.Entity), "type", e.OpType, "sql", ec.SQL, "params", ec.Params)

	if ec.Opts.Recorder != nil {
		ec.Opts.Recorder(Record{
			Type:   e.OpType,
			IsMock: true,
			SQL:    ec.SQL,
			Params: ec.Params,
		})
	}

	mockOutput, newCtx, err := e.MockSource.Data(ec.Ctx, ec.SQL, ec.Params, ec.TC)
	if err != nil {
		return nil, WrapMockError(err)
	}
	if newCtx != nil {
		ec.Ctx = newCtx
	}

	switch e.OpType {
	case config.OpTypeOne:
		return e.ProcessOneResult(ec.Ctx, ec.OriginalParams, mockOutput, ec.TC)
	case config.OpTypeMany:
		return e.ProcessManyResult(ec.Ctx, ec.OriginalParams, mockOutput, ec.TC)
	default:
		return nil, fmt.Errorf("unsupported %s type: %s", e.Entity, e.OpType)
	}
}

// ProcessOneResult processes a single row result.
func (e *BaseExecutor[R]) ProcessOneResult(ctx, originalParams map[string]any, output any, tc *js.TransformContext) (*R, error) {
	entry, isEmpty, err := e.convertToEntry(output)
	if err != nil {
		return nil, err
	}

	// Check empty before post transform (for ErrNotFound)
	if isEmpty && e.Transforms.PostAll == nil {
		return e.Builder.OnEmptyOne(tc)
	}

	result, err := e.applyPostOne(ctx, originalParams, entry, tc)
	if err != nil {
		return nil, err
	}

	// Check empty after post transform (result could still be nil)
	if result == nil && isEmpty {
		return e.Builder.OnEmptyOne(tc)
	}

	return e.Builder.BuildResult(result, tc), nil
}

// ProcessManyResult processes multiple row results.
func (e *BaseExecutor[R]) ProcessManyResult(ctx, originalParams map[string]any, output any, tc *js.TransformContext) (*R, error) {
	entries, err := e.convertToEntries(output)
	if err != nil {
		return nil, err
	}

	result, err := e.applyPostMany(ctx, originalParams, entries, tc)
	if err != nil {
		return nil, err
	}

	return e.Builder.BuildResult(result, tc), nil
}

func (e *BaseExecutor[R]) applyPostOne(ctx, params map[string]any, entry map[string]any, tc *js.TransformContext) (any, error) {
	if e.Transforms.PostAll == nil {
		if entry == nil {
			return nil, nil
		}
		return entry, nil
	}
	postResult, err := e.Transforms.PostAll.ApplyPost(ctx, params, entry, tc)
	if err != nil {
		return nil, WrapPostError(err)
	}
	return postResult.Output, nil
}

func (e *BaseExecutor[R]) applyPostMany(ctx, params map[string]any, entries []map[string]any, tc *js.TransformContext) (any, error) {
	var result any = entries
	currentEntries := entries
	currentCtx := ctx

	if e.Transforms.PostEach != nil {
		eachResult, newCtx, err := e.Transforms.PostEach.ApplyPostToEachRow(currentCtx, params, currentEntries, tc)
		if err != nil {
			return nil, WrapPostError(err)
		}
		result = eachResult
		currentCtx = newCtx

		currentEntries = make([]map[string]any, len(eachResult))
		for i, r := range eachResult {
			m, ok := r.(map[string]any)
			if !ok {
				return nil, WrapPostError(errors.New("post.each transform must return object"))
			}
			currentEntries[i] = m
		}
	}

	if e.Transforms.PostAll != nil {
		postResult, err := e.Transforms.PostAll.ApplyPostToAllRows(currentCtx, params, currentEntries, tc)
		if err != nil {
			return nil, WrapPostError(err)
		}
		result = postResult.Output
	}

	return result, nil
}

func (e *BaseExecutor[R]) convertToEntry(output any) (entry map[string]any, isEmpty bool, err error) {
	switch v := output.(type) {
	case nil:
		return nil, true, nil
	case map[string]any:
		return v, false, nil
	case []any:
		if len(v) == 0 {
			return nil, true, nil
		}
		if m, ok := v[0].(map[string]any); ok {
			return m, false, nil
		}
		return nil, false, fmt.Errorf("%s result for 'one' must be object or null", e.Entity)
	default:
		return nil, false, fmt.Errorf("%s result for 'one' must be object or null", e.Entity)
	}
}

func (e *BaseExecutor[R]) convertToEntries(output any) ([]map[string]any, error) {
	switch v := output.(type) {
	case []map[string]any:
		return v, nil
	case []any:
		entries := make([]map[string]any, len(v))
		for i, item := range v {
			if m, ok := item.(map[string]any); ok {
				entries[i] = m
			} else {
				return nil, fmt.Errorf("%s result for 'many' must be array of objects", e.Entity)
			}
		}
		return entries, nil
	case nil:
		return []map[string]any{}, nil
	default:
		return nil, fmt.Errorf("%s result for 'many' must be array", e.Entity)
	}
}
