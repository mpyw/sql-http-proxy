package executor

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/db"
	"github.com/mpyw/sql-http-proxy/internal/js"
)

// MutationResult holds the result of a mutation execution.
type MutationResult struct {
	Data           any
	NoContent      bool
	Status         int
	ResponseHeader http.Header
}

// mutationResultBuilder implements ResultBuilder for MutationResult.
type mutationResultBuilder struct{}

func (b *mutationResultBuilder) BuildResult(output any, tc *js.TransformContext) *MutationResult {
	if output == nil {
		return &MutationResult{
			NoContent:      true,
			Status:         tc.Response.Status(),
			ResponseHeader: tc.Response.ToHTTPHeader(),
		}
	}
	return &MutationResult{
		Data:           output,
		Status:         tc.Response.Status(),
		ResponseHeader: tc.Response.ToHTTPHeader(),
	}
}

func (b *mutationResultBuilder) OnEmptyOne(tc *js.TransformContext) (*MutationResult, error) {
	return &MutationResult{
		NoContent:      true,
		Status:         tc.Response.Status(),
		ResponseHeader: tc.Response.ToHTTPHeader(),
	}, nil
}

// MutationExecutor executes mutations (mock or DB).
type MutationExecutor struct {
	base     *BaseExecutor[MutationResult]
	mutation config.Mutation
}

// NewMutationExecutor creates a new MutationExecutor with pre-compiled transforms.
func NewMutationExecutor(db *sqlx.DB, mutation config.Mutation, opts CompileTransformOptions) (*MutationExecutor, error) {
	transforms, err := CompileTransforms(mutation.Transform, opts)
	if err != nil {
		return nil, err
	}

	mockSource, err := CompileMock(mutation.Mock, opts)
	if err != nil {
		return nil, err
	}

	return &MutationExecutor{
		base: &BaseExecutor[MutationResult]{
			DB:         db,
			SQL:        mutation.SQL,
			OpType:     config.OpType(mutation.Type),
			Entity:     config.EntityMutation,
			Transforms: transforms,
			MockSource: mockSource,
			Builder:    &mutationResultBuilder{},
		},
		mutation: mutation,
	}, nil
}

// Execute runs the mutation with the given parameters.
func (e *MutationExecutor) Execute(reqCtx context.Context, params map[string]any, opts Options) (*MutationResult, error) {
	return e.base.ExecuteBase(reqCtx, params, opts, e)
}

// ExecuteDB implements DBExecutor interface.
func (e *MutationExecutor) ExecuteDB(reqCtx context.Context, ec *ExecContext[MutationResult]) (*MutationResult, error) {
	boundSQL, args, err := db.BindParams(ec.Base.DB, ec.SQL, ec.Params)
	if err != nil {
		return nil, err
	}

	if ec.Opts.Recorder != nil {
		ec.Opts.Recorder(Record{
			Type:   config.OpType(e.mutation.Type),
			IsMock: false,
			SQL:    boundSQL,
			Args:   args,
		})
	}

	isMySQL := ec.Base.DB.DriverName() == "mysql"

	switch e.mutation.Type {
	case config.MutationTypeNone:
		slog.Debug("Executing mutation", "type", "none", "sql", boundSQL, "args", args)
		_, err := db.Exec(reqCtx, ec.Base.DB, boundSQL, args...)
		if err != nil {
			return nil, err
		}
		return &MutationResult{
			NoContent:      true,
			Status:         ec.TC.Response.Status(),
			ResponseHeader: ec.TC.Response.ToHTTPHeader(),
		}, nil

	case config.MutationTypeOne:
		slog.Debug("Executing mutation", "type", "one", "sql", boundSQL, "args", args)
		if isMySQL {
			return e.execMySQLOne(reqCtx, ec, boundSQL, args)
		}
		result, err := db.QueryOne(reqCtx, ec.Base.DB, boundSQL, args...)
		if err != nil {
			return nil, err
		}
		var output any
		if result != nil {
			output = result
		}
		return ec.Base.ProcessOneResult(ec.Ctx, ec.OriginalParams, output, ec.TC)

	case config.MutationTypeMany:
		slog.Debug("Executing mutation", "type", "many", "sql", boundSQL, "args", args)
		if isMySQL {
			return e.execMySQLMany(reqCtx, ec, boundSQL, args)
		}
		results, err := db.QueryMany(reqCtx, ec.Base.DB, boundSQL, args...)
		if err != nil {
			return nil, err
		}
		return ec.Base.ProcessManyResult(ec.Ctx, ec.OriginalParams, results, ec.TC)

	default:
		return nil, ErrUnsupportedMutationType
	}
}

func (e *MutationExecutor) execMySQLOne(reqCtx context.Context, ec *ExecContext[MutationResult], boundSQL string, args []any) (*MutationResult, error) {
	execResult, err := db.Exec(reqCtx, ec.Base.DB, boundSQL, args...)
	if err != nil {
		return nil, err
	}
	ec.Ctx["rowsAffected"] = execResult.RowsAffected
	if execResult.LastInsertId != nil {
		ec.Ctx["lastInsertId"] = *execResult.LastInsertId
	}
	return ec.Base.ProcessOneResult(ec.Ctx, ec.OriginalParams, nil, ec.TC)
}

func (e *MutationExecutor) execMySQLMany(reqCtx context.Context, ec *ExecContext[MutationResult], boundSQL string, args []any) (*MutationResult, error) {
	execResult, err := db.Exec(reqCtx, ec.Base.DB, boundSQL, args...)
	if err != nil {
		return nil, err
	}
	ec.Ctx["rowsAffected"] = execResult.RowsAffected
	if execResult.LastInsertId != nil {
		ec.Ctx["lastInsertId"] = *execResult.LastInsertId
	}
	return ec.Base.ProcessManyResult(ec.Ctx, ec.OriginalParams, nil, ec.TC)
}

// ErrUnsupportedMutationType is returned when mutation type is not supported.
var ErrUnsupportedMutationType = errors.New("unsupported mutation type")
