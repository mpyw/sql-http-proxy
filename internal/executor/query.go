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

// QueryExecutor executes queries (mock or DB).
type QueryExecutor struct {
	db         *sqlx.DB
	query      config.Query
	transforms *Transforms
}

// NewQueryExecutor creates a new QueryExecutor with pre-compiled transforms.
// configDir is the directory of the YAML config file (for resolving relative paths in mock).
func NewQueryExecutor(db *sqlx.DB, query config.Query, configDir string) (*QueryExecutor, error) {
	transforms, err := CompileTransforms(query.Transform, configDir)
	if err != nil {
		return nil, err
	}

	return &QueryExecutor{
		db:         db,
		query:      query,
		transforms: transforms,
	}, nil
}

// IsMock returns true if this executor uses mock instead of DB.
func (e *QueryExecutor) IsMock() bool {
	return e.transforms.IsMock()
}

// Execute runs the query with the given parameters.
// Returns ErrNotFound for "one" type queries when no row is found and handle_not_found is false.
func (e *QueryExecutor) Execute(reqCtx context.Context, params map[string]any, opts Options) (any, error) {
	// Shared context for pre/mock/post transforms (mutable)
	ctx := make(map[string]any)

	// Keep original params for post-transform
	originalParams := maps.Clone(params)

	// SQL can be modified by pre-transform
	currentSQL := e.query.SQL

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

func (e *QueryExecutor) executeMock(ctx map[string]any, sqlStr string, params, originalParams map[string]any, opts Options) (any, error) {
	slog.Debug("Executing mock query", "type", e.query.Type, "sql", sqlStr, "params", params)

	if opts.Recorder != nil {
		opts.Recorder(Record{
			Type:   string(e.query.Type),
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

	switch e.query.Type {
	case config.QueryTypeOne:
		return e.processOneResult(ctx, originalParams, mockOutput)
	case config.QueryTypeMany:
		return e.processManyResult(ctx, originalParams, mockOutput)
	default:
		return nil, errors.New("unsupported query type")
	}
}

func (e *QueryExecutor) executeDB(reqCtx context.Context, ctx map[string]any, currentSQL string, params, originalParams map[string]any, opts Options) (any, error) {
	boundSQL, args, err := db.BindParams(e.db, currentSQL, params)
	if err != nil {
		return nil, err
	}

	if opts.Recorder != nil {
		opts.Recorder(Record{
			Type:   string(e.query.Type),
			IsMock: false,
			SQL:    boundSQL,
			Args:   args,
		})
	}

	switch e.query.Type {
	case config.QueryTypeOne:
		slog.Debug("Executing query", "type", "one", "sql", boundSQL, "args", args)
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

	case config.QueryTypeMany:
		slog.Debug("Executing query", "type", "many", "sql", boundSQL, "args", args)
		results, err := db.QueryMany(reqCtx, e.db, boundSQL, args...)
		if err != nil {
			return nil, err
		}
		return e.processManyResult(ctx, originalParams, results)

	default:
		return nil, errors.New("unsupported query type")
	}
}

func (e *QueryExecutor) processOneResult(ctx, originalParams map[string]any, output any) (any, error) {
	var entry map[string]any
	if output == nil {
		if !e.query.HandleNotFound {
			return nil, ErrNotFound
		}
		entry = nil
	} else if m, ok := output.(map[string]any); ok {
		entry = m
	} else {
		return nil, errors.New("query result for 'one' must be object or null")
	}

	var result any = entry
	if e.transforms.PostAll != nil {
		var err error
		result, err = e.transforms.PostAll.Apply(ctx, originalParams, entry)
		if err != nil {
			return nil, WrapPostError(err)
		}
	}
	return result, nil
}

func (e *QueryExecutor) processManyResult(ctx, originalParams map[string]any, output any) (any, error) {
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
				return nil, errors.New("query result for 'many' must be array of objects")
			}
		}
	case nil:
		entries = []map[string]any{}
	default:
		return nil, errors.New("query result for 'many' must be array")
	}

	var result any = entries
	currentEntries := entries

	if e.transforms.PostEach != nil {
		eachResult, err := e.transforms.PostEach.ApplyToEachRow(ctx, originalParams, currentEntries)
		if err != nil {
			return nil, WrapPostError(err)
		}
		result = eachResult

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
		var err error
		result, err = e.transforms.PostAll.ApplyToAllRows(ctx, originalParams, currentEntries)
		if err != nil {
			return nil, WrapPostError(err)
		}
	}

	return result, nil
}
