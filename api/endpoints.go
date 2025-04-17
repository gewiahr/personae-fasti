package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"personae-fasti/api/resp"
	"personae-fasti/data"
	"strings"
)

func (api *APIServer) SetHandlers(router *http.ServeMux) {

	router.HandleFunc("GET /login/{accesskey}", api.HTTPWrapper(api.handleLogin))
	router.HandleFunc("GET /records", api.HTTPWrapper(api.PlayerWrapper(api.handleGetRecords)))
	router.HandleFunc("POST /record", api.HTTPWrapper(api.PlayerWrapper(api.handlePostRecord)))

}

// func (api *APIServer) handleHome(w http.ResponseWriter, r *http.Request) *APIError {
// 	return api.Respond(r, w, http.StatusOK, nil)
// }

func (api *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) *APIError {
	accesskey := r.PathValue("accesskey")

	player, err := api.storage.GetPlayerByAccessKey(accesskey)
	if err != nil {
		if err == sql.ErrNoRows {
			return api.HandleError(fmt.Errorf("login failed: no user info for the passkey %d", strings.ToLower(accesskey))).WithCode(http.StatusUnauthorized)
		} else {
			return api.HandleError(err)
		}
	}

	loginInfo := resp.LoginInfo{
		AccessKey: player.AccessKey,
		Player: resp.PlayerInfo{
			ID:       player.ID,
			Username: player.Username,
		},
		CurrentGame: resp.GameInfo{
			ID:    player.CurrentGame.ID,
			Title: player.CurrentGame.Name,
		},
	}

	return api.Respond(r, w, http.StatusOK, loginInfo)
}

func (api *APIServer) handleGetRecords(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	records, err := api.storage.GetCurrentGameRecords(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	players, err := api.storage.GetCurrentGamePlayers(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	var playersInfo []resp.PlayerInfo
	for _, player := range players {
		playersInfo = append(playersInfo, resp.PlayerInfo{
			ID:       player.ID,
			Username: player.Username,
		})
	}

	gameRecords := resp.GameRecords{
		Records: records,
		Players: playersInfo,
		CurrentGame: resp.GameInfo{
			ID:    p.CurrentGame.ID,
			Title: p.CurrentGame.Name,
		},
	}

	return api.Respond(r, w, http.StatusOK, gameRecords)
}

func (api *APIServer) handlePostRecord(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	var record data.Record
	err := ReadJsonBody(r, &record)
	if err != nil {
		return api.HandleError(err)
	}

	err = api.storage.InsertNewRecord(&record, p, p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	records, err := api.storage.GetCurrentGameRecords(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	players, err := api.storage.GetCurrentGamePlayers(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	var playersInfo []resp.PlayerInfo
	for _, player := range players {
		playersInfo = append(playersInfo, resp.PlayerInfo{
			ID:       player.ID,
			Username: player.Username,
		})
	}

	gameRecords := resp.GameRecords{
		Records: records,
		Players: playersInfo,
		CurrentGame: resp.GameInfo{
			ID:    p.CurrentGame.ID,
			Title: p.CurrentGame.Name,
		},
	}

	return api.Respond(r, w, http.StatusOK, gameRecords)
}
