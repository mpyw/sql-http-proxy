package server

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/executor"
	"github.com/mpyw/sql-http-proxy/internal/js"
	"github.com/mpyw/sql-http-proxy/internal/server/body"
)

// MutationHandler handles HTTP requests for mutations.
type MutationHandler struct {
	exec     *executor.MutationExecutor
	method   string
	parser   *body.Parser
	recorder QueryRecorder
}

// NewMutationHandler creates a new MutationHandler.
// db can be nil if mock is configured.
func NewMutationHandler(db *sqlx.DB, mutation config.Mutation) (*MutationHandler, error) {
	return NewMutationHandlerWithOptions(db, mutation, HandlerOptions{})
}

// NewMutationHandlerWithOptions creates a new MutationHandler with options.
// db can be nil if mock is configured.
func NewMutationHandlerWithOptions(db *sqlx.DB, mutation config.Mutation, opts HandlerOptions) (*MutationHandler, error) {
	exec, err := executor.NewMutationExecutor(db, mutation, opts.ConfigDir)
	if err != nil {
		return nil, err
	}

	return &MutationHandler{
		exec:     exec,
		method:   mutation.GetMethod(),
		parser:   body.NewParser(mutation.GetAccepts()),
		recorder: opts.Recorder,
	}, nil
}

// ServeHTTP implements http.Handler.
func (h *MutationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != h.method {
		w.Header().Set("Allow", h.method)
		res := &responder{w: w}
		res.Error(http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	slog.Info("Accepting request", "method", r.Method, "uri", r.RequestURI)
	res := &responder{w: w}

	params, err := h.parser.Parse(r)
	if err != nil {
		if errors.Is(err, body.ErrUnsupportedMediaType) {
			res.Error(http.StatusUnsupportedMediaType, err)
		} else {
			res.Error(http.StatusBadRequest, err)
		}
		return
	}

	var execOpts executor.Options
	if h.recorder != nil {
		execOpts.Recorder = func(rec executor.Record) {
			h.recorder(QueryRecord{
				Type:   rec.Type,
				IsMock: rec.IsMock,
				SQL:    rec.SQL,
				Params: rec.Params,
				Args:   rec.Args,
			})
		}
	}

	result, err := h.exec.Execute(r.Context(), params, execOpts)
	if err != nil {
		h.handleError(res, err)
		return
	}

	if result.NoContent {
		res.Respond(http.StatusNoContent, nil)
	} else {
		res.Respond(http.StatusOK, result.Data)
	}
}

func (h *MutationHandler) handleError(res *responder, err error) {
	var transformErr *js.TransformError
	if errors.As(err, &transformErr) {
		status := transformErr.Status
		if status == 0 {
			status = defaultStatusForPhase(err)
		}
		res.Respond(status, transformErr.Body)
		return
	}
	res.Error(http.StatusInternalServerError, err)
}
