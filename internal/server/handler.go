// Package server provides HTTP server functionality for sql-http-proxy.
package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"

	"github.com/mpyw/sql-http-proxy/internal/executor"
	"github.com/mpyw/sql-http-proxy/internal/js"
	"github.com/mpyw/sql-http-proxy/internal/mock"
	"github.com/mpyw/sql-http-proxy/internal/server/body"
)

// QueryRecord represents an executed query record.
type QueryRecord struct {
	Type   string         // "one" or "many"
	IsMock bool           // true if mock query
	SQL    string         // executed SQL (bound for DB, original for mock)
	Params map[string]any // named parameters (for mock)
	Args   []any          // positional arguments (for DB)
}

// QueryRecorder is called when a query is executed.
type QueryRecorder func(record QueryRecord)

// HandlerOptions contains optional settings for CreateHandler.
type HandlerOptions struct {
	Recorder    QueryRecorder
	ConfigDir   string              // Directory of config file for resolving relative paths
	Helpers     *js.CompiledHelpers // Global JavaScript helpers
	ValueParser *mock.ValueParser   // Global CSV value parser
}

// CreateNotFoundHandler creates a 404 handler.
func CreateNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := &responder{w: w}
		res.Error(http.StatusNotFound, errors.New("not found"))
	})
}

// Executor is the interface for query/mutation executors.
type Executor[R any] interface {
	Execute(ctx context.Context, params map[string]any, opts executor.Options) (*R, error)
}

// ResultProcessor handles result-specific response building.
type ResultProcessor[R any] interface {
	// BuildResponse returns (status, output, responseHeader) from the result.
	BuildResponse(result *R) (status int, output any, headers http.Header)
}

// BaseHandler contains common handler logic.
type BaseHandler[R any] struct {
	exec          Executor[R]
	method        string
	pathParams    []string // path parameter names in order
	parser        *body.Parser
	recorder      QueryRecorder
	processor     ResultProcessor[R]
	checkNotFound bool // true for query handlers
}

// ServeHTTP implements http.Handler.
func (h *BaseHandler[R]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Accepting request", "method", r.Method, "uri", r.RequestURI)
	res := &responder{w: w}

	if r.Method != h.method {
		if h.method != http.MethodGet && h.method != http.MethodHead {
			w.Header().Set("Allow", h.method)
		}
		res.Error(http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	params, err := h.parseParams(r)
	if err != nil {
		var status int
		switch {
		case errors.Is(err, body.ErrBodyTooLarge):
			status = http.StatusRequestEntityTooLarge
		case errors.Is(err, body.ErrUnsupportedMediaType):
			status = http.StatusUnsupportedMediaType
		default:
			status = http.StatusBadRequest
		}
		res.Error(status, err)
		return
	}

	result, err := h.exec.Execute(r.Context(), params, h.buildExecOptions(r))
	if err != nil {
		h.handleError(res, err)
		return
	}

	status, output, headers := h.processor.BuildResponse(result)
	h.applyResponseHeaders(w, headers)
	res.Respond(status, output)
}

func (h *BaseHandler[R]) parseParams(r *http.Request) (map[string]any, error) {
	var params map[string]any
	var err error

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		params = extractNamedParams(r.URL.Query())
	default:
		params, err = h.parser.Parse(r)
		if err != nil {
			return nil, err
		}
	}

	// Merge path parameters (path params take priority over query/body params)
	for _, name := range h.pathParams {
		if value := chi.URLParam(r, name); value != "" {
			params[name] = value
		}
	}

	return params, nil
}

func (h *BaseHandler[R]) handleError(res *responder, err error) {
	if h.checkNotFound && errors.Is(err, executor.ErrNotFound) {
		res.Error(http.StatusNotFound, errors.New("not found"))
		return
	}
	if transformErr := (*js.TransformError)(nil); errors.As(err, &transformErr) {
		status := lo.CoalesceOrEmpty(transformErr.Status, defaultStatusForPhase(err))
		res.Respond(status, transformErr.Body)
		return
	}
	res.Error(http.StatusInternalServerError, err)
}

// buildExecOptions creates executor.Options with optional recorder wrapping.
func (h *BaseHandler[R]) buildExecOptions(r *http.Request) executor.Options {
	opts := executor.Options{
		HTTPRequest: r,
	}
	if h.recorder != nil {
		opts.Recorder = func(rec executor.Record) {
			h.recorder(QueryRecord{
				Type:   rec.Type,
				IsMock: rec.IsMock,
				SQL:    rec.SQL,
				Params: rec.Params,
				Args:   rec.Args,
			})
		}
	}
	return opts
}

// applyResponseHeaders copies response headers to the http.ResponseWriter.
func (h *BaseHandler[R]) applyResponseHeaders(w http.ResponseWriter, headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}

func defaultStatusForPhase(err error) int {
	if phaseErr := (*executor.PhaseError)(nil); errors.As(err, &phaseErr) && phaseErr.Phase == "pre" {
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

// extractPathParams extracts path parameter names from a path pattern.
// Example: /users/:id/posts/:post_id -> ["id", "post_id"]
func extractPathParams(path string) []string {
	var params []string
	i := 0
	for i < len(path) {
		if path[i] == ':' {
			i++
			start := i
			// Read parameter name (alphanumeric and underscore)
			for i < len(path) && (path[i] == '_' || isPathParamChar(path[i])) {
				i++
			}
			if i > start {
				params = append(params, path[start:i])
			}
		} else {
			i++
		}
	}
	return params
}

func isPathParamChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
