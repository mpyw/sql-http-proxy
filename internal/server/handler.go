// Namespace note: baseHandler is one type assembled across three files - the
// generic machinery here, and the query/mutation specializations that fill in
// its fields. mutation_handler.go and query_handler.go therefore join this
// file's namespace with //declscope:namespace handler.

// Package server provides HTTP server functionality for sql-http-proxy.
package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/executor"
	"github.com/mpyw/sql-http-proxy/internal/js"
	"github.com/mpyw/sql-http-proxy/internal/mock"
	"github.com/mpyw/sql-http-proxy/internal/server/body"
)

// handledQueryRecord represents an executed query record.
type handledQueryRecord struct {
	Type   config.OpType  // OpTypeOne, OpTypeMany, or OpTypeNone
	IsMock bool           // true if mock query
	SQL    string         // executed SQL (bound for DB, original for mock)
	Params map[string]any // named parameters (for mock)
	Args   []any          // positional arguments (for DB)
}

// handledQueryRecorder is called when a query is executed.
type handledQueryRecorder func(record handledQueryRecord)

// handlerOptions contains optional settings shared by the handler
// constructors. Shared on purpose: server.go (NewServeMux) fills it once and
// hands it to every newQueryHandlerWithOptions / newMutationHandlerWithOptions
// call.
//
//declscope:package
type handlerOptions struct {
	configDir   string              // Directory of config file for resolving relative paths
	helpers     *js.CompiledHelpers // Global JavaScript helpers
	valueParser *mock.ValueParser   // Global CSV value parser
	//declscope:private
	recorder handledQueryRecorder
}

// createNotFoundHandler creates a 404 handler.
// Shared on purpose: server.go (NewServeMux) installs it as the router's
// NotFound handler.
//
//declscope:package
func createNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := &responder{w: w}
		res.Error(http.StatusNotFound, errors.New("not found"))
	})
}

// HandlerExecutor is the handler-side interface for query/mutation executors.
type HandlerExecutor[R any] interface {
	Execute(ctx context.Context, params map[string]any, opts executor.Options) (*R, error)
}

// HandlerResultProcessor handles result-specific response building.
type HandlerResultProcessor[R any] interface {
	// BuildResponse returns (status, output, responseHeader) from the result.
	BuildResponse(result *R) (status int, output any, headers http.Header)
}

// baseHandler contains common handler logic.
type baseHandler[R any] struct {
	exec          HandlerExecutor[R]
	method        string
	pathParams    []string // path parameter names in order
	parser        *body.Parser
	recorder      handledQueryRecorder
	processor     HandlerResultProcessor[R]
	checkNotFound bool          // true for query handlers
	delay         time.Duration // artificial delay before response
}

// ServeHTTP implements http.Handler.
func (h *baseHandler[R]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// Apply artificial delay if configured
	if h.delay > 0 {
		time.Sleep(h.delay)
	}

	status, output, headers := h.processor.BuildResponse(result)
	h.applyResponseHeaders(w, headers)
	res.Respond(status, output)
}

func (h *baseHandler[R]) parseParams(r *http.Request) (map[string]any, error) {
	var params map[string]any
	var err error

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		params = h.extractNamedParams(r.URL.Query())
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

func (h *baseHandler[R]) handleError(res *responder, err error) {
	if h.checkNotFound && errors.Is(err, executor.ErrNotFound) {
		res.Error(http.StatusNotFound, errors.New("not found"))
		return
	}
	if transformErr, ok := errors.AsType[*js.TransformError](err); ok {
		status := lo.CoalesceOrEmpty(transformErr.Status, h.defaultStatusForPhase(err))
		res.Respond(status, transformErr.Body)
		return
	}
	res.Error(http.StatusInternalServerError, err)
}

// buildExecOptions creates executor.Options with optional recorder wrapping.
func (h *baseHandler[R]) buildExecOptions(r *http.Request) executor.Options {
	opts := executor.Options{
		HTTPRequest: r,
	}
	if h.recorder != nil {
		opts.Recorder = func(rec executor.Record) {
			h.recorder(handledQueryRecord{
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
func (h *baseHandler[R]) applyResponseHeaders(w http.ResponseWriter, headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}

// defaultStatusForPhase picks the handler's fallback HTTP status for a
// transform error that names no status of its own.
func (h *baseHandler[R]) defaultStatusForPhase(err error) int {
	if phaseErr, ok := errors.AsType[*executor.PhaseError](err); ok && phaseErr.Phase == "pre" {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

// extractNamedParams flattens query-string values into handler params.
func (h *baseHandler[R]) extractNamedParams(values url.Values) map[string]any {
	params := make(map[string]any, len(values))
	for key, vals := range values {
		if len(vals) > 0 {
			params[key] = vals[0]
		}
	}
	return params
}

// handlerPathParams extracts path parameter names from a chi path pattern.
// Example: /users/{id}/posts/{post_id} -> ["id", "post_id"]
// Example: /users/{id:[0-9]+} -> ["id"]
func handlerPathParams(path string) []string {
	var params []string
	i := 0
	for i < len(path) {
		if path[i] == '{' {
			i++
			start := i
			// Read parameter name until colon (regex) or closing brace
			for i < len(path) && path[i] != '}' && path[i] != ':' {
				i++
			}
			if i > start {
				params = append(params, path[start:i])
			}
			// Skip to closing brace
			for i < len(path) && path[i] != '}' {
				i++
			}
			if i < len(path) {
				i++ // skip closing brace
			}
		} else {
			i++
		}
	}
	return params
}
