package api

import (
	"net/http"
)

func (s *APIServer) HTTPWrapper(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if APIErr := f(w, r); APIErr != nil {
			s.Respond(r, w, APIErr.Code, APIErr)
		}
	}
}
