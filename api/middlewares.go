package api

import (
	"net/http"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"
)

// func (api *APIServer) HandlerWrapper(f HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if APIErr := f(w, r); APIErr != nil {
// 			api.Respond(r, w, APIErr.Code, APIErr)
// 		}
// 	}
// }

// func (api *APIServer) PlayerWrapper(f AuthHandlerFunc) HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) *APIError {
// 		var player *domain.Player
// 		var err error

// 		//defer api.storage.Log()

// 		//accesskey := r.Header.Get("AccessKey")
// 		token := r.Header.Get("Authorization")
// 		if token == "" {
// 			return api.HandleErrorString("authorization is invalid").WithCode(http.StatusUnauthorized)
// 		}

// 		tokenArray := strings.Split(token, " ")
// 		if len(tokenArray) == 1 {
// 			player, err = api.storage.GetPlayerByTGToken(tokenArray[0])
// 		} else if len(tokenArray) == 2 {
// 			player, err = api.storage.GetPlayerByTGToken(tokenArray[1])
// 		} else {
// 			return api.HandleError(err)
// 		}

// 		if err == sql.ErrNoRows {
// 			return api.HandleErrorString("token is invalid").WithCode(http.StatusUnauthorized)
// 		}

// 		// player, err := api.storage.GetPlayerByAccessKey(accesskey)
// 		// if err != nil {
// 		// 	if err == sql.ErrNoRows {
// 		// 		return api.HandleError(fmt.Errorf("login failed: no user info for the passkey %s", accesskey)).WithCode(http.StatusUnauthorized)
// 		// 	} else {
// 		// 		return api.HandleError(err)
// 		// 	}
// 		// }

// 		if APIErr := f(w, r, player); APIErr != nil {
// 			return api.Respond(r, w, APIErr.Code, APIErr)
// 		}

// 		return nil
// 	}
// }

// HandlerWithPlayer is a pure handler that receives the authenticated player directly.
//type HandlerWithPlayer[Body any] func(player *domain.Player, req HandlerRequest[Body]) Response

func Adapt[Body any](fn HandlerFunc[Body]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, apiErr := readBody[Body](w, r)
		if apiErr != nil {
			apiErr.Respond(w)
			return
		}

		req := httputils.RequestData[Body]{
			Context: r.Context(),
			Player:  nil,
			Body:    body,
			Query:   r.URL.Query(),
			Request: r,
		}

		resp := fn(req)
		if err := resp.Respond(w); err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		}
	}
}

func AuthAdapt[Body any](authService *service.AuthService, fn HandlerFunc[Body]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractBearer(r)
		player, err := authService.AuthenticateToken(r.Context(), token)
		if err != nil {
			httputils.RespondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		body, apiErr := readBody[Body](w, r)
		if apiErr != nil {
			apiErr.Respond(w)
			return
		}

		req := httputils.RequestData[Body]{
			Context: r.Context(),
			Player:  player,
			Body:    body,
			Query:   r.URL.Query(),
			Request: r,
		}

		resp := fn(req)
		if err := resp.Respond(w); err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		}
	}
}
