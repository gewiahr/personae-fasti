package api

import (
	"context"
	"net/http"
	"personae-fasti/internal/dto"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"
	"time"
)

func Adapt[Body any](fn HandlerFunc[Body]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, apiErr := readBody[Body](w, r)
		if apiErr != nil {
			respondAndTrack(w, r, *apiErr)
			return
		}

		req := httputils.RequestData[Body]{
			Context: r.Context(),
			Player:  nil,
			Body:    body,
			Query:   r.URL.Query(),
			Request: r,
		}

		respondAndTrack(w, r, fn(req))
	}
}

func AuthAdapt[Body any](authService *service.AuthService, fn HandlerFunc[Body]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractBearer(r)
		player, err := authService.AuthenticateToken(r.Context(), token)
		if err != nil {
			respondAndTrack(w, r, e.ErrToApiError(err))
			return
		}
		setRequestLogPlayer(r.Context(), player.ID, player.CurrentGameID)

		body, apiErr := readBody[Body](w, r)
		if apiErr != nil {
			respondAndTrack(w, r, *apiErr)
			return
		}

		req := httputils.RequestData[Body]{
			Context: r.Context(),
			Player:  player,
			Body:    body,
			Query:   r.URL.Query(),
			Request: r,
		}

		respondAndTrack(w, r, fn(req))
	}
}

func ImageAdapt(authService *service.AuthService, maxRequestBytes int64, timeout time.Duration, fn HandlerFunc[dto.NoBody]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		r = r.WithContext(ctx)
		token := extractBearer(r)
		player, err := authService.AuthenticateToken(r.Context(), token)
		if err != nil {
			respondAndTrack(w, r, e.ErrToApiError(err))
			return
		}
		setRequestLogPlayer(r.Context(), player.ID, player.CurrentGameID)
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
		req := httputils.RequestData[dto.NoBody]{
			Context: r.Context(),
			Player:  player,
			Query:   r.URL.Query(),
			Request: r,
		}
		respondAndTrack(w, r, fn(req))
	}
}

func respondAndTrack(w http.ResponseWriter, r *http.Request, responder httputils.Responder) {
	switch apiErr := responder.(type) {
	case e.ApiError:
		apiErr.RequestID = requestID(r.Context())
		setRequestLogError(r.Context(), apiErr.Code, apiErr.Internal())
		responder = apiErr
	case *e.ApiError:
		apiErr.RequestID = requestID(r.Context())
		setRequestLogError(r.Context(), apiErr.Code, apiErr.Internal())
	}
	if err := responder.Respond(w); err != nil {
		setRequestLogError(r.Context(), "response_write_failed", err)
		http.Error(w, `{"error":"Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
	}
}
