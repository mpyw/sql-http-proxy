package server

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/executor"
	"github.com/mpyw/sql-http-proxy/internal/js"
	"github.com/mpyw/sql-http-proxy/internal/server/body"
)

// QueryHandler handles HTTP requests for queries.
type QueryHandler struct {
	exec     *executor.QueryExecutor
	method   string
	parser   *body.Parser
	recorder QueryRecorder
}

// NewQueryHandler creates a new QueryHandler.
// db can be nil if mock is configured.
func NewQueryHandler(db *sqlx.DB, query config.Query) (*QueryHandler, error) {
	return NewQueryHandlerWithOptions(db, query, HandlerOptions{})
}

// NewQueryHandlerWithOptions creates a new QueryHandler with options.
// db can be nil if mock is configured.
func NewQueryHandlerWithOptions(db *sqlx.DB, query config.Query, opts HandlerOptions) (*QueryHandler, error) {
	exec, err := executor.NewQueryExecutor(db, query, opts.ConfigDir)
	if err != nil {
		return nil, err
	}

	return &QueryHandler{
		exec:     exec,
		method:   query.GetMethod(),
		parser:   body.NewParser(query.GetAccepts()),
		recorder: opts.Recorder,
	}, nil
}

// ServeHTTP implements http.Handler.
func (h *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Accepting request", "uri", r.RequestURI)
	res := &responder{w: w}

	if r.Method != h.method {
		res.Error(http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	// GET/HEAD use query params, others use body
	var params map[string]any
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		params = extractNamedParams(r.URL.Query())
	} else {
		var err error
		params, err = h.parser.Parse(r)
		if err != nil {
			if errors.Is(err, body.ErrUnsupportedMediaType) {
				res.Error(http.StatusUnsupportedMediaType, err)
			} else {
				res.Error(http.StatusBadRequest, err)
			}
			return
		}
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
	res.Respond(http.StatusOK, result)
}

func (h *QueryHandler) handleError(res *responder, err error) {
	if errors.Is(err, executor.ErrNotFound) {
		res.Error(http.StatusNotFound, errors.New("not found"))
		return
	}
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

func defaultStatusForPhase(err error) int {
	var phaseErr *executor.PhaseError
	if errors.As(err, &phaseErr) && phaseErr.Phase == "pre" {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func extractNamedParams(values url.Values) map[string]any {
	params := make(map[string]any, len(values))
	for key, vals := range values {
		if len(vals) > 0 {
			params[key] = vals[0]
		}
	}
	return params
}
