package handler

import (
	"net/http"
	"personae-fasti/internal/dto"
	"personae-fasti/internal/mapper"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"
)

type PlayerHandler struct {
	svc *service.PlayerService
}

func NewPlayerHandler(svc *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

// GetPlayerCurrentGame handles GET /player/currentGame (protected).
func (h *PlayerHandler) GetPlayerCurrentGame(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	currentGame, err := h.svc.GetPlayerCurrentGame(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	currentGameFull := mapper.GameToGameFull(currentGame)
	return httputils.Response{Status: http.StatusOK, Body: currentGameFull}
}

// ChangePlayerCurrentGame handles PUT /player/currentGame/{gameExt} (protected).
func (h *PlayerHandler) ChangePlayerCurrentGame(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	gameExt := req.Request.PathValue("gameExt")
	if gameExt == "" {
		return e.NewApiError(http.StatusBadRequest, "gameID cannot be empty")
	}

	if gameExt == req.Player.CurrentGame.ExtID {
		return e.NewApiError(http.StatusBadRequest, "error changing current game: this game is already current")
	}

	currentGame, err := h.svc.ChangePlayerCurrentGame(req.Context, req.Player, gameExt)
	if err != nil {
		return e.ErrToApiError(err)
	}

	currentGameFull := mapper.GameToGameFull(currentGame)
	return httputils.Response{Status: http.StatusOK, Body: currentGameFull}
}

// GetPlayerSettings handles GET /player/settings (protected).
func (h *PlayerHandler) GetPlayerSettings(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	player, err := h.svc.GetPlayerSettings(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	var currentGame *dto.GameFull
	if player.CurrentGame != nil {
		currentGame = mapper.GameToGameFull(player.CurrentGame)
	}

	var invites []dto.GameInvite
	if len(player.Invites) > 0 {
		playerInvites, err := h.svc.GetPlayerInvites(req.Context, player)
		if err != nil {
			return e.ErrToApiError(err)
		}

		invites = mapper.InviteToGameInviteArray(playerInvites)
	}

	return httputils.Response{Status: http.StatusOK, Body: dto.PlayerSettingsResponse{
		CurrentGame:   currentGame,
		PlayerGames:   mapper.GameToGameBriefArray(player.Games),
		PlayerInvites: invites,
	}}
}

// GameInviteAccept handles POST /player/invite/accept/{inviteCode} (protected).
func (h *PlayerHandler) AcceptGameInvite(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	inviteCode := req.Request.PathValue("inviteCode")
	if inviteCode == "" {
		return e.NewApiError(http.StatusBadRequest, "incorrect invite code")
	}

	err := h.svc.AcceptGameInvite(req.Context, req.Player, inviteCode)
	if err != nil {
		return e.ErrToApiError(err)
	}

	return httputils.Response{Status: http.StatusOK, Body: nil}
}

// GameInviteRefuse handles POST /player/invite/refuse/{inviteCode} (protected).
func (h *PlayerHandler) RefuseGameInvite(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	inviteCode := req.Request.PathValue("inviteCode")
	if inviteCode == "" {
		return e.NewApiError(http.StatusBadRequest, "incorrect invite code")
	}

	err := h.svc.RefuseGameInvite(req.Context, req.Player, inviteCode)
	if err != nil {
		return e.ErrToApiError(err)
	}

	return httputils.Response{Status: http.StatusOK, Body: nil}
}

// // GET /player/username/checkAvailability/{username}
// func (api *APIServer) handleCheckUsernameAvailability(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
// 	username := r.PathValue("username")
// 	if username == "" {
// 		return api.HandleErrorString("username cannot be empty")
// 	}

// 	available, err := api.storage.CheckUsernameAvailability(p, username)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	return api.Respond(r, w, http.StatusOK, struct {
// 		Available     bool   `json:"available"`
// 		CheckUsername string `json:"checkUsername"`
// 	}{
// 		Available:     available,
// 		CheckUsername: username,
// 	})
// }

// // PATCH /player/username
// func (api *APIServer) handleChangePlayerUsername(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
// 	var usernameChange reqData.UsernameChange
// 	err := ReadJsonBody(r, &usernameChange)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	p, err = api.storage.ChangeUsername(p, usernameChange.NewUsername)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	return api.Respond(r, w, http.StatusOK, respData.LoginPlayerInfo{
// 		ID:       p.ID,
// 		Username: p.Username,
// 		Settings: &respData.LoginPlayerInfoSettings{
// 			CouldChangeUsername: false,
// 		},
// 	})
// }
