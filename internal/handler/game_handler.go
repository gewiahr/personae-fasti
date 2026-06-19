package handler

import (
	"net/http"
	"personae-fasti/internal/dto"
	"personae-fasti/internal/mapper"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"
)

type GameHandler struct {
	svc *service.GameService
}

func NewGameHandler(svc *service.GameService) *GameHandler {
	return &GameHandler{svc: svc}
}

// GetPlayerCurrentGame handles GET /player/currentGame (protected).
func (h *GameHandler) GetPlayerCurrentGame(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	currentGame, err := h.svc.GetPlayerCurrentGame(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	currentGameFull := mapper.GameToGameFull(currentGame)
	return httputils.Response{Status: http.StatusOK, Body: currentGameFull} //api.Respond(r, w, http.StatusOK, currentGameInfo)
}
