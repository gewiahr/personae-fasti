package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

func (api *APIServer) HTTPWrapper(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if APIErr := f(w, r); APIErr != nil {
			api.Respond(r, w, APIErr.Code, APIErr)
		}
	}
}

func (api *APIServer) PlayerWrapper(f APIFuncAuth) APIFunc {
	return func(w http.ResponseWriter, r *http.Request) *APIError {
		accesskey := r.Header.Get("AccessKey")
		player, err := api.storage.GetPlayerByAccessKey(accesskey)
		if err != nil {
			if err == sql.ErrNoRows {
				return api.HandleError(fmt.Errorf("login failed: no user info for the passkey %d", strings.ToLower(accesskey))).WithCode(http.StatusUnauthorized)
			} else {
				return api.HandleError(err)
			}
		}

		if APIErr := f(w, r, player); APIErr != nil {
			return api.Respond(r, w, APIErr.Code, APIErr)
		}

		return nil
	}
}
