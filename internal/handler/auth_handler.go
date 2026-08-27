package handler

import (
	"fmt"
	"net/http"

	"personae-fasti/internal/dto"
	"personae-fasti/internal/mapper"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// LogInByToken handles GET /login (protected)
func (h *AuthHandler) LoginByToken(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	if req.Player == nil {
		return e.ErrToApiError(e.ErrNotFound)
	}

	loginInfo := mapper.FormLoginInfoResponse(req.Request.Header.Get("Authorization"), req.Player)
	return httputils.Response{Status: http.StatusOK, Body: loginInfo}
}

// LogInByUsername handles GET /login/{username}
func (h *AuthHandler) LoginByUsername(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	username := req.Request.PathValue("username")

	available, err := h.svc.GetLoginUsernameAvailability(req.Context, username)
	if err != nil {
		return e.ErrToApiError(err)
	}

	return httputils.Response{Status: http.StatusOK, Body: struct {
		Available     bool   `json:"available"`
		CheckUsername string `json:"checkUsername"`
	}{
		Available:     available,
		CheckUsername: username,
	}}
}

// POST /login
func (h *AuthHandler) Login(req httputils.RequestData[dto.LoginRequest]) httputils.Responder {
	if req.Body.LoginSource != "Web" {
		return e.NewApiError(http.StatusPreconditionFailed, "wrong login source")
	}

	player, err := h.svc.AuthenticatePlayerWeb(req.Context, req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	token, err := h.svc.EmitWebToken(req.Context, player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	loginInfo := mapper.FormLoginInfoResponse(fmt.Sprintf("Web %s", token), player)

	return httputils.Response{Status: http.StatusOK, Body: loginInfo}
}

// SignUp handles POST /signup
func (h *AuthHandler) SignUp(req httputils.RequestData[dto.SignUpRequest]) httputils.Responder {
	player, err := h.svc.CreatePlayer(req.Context, req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	token, err := h.svc.EmitWebToken(req.Context, player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	return httputils.Response{Status: http.StatusOK, Body: mapper.FormLoginInfoResponse(fmt.Sprintf("Web %s", token), player)}
}
