// Package server provides HTTP server functionality for sql-http-proxy.
package server

import (
	"errors"
	"net/http"

	"github.com/mpyw/sql-http-proxy/internal/js"
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
	Recorder QueryRecorder
}

// handleTransformError handles errors from transform, using status from TransformError if available.
func handleTransformError(res *responder, err error, defaultStatus int) {
	var transformErr *js.TransformError
	if errors.As(err, &transformErr) {
		status := transformErr.Status
		if status == 0 {
			status = defaultStatus
		}
		res.RespondBody(status, transformErr.Body)
		return
	}
	res.Error(defaultStatus, err)
}

// CreateNotFoundHandler creates a 404 handler.
func CreateNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := &responder{w: w}
		res.Error(http.StatusNotFound, errors.New("not found"))
	})
}
