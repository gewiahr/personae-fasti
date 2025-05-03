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
	router.HandleFunc("PUT /record", api.HTTPWrapper(api.PlayerWrapper(api.handleChangeRecord)))

	router.HandleFunc("GET /chars", api.HTTPWrapper(api.PlayerWrapper(api.handleGetChars)))
	router.HandleFunc("GET /char/{id}", api.HTTPWrapper(api.PlayerWrapper(api.handleGetCharByID)))
	router.HandleFunc("POST /char", api.HTTPWrapper(api.PlayerWrapper(api.handleCreateChar)))
	router.HandleFunc("PUT /char", api.HTTPWrapper(api.PlayerWrapper(api.handleUpdateChar)))

	router.HandleFunc("GET /npcs", api.HTTPWrapper(api.PlayerWrapper(api.handleGetNPCs)))
	router.HandleFunc("GET /npc/{id}", api.HTTPWrapper(api.PlayerWrapper(api.handleGetNPCByID)))
	router.HandleFunc("POST /npc", api.HTTPWrapper(api.PlayerWrapper(api.handleCreateNPC)))
	router.HandleFunc("PUT /npc", api.HTTPWrapper(api.PlayerWrapper(api.handleUpdateNPC)))

	router.HandleFunc("GET /locations", api.HTTPWrapper(api.PlayerWrapper(api.handleGetLocations)))
	router.HandleFunc("GET /location/{id}", api.HTTPWrapper(api.PlayerWrapper(api.handleGetLocationByID)))
	router.HandleFunc("POST /location", api.HTTPWrapper(api.PlayerWrapper(api.handleCreateLocation)))
	router.HandleFunc("PUT /location", api.HTTPWrapper(api.PlayerWrapper(api.handleUpdateLocation)))

	router.HandleFunc("GET /suggestions", api.HTTPWrapper(api.PlayerWrapper(api.handleGetSuggestions)))

	router.HandleFunc("GET /player/settings", api.HTTPWrapper(api.PlayerWrapper(api.handleGetPlayerSetings)))
	router.HandleFunc("PUT /player/game", api.HTTPWrapper(api.PlayerWrapper(api.handleChangePlayerGame)))
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

	sessions, err := api.storage.GetCurrentGameSessions(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	gameRecords := respData.FormGameRecords(p, records, players, sessions)

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

	sessions, err := api.storage.GetCurrentGameSessions(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	gameRecords := respData.FormGameRecords(p, records, players, sessions)

	return api.Respond(r, w, http.StatusCreated, gameRecords)
}

// PUT /record
func (api *APIServer) handleChangeRecord(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	var recordUpdate reqData.RecordUpdate
	err := ReadJsonBody(r, &recordUpdate)
	if err != nil {
		return api.HandleError(err)
	}

	err = api.storage.UpdateRecord(&recordUpdate, p)
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

	sessions, err := api.storage.GetCurrentGameSessions(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	gameRecords := respData.FormGameRecords(p, records, players, sessions)

	return api.Respond(r, w, http.StatusOK, gameRecords)
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
		Records: char.Records, // ** change to mention API type ** //
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
	return api.Respond(r, w, http.StatusCreated, charFullInfo)
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

// GET /npcs
func (api *APIServer) handleGetNPCs(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	npcs, err := api.storage.GetCurrentGameNPCs(p.CurrentGame)
	if err != nil {
		api.HandleError(err)
	}

	gameNPCs := respData.GameNPCs{
		NPCs:        respData.NPCToNPCInfoArray(npcs),
		CurrentGame: *respData.GameToGameInfo(p.CurrentGame),
	}

	return api.Respond(r, w, http.StatusOK, gameNPCs)
}

// GET /npc/{id}
func (api *APIServer) handleGetNPCByID(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	npcID := getPathValueInt(r, "id")
	if npcID < 0 {
		return api.HandleError(fmt.Errorf("error parsing id: npc id is invalid"))
	}

	npc, err := api.storage.GetNPCByID(npcID)
	if err != nil {
		return api.HandleError(err)
	} else if npc == nil {
		return api.HandleErrorString(fmt.Sprintf("no npc with id %d", npcID)).WithCode(http.StatusNotFound)
	} else if npc.GameID != p.CurrentGameID {
		return api.HandleErrorString(fmt.Sprintf("npc %d is not allowed to request for the game %d", npc.ID, p.CurrentGameID)).WithCode(http.StatusForbidden)
	}
	// ++ Add char check ++//

	npcPage := respData.NPCPage{
		NPC:     *respData.NPCToNPCFullInfo(npc),
		Records: npc.Records, // ** change to mention API type ** //
	}

	return api.Respond(r, w, http.StatusOK, npcPage)
}

// POST /npc
func (api *APIServer) handleCreateNPC(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	var npcCreate reqData.NPCCreate
	err := ReadJsonBody(r, &npcCreate)
	if err != nil {
		return api.HandleError(err)
	}

	npc, err := api.storage.CreateNPC(&npcCreate, p)
	if err != nil {
		return api.HandleError(err)
	}

	npcFullInfo := respData.NPCToNPCFullInfo(npc)
	return api.Respond(r, w, http.StatusCreated, npcFullInfo)
}

// PUT /npc
func (api *APIServer) handleUpdateNPC(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	var npcUpdate reqData.NPCUpdate
	err := ReadJsonBody(r, &npcUpdate)
	if err != nil {
		return api.HandleError(err)
	}

	npc, err := api.storage.GetNPCByID(npcUpdate.ID)
	if err != nil {
		return api.HandleError(err)
	} else if npc == nil {
		return api.HandleErrorString(fmt.Sprintf("no npc with id %d", npcUpdate.ID)).WithCode(http.StatusNotFound)
	} else if npc.GameID != p.CurrentGameID {
		return api.HandleErrorString(fmt.Sprintf("npc %d is not allowed to request for the game %d", npc.ID, p.CurrentGameID)).WithCode(http.StatusForbidden)
	}
	// ++ Add char check ++//

	npc, err = api.storage.UpdateNPC(&npcUpdate, npc)
	if err != nil {
		api.HandleError(err)
	}

	npcFullInfo := respData.NPCToNPCFullInfo(npc)
	return api.Respond(r, w, http.StatusOK, npcFullInfo)
}

// GET /locations
func (api *APIServer) handleGetLocations(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	locations, err := api.storage.GetCurrentGameLocations(p.CurrentGame)
	if err != nil {
		api.HandleError(err)
	}

	gameLocations := respData.GameLocations{
		Locations:   respData.LocationToLocationInfoArray(locations),
		CurrentGame: *respData.GameToGameInfo(p.CurrentGame),
	}

	return api.Respond(r, w, http.StatusOK, gameLocations)
}

// GET /location/{id}
func (api *APIServer) handleGetLocationByID(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	locationID := getPathValueInt(r, "id")
	if locationID < 0 {
		return api.HandleError(fmt.Errorf("error parsing id: location id is invalid"))
	}

	location, err := api.storage.GetLocationByID(locationID)
	if err != nil {
		return api.HandleError(err)
	} else if location == nil {
		return api.HandleErrorString(fmt.Sprintf("no location with id %d", locationID)).WithCode(http.StatusNotFound)
	} else if location.GameID != p.CurrentGameID {
		return api.HandleErrorString(fmt.Sprintf("location %d is not allowed to request for the game %d", location.ID, p.CurrentGameID)).WithCode(http.StatusForbidden)
	}
	// ++ Add char check ++//

	locationPage := respData.LocationPage{
		Location: *respData.LocationToLocationFullInfo(location),
		Records:  location.Records, // ** change to mention API type ** //
	}

	return api.Respond(r, w, http.StatusOK, locationPage)
}

// POST /location
func (api *APIServer) handleCreateLocation(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	var locationCreate reqData.LocationCreate
	err := ReadJsonBody(r, &locationCreate)
	if err != nil {
		return api.HandleError(err)
	}

	location, err := api.storage.CreateLocation(&locationCreate, p)
	if err != nil {
		return api.HandleError(err)
	}

	locationFullInfo := respData.LocationToLocationFullInfo(location)
	return api.Respond(r, w, http.StatusCreated, locationFullInfo)
}

// PUT /location
func (api *APIServer) handleUpdateLocation(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	var locationUpdate reqData.LocationUpdate
	err := ReadJsonBody(r, &locationUpdate)
	if err != nil {
		return api.HandleError(err)
	}

	location, err := api.storage.GetLocationByID(locationUpdate.ID)
	if err != nil {
		return api.HandleError(err)
	} else if location == nil {
		return api.HandleErrorString(fmt.Sprintf("no location with id %d", locationUpdate.ID)).WithCode(http.StatusNotFound)
	} else if location.GameID != p.CurrentGameID {
		return api.HandleErrorString(fmt.Sprintf("location %d is not allowed to request for the game %d", location.ID, p.CurrentGameID)).WithCode(http.StatusForbidden)
	}
	// ++ Add char check ++//

	location, err = api.storage.UpdateLocation(&locationUpdate, location)
	if err != nil {
		api.HandleError(err)
	}

	locationFullInfo := respData.LocationToLocationFullInfo(location)
	return api.Respond(r, w, http.StatusOK, locationFullInfo)
}

// GET /suggestions
func (api *APIServer) handleGetSuggestions(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	suggestions, err := api.storage.GetSuggestions(p)
	if err != nil {
		api.HandleError(err)
	}

	return api.Respond(r, w, http.StatusOK, respData.SuggestionData{Suggestions: suggestions})
}

// GET /player/settings
func (api *APIServer) handleGetPlayerSetings(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	playerGames, err := api.storage.GetPlayerGames(p)
	if err != nil {
		api.HandleError(err)
	}

	playerSettings := respData.FormPlayerSettings(playerGames, *p.CurrentGame)
	return api.Respond(r, w, http.StatusOK, playerSettings)
}

// PUT /player/game
func (api *APIServer) handleChangePlayerGame(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	var currentGameChange reqData.GameChange
	err := ReadJsonBody(r, &currentGameChange)
	if err != nil {
		return api.HandleError(err)
	}

	currentGame, err := api.storage.ChangeCurrentGame(p, currentGameChange.GameID)
	if err != nil {
		return api.HandleError(err)
	}

	currentGameInfo := respData.GameToGameInfo(currentGame)
	return api.Respond(r, w, http.StatusOK, currentGameInfo)
}
