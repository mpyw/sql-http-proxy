package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/samber/lo"
)

type responder struct {
	w http.ResponseWriter
}

// validateStatus ensures status code is in valid HTTP range (100-599).
// Returns the validated status or 500 if invalid.
func validateStatus(status int) int {
	if status < 100 || status > 599 {
		slog.Warn("Invalid HTTP status code, falling back to 500", "status", status)
		return http.StatusInternalServerError
	}
	return status
}

func (r *responder) Respond(status int, payload any) {
	status = validateStatus(status)
	// For 204 No Content, don't write any body
	if status == http.StatusNoContent {
		r.w.WriteHeader(status)
		return
	}

	msg, err := json.Marshal(payload)
	if err != nil {
		r.Error(http.StatusInternalServerError, err)
		return
	}
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(status)
	if _, err = r.w.Write(msg); err != nil {
		slog.Error("Failed to write response", "error", err)
		return
	}
}

func (r *responder) Error(status int, err error) {
	status = validateStatus(status)
	msg := lo.Must(json.Marshal(err.Error()))
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(status)
	if _, err = fmt.Fprintf(r.w, "{\"error\":%s}", msg); err != nil {
		slog.Error("Failed to write error response", "error", err)
		return
	}
}
