package server

import (
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/samber/lo"
)

// marshalOptions pins encoding/json/v2 to the response bytes v1 produced.
// Responses are this project's public API, so v2's stricter or tidier defaults
// are opted out of one by one rather than accepted wholesale:
//
//   - Deterministic sorts map keys. v1 always sorted them, and rows arrive as
//     map[string]any, so without this the column order in a response would
//     vary between identical requests.
//   - AllowInvalidUTF8 keeps v1's lenient handling of column values that are
//     not valid UTF-8: invalid bytes are emitted as U+FFFD instead of failing.
//     v2 rejects them by default, which would turn rows that serve fine today
//     (latin1 or binary columns) into a 500.
//   - FormatNilSliceAsNull and FormatNilMapAsNull keep v1's `null` for a nil
//     slice or map. v2 renders them as `[]` and `{}`.
//
// The one difference left is deliberate: v2 does not escape <, > and & , and
// responses are served as application/json where v1's escaping bought nothing.
var marshalOptions = json.JoinOptions(
	json.Deterministic(true),
	jsontext.AllowInvalidUTF8(true),
	json.FormatNilSliceAsNull(true),
	json.FormatNilMapAsNull(true),
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
