package server

import (
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/samber/lo"
)

// marshalOptions configures response encoding.
//
// Deterministic sorts map keys, matching encoding/json v1 (which always sorted
// them) so that responses built from map[string]any rows stay byte-stable.
//
// AllowInvalidUTF8 keeps v1's lenient handling of column values that are not
// valid UTF-8: the invalid bytes are emitted as U+FFFD instead of failing the
// whole response. v2 would reject them by default, which would turn rows that
// used to serve fine into a 500.
var marshalOptions = json.JoinOptions(
	json.Deterministic(true),
	jsontext.AllowInvalidUTF8(true),
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

	// Marshal before writing the header so a marshal failure can still be
	// reported as a 500 instead of truncating an already-committed response.
	msg, err := json.Marshal(payload, marshalOptions)
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
	// Log server errors (5xx)
	if status >= 500 {
		slog.Error("Internal server error", "status", status, "error", err)
	}
	// Marshaling a string with AllowInvalidUTF8 cannot fail, so lo.Must is safe
	// even for error messages carrying non-UTF-8 bytes from a driver.
	msg := lo.Must(json.Marshal(err.Error(), marshalOptions))
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(status)
	if _, err = fmt.Fprintf(r.w, "{\"error\":%s}", msg); err != nil {
		slog.Error("Failed to write error response", "error", err)
		return
	}
}
