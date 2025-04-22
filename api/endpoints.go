package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"personae-fasti/api/models/reqData"
	"personae-fasti/api/models/respData"
	"personae-fasti/data"
	"strings"
)

func (api *APIServer) SetHandlers(router *http.ServeMux) {

	router.HandleFunc("GET /login/{accesskey}", api.HTTPWrapper(api.handleLogin))

	router.HandleFunc("GET /records", api.HTTPWrapper(api.PlayerWrapper(api.handleGetRecords)))
	router.HandleFunc("POST /record", api.HTTPWrapper(api.PlayerWrapper(api.handlePostRecord)))

	router.HandleFunc("GET /chars", api.HTTPWrapper(api.PlayerWrapper(api.handleGetChars)))
	router.HandleFunc("GET /char/{id}", api.HTTPWrapper(api.PlayerWrapper(api.handleGetCharByID)))
	router.HandleFunc("POST /char", api.HTTPWrapper(api.PlayerWrapper(api.handleCreateChar)))
	router.HandleFunc("PUT /char", api.HTTPWrapper(api.PlayerWrapper(api.handleUpdateChar)))

	router.HandleFunc("GET /suggestions", api.HTTPWrapper(api.PlayerWrapper(api.handleGetSuggestions)))
}

// func (api *APIServer) handleHome(w http.ResponseWriter, r *http.Request) *APIError {
// 	return api.Respond(r, w, http.StatusOK, nil)
// }

// GET /login/{accesskey}
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

	loginInfo := respData.LoginInfo{
		AccessKey: player.AccessKey,
		Player: respData.PlayerInfo{
			ID:       player.ID,
			Username: player.Username,
		},
		CurrentGame: respData.GameInfo{
			ID:    player.CurrentGame.ID,
			Title: player.CurrentGame.Name,
		},
	}

	return api.Respond(r, w, http.StatusOK, loginInfo)
}

// GET /records
func (api *APIServer) handleGetRecords(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	records, err := api.storage.GetCurrentGameRecords(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	players, err := api.storage.GetCurrentGamePlayers(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	var playersInfo []respData.PlayerInfo
	for _, player := range players {
		playersInfo = append(playersInfo, respData.PlayerInfo{
			ID:       player.ID,
			Username: player.Username,
		})
	}

	gameRecords := respData.GameRecords{
		Records: records,
		Players: playersInfo,
		CurrentGame: respData.GameInfo{
			ID:    p.CurrentGame.ID,
			Title: p.CurrentGame.Name,
		},
	}

	return api.Respond(r, w, http.StatusOK, gameRecords)
}

// POST /record
func (api *APIServer) handlePostRecord(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	var recordInsert reqData.RecordInsert
	err := ReadJsonBody(r, &recordInsert)
	if err != nil {
		return api.HandleError(err)
	}

	err = api.storage.InsertNewRecord(&recordInsert, p)
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

	gameRecords := respData.GameRecords{
		Records:     records,
		Players:     respData.PlayersToPlayersInfoArray(players),
		CurrentGame: *respData.GameToGameInfo(p.CurrentGame),
	}

	return api.Respond(r, w, http.StatusCreated, gameRecords)
}

// GET /chars
func (api *APIServer) handleGetChars(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	chars, err := api.storage.GetCurrentGameChars(p.CurrentGame)
	if err != nil {
		api.HandleError(err)
	}

	players, err := api.storage.GetCurrentGamePlayers(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	gameChars := respData.GameChars{
		Chars:       respData.CharToCharInfoArray(chars),
		Players:     respData.PlayersToPlayersInfoArray(players),
		CurrentGame: *respData.GameToGameInfo(p.CurrentGame),
	}

	return api.Respond(r, w, http.StatusOK, gameChars)
}

// GET /char/{id}
func (api *APIServer) handleGetCharByID(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	charID := getPathValueInt(r, "id")
	if charID < 0 {
		return api.HandleError(fmt.Errorf("error parsing id: char id is invalid"))
	}

	char, err := api.storage.GetCharByID(charID)
	if err != nil {
		return api.HandleError(err)
	} else if char == nil {
		return api.HandleErrorString(fmt.Sprintf("no character with id %d", charID)).WithCode(http.StatusNotFound)
	} else if char.GameID != p.CurrentGameID {
		return api.HandleErrorString(fmt.Sprintf("char %d is not allowed to request for the game %d", char.ID, p.CurrentGameID)).WithCode(http.StatusForbidden)
	}
	// ++ Add char check ++//

	charPage := respData.CharPage{
		Char:    *respData.CharToCharFullInfo(char),
		Records: []data.Record{},
	}

	return api.Respond(r, w, http.StatusOK, charPage)
}

// POST /char
func (api *APIServer) handleCreateChar(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	var charCreate reqData.CharCreate
	err := ReadJsonBody(r, &charCreate)
	if err != nil {
		return api.HandleError(err)
	}

	char, err := api.storage.CreateChar(&charCreate, p)
	if err != nil {
		return api.HandleError(err)
	}

	charFullInfo := respData.CharToCharFullInfo(char)
	return api.Respond(r, w, http.StatusOK, charFullInfo)
}

// PUT /char
func (api *APIServer) handleUpdateChar(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	var charUpdate reqData.CharUpdate
	err := ReadJsonBody(r, &charUpdate)
	if err != nil {
		return api.HandleError(err)
	}

	char, err := api.storage.GetCharByID(charUpdate.ID)
	if err != nil {
		return api.HandleError(err)
	} else if char == nil {
		return api.HandleErrorString(fmt.Sprintf("no character with id %d", charUpdate.ID)).WithCode(http.StatusNotFound)
	} else if char.GameID != p.CurrentGameID {
		return api.HandleErrorString(fmt.Sprintf("char %d is not allowed to request for the game %d", char.ID, p.CurrentGameID)).WithCode(http.StatusForbidden)
	}
	// ++ Add char check ++//

	char, err = api.storage.UpdateChar(&charUpdate, char)
	if err != nil {
		api.HandleError(err)
	}

	charFullInfo := respData.CharToCharFullInfo(char)
	return api.Respond(r, w, http.StatusOK, charFullInfo)
}

// GET /suggestions
func (api *APIServer) handleGetSuggestions(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	suggestions, err := api.storage.GetSuggestions(p)
	if err != nil {
		api.HandleError(err)
	}

	return api.Respond(r, w, http.StatusOK, respData.SuggestionData{Suggestions: suggestions})
}
