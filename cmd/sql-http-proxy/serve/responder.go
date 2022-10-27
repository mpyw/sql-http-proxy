package serve

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type responder struct {
	w http.ResponseWriter
	r *http.Request
}

func (r *responder) Respond(status int, payload any) {
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
	if _, err = r.w.Write([]byte(fmt.Sprintf("{\"error\":%s}", msg))); err != nil {
		log.Println(err)
		return
	}
}
