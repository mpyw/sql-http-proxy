package executor

import (
	"context"
	"errors"
	"log/slog"
	"maps"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/db"
)

// MutationResult holds the result of a mutation execution.
type MutationResult struct {
	Data      any
	NoContent bool // true if result should be 204 No Content
}

// MutationExecutor executes mutations (mock or DB).
type MutationExecutor struct {
	db         *sqlx.DB
	mutation   config.Mutation
	transforms *Transforms
}

// NewMutationExecutor creates a new MutationExecutor with pre-compiled transforms.
func NewMutationExecutor(db *sqlx.DB, mutation config.Mutation, opts CompileTransformOptions) (*MutationExecutor, error) {
	transforms, err := CompileTransforms(mutation.Transform, opts)
	if err != nil {
		return nil, err
	}

	return &MutationExecutor{
		db:         db,
		mutation:   mutation,
		transforms: transforms,
	}, nil
}

// Execute runs the mutation with the given parameters.
func (e *MutationExecutor) Execute(reqCtx context.Context, params map[string]any, opts Options) (*MutationResult, error) {
	// Shared context for pre/mock/post transforms (mutable)
	ctx := make(map[string]any)

	// Keep original params for post-transform
	originalParams := maps.Clone(params)

	// SQL can be modified by pre-transform
	currentSQL := e.mutation.SQL

	// Apply pre-transform
	if e.transforms.Pre != nil {
		result, err := e.transforms.Pre.ApplyPre(ctx, currentSQL, params)
		if err != nil {
			return nil, WrapPreError(err)
		}
		params = result.Output
		currentSQL = result.SQL
		ctx = result.Ctx
	}

	// Execute mock or DB
	if e.transforms.Mock != nil {
		return e.executeMock(ctx, currentSQL, params, originalParams, opts)
	}
	return e.executeDB(reqCtx, ctx, currentSQL, params, originalParams, opts)
}

func (e *MutationExecutor) executeMock(ctx map[string]any, sqlStr string, params, originalParams map[string]any, opts Options) (*MutationResult, error) {
	slog.Debug("Executing mock mutation", "type", e.mutation.Type, "sql", sqlStr, "params", params)

	if opts.Recorder != nil {
		opts.Recorder(Record{
			Type:   string(e.mutation.Type),
			IsMock: true,
			SQL:    sqlStr,
			Params: params,
		})
	}

	mockOutput, newCtx, err := e.transforms.Mock.Data(ctx, sqlStr, params)
	if err != nil {
		return nil, WrapMockError(err)
	}
	if newCtx != nil {
		ctx = newCtx
	}

	switch e.mutation.Type {
	case config.MutationTypeNone:
		if mockOutput != nil {
			slog.Warn("Mock for 'none' mutation returned a value, but it will be ignored")
		}
		return &MutationResult{NoContent: true}, nil
	case config.MutationTypeOne:
		return e.processOneResult(ctx, originalParams, mockOutput)
	case config.MutationTypeMany:
		return e.processManyResult(ctx, originalParams, mockOutput)
	default:
		return nil, errors.New("unsupported mutation type")
	}
}

