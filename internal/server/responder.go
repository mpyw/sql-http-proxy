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

func (r *responder) Respond(status int, payload any) {
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
	msg := lo.Must(json.Marshal(err.Error()))
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(status)
	if _, err = fmt.Fprintf(r.w, "{\"error\":%s}", msg); err != nil {
		slog.Error("Failed to write error response", "error", err)
		return
	}
}
