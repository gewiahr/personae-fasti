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

// StartNewGameSession handles POST /game/session/new (protected).
func (h *GameHandler) StartNewGameSession(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	session, err := h.svc.StartNewGameSession(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	currentGameSessionBrief := mapper.SessionToSessionBrief(session)
	return httputils.Response{Status: http.StatusCreated, Body: currentGameSessionBrief}
}

// EditGameSession handles PATCH /game/session (protected).
func (h *GameHandler) EditGameSession(req httputils.RequestData[dto.SessionUpdate]) httputils.Responder {
	session, err := h.svc.EditGameSession(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	currentGameSessionBrief := mapper.SessionToSessionBrief(session)
	return httputils.Response{Status: http.StatusCreated, Body: currentGameSessionBrief}
}

// RemoveLastGameSession handles DELETE /game/session/remove (protected).
func (h *GameHandler) RemoveLastGameSession(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	err := h.svc.RemoveLastGameSession(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	return httputils.Response{Status: http.StatusCreated, Body: nil}
}

// InvitePlayer handles POST /game/invite/{username} (protected).
func (h *GameHandler) InvitePlayer(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	username := req.Request.PathValue("username")
	if username == "" {
		return e.NewApiError(http.StatusBadRequest, "error parsing username: username is invalid")
	}

	err := h.svc.InvitePlayer(req.Context, req.Player, username)
	if err != nil {
		return e.ErrToApiError(err)
	}

	return httputils.Response{Status: http.StatusCreated, Body: nil}
}

// RemovePlayerInvite handles DELETE /game/invite/{username} (protected).
func (h *GameHandler) RemovePlayerInvite(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	username := req.Request.PathValue("username")
	if username == "" {
		return e.NewApiError(http.StatusBadRequest, "error parsing username: username is invalid")
	}

	err := h.svc.RemoveInvite(req.Context, req.Player, username)
	if err != nil {
		return e.ErrToApiError(err)
	}

	return httputils.Response{Status: http.StatusOK, Body: nil}
}

// PutGameSettings handles PUT /game/settings (protected).
func (h *GameHandler) PutGameSettings(req httputils.RequestData[dto.GameSettingsUpdate]) httputils.Responder {
	currentGame, err := h.svc.UpdateGameSettings(req.Context, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	currentGameInfo := mapper.GameToGameFull(currentGame)
	return httputils.Response{Status: http.StatusOK, Body: currentGameInfo}

	//api.Respond(r, w, http.StatusOK, currentGameInfo)
}

// // PUT /game/settings
// func (api *APIServer) handlePutGameSettings(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
// 	var gameSettingsUpdate reqData.GameSettingsUpdate
// 	err := ReadJsonBody(r, &gameSettingsUpdate)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	currentGameInfo := respData.GameToGameFullInfo(currentGame)
// 	return api.Respond(r, w, http.StatusOK, currentGameInfo)
// }
