package handler

import (
	"net/http"

	"personae-fasti/internal/dto"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"
)

type AppHandler struct {
	svc *service.AppService
}

func NewAppHandler(svc *service.AppService) *AppHandler {
	return &AppHandler{svc: svc}
}

// PostFeedback handles POST /feedback (protected)
func (h *AppHandler) LoginByToken(req httputils.RequestData[dto.ServiceFeedback]) httputils.Responder {
	_, err := h.svc.PostFeedback(req.Context, req.Player.ID, req.Player.CurrentGameID, req.Body.Type, req.Body.Text)
	if err != nil {
		return e.ErrToApiError(err)
	}

	return httputils.Response{Status: http.StatusCreated, Body: nil}
}
