package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
		log.Println(err)
		return
	}
}

func (r *responder) Error(status int, err error) {
	msg, _ := json.Marshal(err.Error())
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(status)
	if _, err = fmt.Fprintf(r.w, "{\"error\":%s}", msg); err != nil {
		log.Println(err)
		return
	}
}

// RespondBody writes the body as JSON response.
// All types are JSON serialized (string -> "...", object -> {...}, etc.)
func (r *responder) RespondBody(status int, body any) {
	r.Respond(status, body)
}
