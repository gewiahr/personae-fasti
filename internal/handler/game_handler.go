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

// GetGameByExt handles GET /game/{ext} (protected).
func (h *GameHandler) GetPlayerGameByExt(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	gameExt := req.Request.PathValue("ext")
	if gameExt == "" {
		return e.NewApiError(http.StatusBadRequest, "error parsing id: game ext is invalid")
	}

	game, err := h.svc.GetPlayerGameByExt(req.Context, req.Player.ID, gameExt) //api.storage.GetPlayerGame(gameID, p.ID)
	if err != nil {
		return e.ErrToApiError(err)
	}

	currentGameFull := mapper.GameToGameFull(game)
	return httputils.Response{Status: http.StatusOK, Body: currentGameFull}
}

// CreateNewGame handles POST /game (protected).
func (h *GameHandler) CreateNewGame(req httputils.RequestData[dto.GameCreate]) httputils.Responder {
	game, err := h.svc.CreateGame(req.Context, req.Player, req.Body.Title)
	if err != nil {
		return e.ErrToApiError(err)
	}

	currentGameFull := mapper.GameToGameFull(game)
	return httputils.Response{Status: http.StatusCreated, Body: currentGameFull}
}

// EditGame handles PUT /game (protected).
func (h *GameHandler) EditGame(req httputils.RequestData[dto.GameUpdate]) httputils.Responder {
	game, err := h.svc.EditGame(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	currentGameFull := mapper.GameToGameFull(game)
	return httputils.Response{Status: http.StatusCreated, Body: currentGameFull}
}
