package api

import (
	"net/http"
)

func (api *APIServer) SetHandlers(router *http.ServeMux) {

	router.HandleFunc("GET /", api.HTTPWrapper(api.handleHome))

}

func (api *APIServer) handleHome(w http.ResponseWriter, r *http.Request) *APIError {
	return api.Respond(r, w, http.StatusOK, nil)
}
