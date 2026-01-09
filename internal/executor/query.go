package executor

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/db"
	"github.com/mpyw/sql-http-proxy/internal/js"
)

// queryResultBuilder implements ResultBuilder for ExecuteResult.
type queryResultBuilder struct {
	handleNotFound bool
}

func (b *queryResultBuilder) BuildResult(output any, tc *js.TransformContext) *ExecuteResult {
	return &ExecuteResult{
		Output:         output,
		Status:         tc.Response.Status(),
		ResponseHeader: tc.Response.ToHTTPHeader(),
	}
}

func (b *queryResultBuilder) OnEmptyOne(tc *js.TransformContext) (*ExecuteResult, error) {
	if !b.handleNotFound {
		return nil, ErrNotFound
	}
	return &ExecuteResult{
		Output:         nil,
		Status:         tc.Response.Status(),
		ResponseHeader: tc.Response.ToHTTPHeader(),
	}, nil
}

// QueryExecutor executes queries (mock or DB).
type QueryExecutor struct {
	base  *BaseExecutor[ExecuteResult]
	query config.Query
}

// NewQueryExecutor creates a new QueryExecutor with pre-compiled transforms.
func NewQueryExecutor(db *sqlx.DB, query config.Query, opts CompileTransformOptions) (*QueryExecutor, error) {
	transforms, err := CompileTransforms(query.Transform, opts)
	if err != nil {
		return nil, err
	}

	mockSource, err := CompileMock(query.Mock, opts)
	if err != nil {
		return nil, err
	}

	return &QueryExecutor{
		base: &BaseExecutor[ExecuteResult]{
			DB:         db,
			SQL:        query.SQL,
			OpType:     config.OpType(query.Type),
			Entity:     config.EntityQuery,
			Transforms: transforms,
			MockSource: mockSource,
			Builder:    &queryResultBuilder{handleNotFound: query.HandleNotFound},
		},
		query: query,
	}, nil
}

// Execute runs the query with the given parameters.
func (e *QueryExecutor) Execute(reqCtx context.Context, params map[string]any, opts Options) (*ExecuteResult, error) {
	return e.base.ExecuteBase(reqCtx, params, opts, e)
}

// ExecuteDB implements DBExecutor interface.
func (e *QueryExecutor) ExecuteDB(reqCtx context.Context, ec *ExecContext[ExecuteResult]) (*ExecuteResult, error) {
	boundSQL, args, err := db.BindParams(ec.Base.DB, ec.SQL, ec.Params)
	if err != nil {
		return nil, err
	}

	if ec.Opts.Recorder != nil {
		ec.Opts.Recorder(Record{
			Type:   config.OpType(e.query.Type),
			IsMock: false,
			SQL:    boundSQL,
			Args:   args,
		})
	}

	switch e.query.Type {
	case config.QueryTypeOne:
		slog.Debug("Executing query", "type", "one", "sql", boundSQL, "args", args)
		result, err := db.QueryOne(reqCtx, ec.Base.DB, boundSQL, args...)
		if err != nil {
			return nil, err
		}
		var output any
		if result != nil {
			output = result
		}
		return ec.Base.ProcessOneResult(ec.Ctx, ec.OriginalParams, output, ec.TC)

	case config.QueryTypeMany:
		slog.Debug("Executing query", "type", "many", "sql", boundSQL, "args", args)
		results, err := db.QueryMany(reqCtx, ec.Base.DB, boundSQL, args...)
		if err != nil {
			return nil, err
		}
		return ec.Base.ProcessManyResult(ec.Ctx, ec.OriginalParams, results, ec.TC)

	default:
		return nil, ErrUnsupportedQueryType
	}
}

// ErrUnsupportedQueryType is returned when query type is not supported.
var ErrUnsupportedQueryType = errors.New("unsupported query type")
