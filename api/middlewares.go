package api

import (
	"net/http"
	"personae-fasti/internal/dto"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"
)

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

func ImageAdapt(authService *service.AuthService, maxRequestBytes int64, fn HandlerFunc[dto.NoBody]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractBearer(r)
		player, err := authService.AuthenticateToken(r.Context(), token)
		if err != nil {
			httputils.RespondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
		req := httputils.RequestData[dto.NoBody]{
			Context: r.Context(),
			Player:  player,
			Query:   r.URL.Query(),
			Request: r,
		}
		resp := fn(req)
		if err := resp.Respond(w); err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		}
	}
}