func (e *MutationExecutor) executeDB(reqCtx context.Context, ctx map[string]any, currentSQL string, params, originalParams map[string]any, opts Options) (*MutationResult, error) {
	boundSQL, args, err := db.BindParams(e.db, currentSQL, params)
	if err != nil {
		return nil, err
	}

	if opts.Recorder != nil {
		opts.Recorder(Record{
			Type:   string(e.mutation.Type),
			IsMock: false,
			SQL:    boundSQL,
			Args:   args,
		})
	}

	// For MySQL, use Exec to get lastInsertId/rowsAffected
	// For other drivers, use Query to get RETURNING results
	isMySQL := e.db.DriverName() == "mysql"

	switch e.mutation.Type {
	case config.MutationTypeNone:
		slog.Debug("Executing mutation", "type", "none", "sql", boundSQL, "args", args)
		_, err := db.Exec(reqCtx, e.db, boundSQL, args...)
		if err != nil {
			return nil, err
		}
		return &MutationResult{NoContent: true}, nil

	case config.MutationTypeOne:
		slog.Debug("Executing mutation", "type", "one", "sql", boundSQL, "args", args)
		if isMySQL {
			execResult, err := db.Exec(reqCtx, e.db, boundSQL, args...)
			if err != nil {
				return nil, err
			}
			// Set rowsAffected in ctx (always available)
			ctx["rowsAffected"] = execResult.RowsAffected
			// Set lastInsertId in ctx (MySQL only)
			if execResult.LastInsertId != nil {
				ctx["lastInsertId"] = *execResult.LastInsertId
			}
			return e.processOneResult(ctx, originalParams, nil)
		}
		result, err := db.QueryOne(reqCtx, e.db, boundSQL, args...)
		if err != nil {
			return nil, err
		}
		// Explicitly convert typed nil to untyped nil for interface comparison
		var output any
		if result != nil {
			output = result
		}
		return e.processOneResult(ctx, originalParams, output)

	case config.MutationTypeMany:
		slog.Debug("Executing mutation", "type", "many", "sql", boundSQL, "args", args)
		if isMySQL {
			execResult, err := db.Exec(reqCtx, e.db, boundSQL, args...)
			if err != nil {
				return nil, err
			}
			// Set rowsAffected in ctx (always available)
			ctx["rowsAffected"] = execResult.RowsAffected
			// Set lastInsertId in ctx (MySQL only)
			if execResult.LastInsertId != nil {
				ctx["lastInsertId"] = *execResult.LastInsertId
			}
			return e.processManyResult(ctx, originalParams, nil)
		}
		results, err := db.QueryMany(reqCtx, e.db, boundSQL, args...)
		if err != nil {
			return nil, err
		}
		return e.processManyResult(ctx, originalParams, results)

	default:
		return nil, errors.New("unsupported mutation type")
	}
}

func (e *MutationExecutor) processOneResult(ctx, originalParams map[string]any, output any) (*MutationResult, error) {
	var entry map[string]any
	if output == nil {
		entry = nil
	} else if m, ok := output.(map[string]any); ok {
		entry = m
	} else {
		return nil, errors.New("mutation result for 'one' must be object or null")
	}

	if e.transforms.PostAll != nil {
		postResult, err := e.transforms.PostAll.ApplyPost(ctx, originalParams, entry)
		if err != nil {
			return nil, WrapPostError(err)
		}
		if postResult.Output == nil {
			return &MutationResult{NoContent: true}, nil
		}
		return &MutationResult{Data: postResult.Output}, nil
	}

	if entry == nil {
		return &MutationResult{NoContent: true}, nil
	}
	return &MutationResult{Data: entry}, nil
}

func (e *MutationExecutor) processManyResult(ctx, originalParams map[string]any, output any) (*MutationResult, error) {
	var entries []map[string]any

	switch v := output.(type) {
	case []map[string]any:
		entries = v
	case []any:
		entries = make([]map[string]any, len(v))
		for i, item := range v {
			if m, ok := item.(map[string]any); ok {
				entries[i] = m
			} else {
				return nil, errors.New("mutation result for 'many' must be array of objects")
			}
		}
	case nil:
		entries = []map[string]any{}
	default:
		return nil, errors.New("mutation result for 'many' must be array")
	}

	var result any = entries
	currentEntries := entries
	currentCtx := ctx

	if e.transforms.PostEach != nil {
		eachResult, newCtx, err := e.transforms.PostEach.ApplyPostToEachRow(currentCtx, originalParams, currentEntries)
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

	if e.transforms.PostAll != nil {
		postResult, err := e.transforms.PostAll.ApplyPostToAllRows(currentCtx, originalParams, currentEntries)
		if err != nil {
			return nil, WrapPostError(err)
		}
		result = postResult.Output
	}

	return &MutationResult{Data: result}, nil
}
