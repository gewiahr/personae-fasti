package httputils

import (
	"encoding/json"
	"net/http"
	e "personae-fasti/internal/pkg/errorutils"
)

type Responder interface {
	Respond(w http.ResponseWriter) error
}

type Response struct {
	Status int
	Body   any
}

func (r Response) Respond(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	return json.NewEncoder(w).Encode(r.Body)
}

func RespondError(w http.ResponseWriter, status int, msg string) {
	e.NewApiError(status, msg).Respond(w)
}
