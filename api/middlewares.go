package api

import (
	"database/sql"
	"net/http"
	"personae-fasti/data"
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
		var player *data.Player
		var err error
		//accesskey := r.Header.Get("AccessKey")
		token := r.Header.Get("Authorization")

		if token == "" {
			return api.HandleErrorString("authorization is invalid").WithCode(http.StatusUnauthorized)
		}

		tokenArray := strings.Split(token, " ")
		if len(tokenArray) == 1 {
			player, err = api.storage.GetPlayerByTGToken(tokenArray[0])
		} else if len(tokenArray) == 2 {
			player, err = api.storage.GetPlayerByTGToken(tokenArray[1])
		} else {
			return api.HandleError(err)
		}

		if err == sql.ErrNoRows {
			return api.HandleErrorString("token is invalid").WithCode(http.StatusUnauthorized)
		}

		// player, err := api.storage.GetPlayerByAccessKey(accesskey)
		// if err != nil {
		// 	if err == sql.ErrNoRows {
		// 		return api.HandleError(fmt.Errorf("login failed: no user info for the passkey %s", accesskey)).WithCode(http.StatusUnauthorized)
		// 	} else {
		// 		return api.HandleError(err)
		// 	}
		// }

		if APIErr := f(w, r, player); APIErr != nil {
			return api.Respond(r, w, APIErr.Code, APIErr)
		}

		return nil
	}
}
