package api

import (
	"database/sql"
	"fmt"
	"io"
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
	router.HandleFunc("DELETE /record/{id}", api.HTTPWrapper(api.PlayerWrapper(api.handleDeleteRecord)))

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

	router.HandleFunc("POST /game/session/new", api.HTTPWrapper(api.PlayerWrapper(api.handleStartNewGameSession)))

	router.HandleFunc("GET /image/{type}/{id}", api.HTTPWrapper(api.handleGetImage))
	router.HandleFunc("POST /image/{type}/{id}", api.HTTPWrapper(api.handlePostImage))
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
			GMID:  player.CurrentGame.GMID,
		},
	}

	return api.Respond(r, w, http.StatusOK, loginInfo)
}

// GET /records
func (api *APIServer) handleGetRecords(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	records, err := api.storage.GetCurrentGameRecordsForPlayer(p.CurrentGame, p)
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

	records, err := api.storage.GetCurrentGameRecordsForPlayer(p.CurrentGame, p)
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

	records, err := api.storage.GetCurrentGameRecordsForPlayer(p.CurrentGame, p)
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

// DELETE /record
func (api *APIServer) handleDeleteRecord(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	recordID := getPathValueInt(r, "id")
	if recordID < 0 {
		return api.HandleError(fmt.Errorf("error parsing id: record id is invalid"))
	}

	err := api.storage.DeleteRecord(recordID, p)
	if err != nil {
		return api.HandleError(err)
	}

	return api.Respond(r, w, http.StatusOK, nil)
}

// GET /chars
func (api *APIServer) handleGetChars(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	chars, err := api.storage.GetCurrentGameChars(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
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

	records := []data.Record{}
	if len(char.Records) > 0 {
		records, err = api.storage.GetAllowedRecords(char.Records, p.ID)
	}

	charPage := respData.CharPage{
		Char:    *respData.CharToCharFullInfo(char),
		Records: records, // ** change to mention API type ** //
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

	char, err = api.storage.UpdateChar(&charUpdate, char, p)
	if err != nil {
		return api.HandleError(err)
	}

	charFullInfo := respData.CharToCharFullInfo(char)
	return api.Respond(r, w, http.StatusOK, charFullInfo)
}

// GET /npcs
func (api *APIServer) handleGetNPCs(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	npcs, err := api.storage.GetCurrentGameNPCs(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	npcs, err = api.storage.GetAllowedNPCs(npcs, p.ID)
	if err != nil {
		return api.HandleError(err)
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

	records := []data.Record{}
	if len(npc.Records) > 0 {
		records, err = api.storage.GetAllowedRecords(npc.Records, p.ID)
	}

	npcPage := respData.NPCPage{
		NPC:     *respData.NPCToNPCFullInfo(npc),
		Records: records, // ** change to mention API type ** //
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

	npc, err = api.storage.UpdateNPC(&npcUpdate, npc, p)
	if err != nil {
		return api.HandleError(err)
	}

	npcFullInfo := respData.NPCToNPCFullInfo(npc)
	return api.Respond(r, w, http.StatusOK, npcFullInfo)
}

// GET /locations
func (api *APIServer) handleGetLocations(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	locations, err := api.storage.GetCurrentGameLocations(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	locations, err = api.storage.GetAllowedLocations(locations, p.ID)
	if err != nil {
		return api.HandleError(err)
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

	records := []data.Record{}
	if len(location.Records) > 0 {
		records, err = api.storage.GetAllowedRecords(location.Records, p.ID)
	}

	locationPage := respData.LocationPage{
		Location: *respData.LocationToLocationFullInfo(location),
		Records:  records, // ** change to mention API type ** //
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

	location, err = api.storage.UpdateLocation(&locationUpdate, location, p)
	if err != nil {
		return api.HandleError(err)
	}

	locationFullInfo := respData.LocationToLocationFullInfo(location)
	return api.Respond(r, w, http.StatusOK, locationFullInfo)
}

// GET /suggestions
func (api *APIServer) handleGetSuggestions(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	suggestions, err := api.storage.GetSuggestions(p)
	if err != nil {
		return api.HandleError(err)
	}

	return api.Respond(r, w, http.StatusOK, respData.SuggestionData{Suggestions: suggestions})
}

// GET /player/settings
func (api *APIServer) handleGetPlayerSetings(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	playerGames, err := api.storage.GetPlayerGames(p)
	if err != nil {
		return api.HandleError(err)
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

// POST /game/session/new
func (api *APIServer) handleStartNewGameSession(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
	if p.CurrentGame.GMID != p.ID {
		return api.HandleErrorString("only GM may start new session").WithCode(http.StatusForbidden)
	}

	newSession, err := api.storage.StartNewGameSession(p.CurrentGame)
	if err != nil {
		return api.HandleError(err)
	}

	return api.Respond(r, w, http.StatusCreated, newSession)
}

// GET /image/{type}/{id}
func (api *APIServer) handleGetImage(w http.ResponseWriter, r *http.Request) *APIError {
	// ++ add permissions by player ++ //
	imageType := r.PathValue("type")
	imageID := getPathValueInt(r, "id")
	if imageType == "" || imageID == 0 {
		return api.HandleErrorString("image type and id cannot be empty or 0").WithCode(http.StatusBadRequest)
	}

	params := fmt.Sprintf("%s_%d", imageType, imageID)
	uri := fmt.Sprintf("%s/file/%s/%s", api.fileServer.Addr, api.fileServer.Proj, params)

	req, _ := http.NewRequest(r.Method, uri, nil)
	req.Header.Add("Authorization", api.fileServer.Pass)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return api.HandleErrorString(fmt.Sprintf("cannot send image get request: %v", err))
	}

	if res.StatusCode == 200 || res.StatusCode == 201 {
		resBody, _ := io.ReadAll(res.Body)
		return api.Respond(r, w, res.StatusCode, string(resBody))
	} else {
		resBody, _ := io.ReadAll(res.Body)
		return api.HandleErrorString(fmt.Sprintf("file server error: %s", string(resBody))).WithCode(res.StatusCode)
	}
}

// POST /image/{type}/{id}
func (api *APIServer) handlePostImage(w http.ResponseWriter, r *http.Request) *APIError {
	imageType := r.PathValue("type")
	imageID := getPathValueInt(r, "id")

	if imageType == "" || imageID == 0 {
		return api.HandleErrorString("image type and id cannot be empty or 0").WithCode(http.StatusBadRequest)
	}

	// ++ add permissions by player ++ //

	params := fmt.Sprintf("%s_%d", imageType, imageID)
	uri := fmt.Sprintf("%s/file/%s/%s", api.fileServer.Addr, api.fileServer.Proj, params)

	maxSize := int64(4 * 1024 * 1024)
	// body, err := io.ReadAll(io.LimitReader(r.Body, maxSize))
	// if err != nil {
	// 	return api.HandleErrorString(fmt.Sprintf("error reading request body: %s", err))
	// }

	// req, err := http.NewRequest("POST", uri, bytes.NewReader(body))
	// if err != nil {
	// 	return api.HandleErrorString(fmt.Sprintf("error creating forward request: %s", err))
	// }

	// Create pipe for streaming
	pr, pw := io.Pipe()
	defer pr.Close()

	// Prepare the outgoing request
	req, _ := http.NewRequest(r.Method, uri, pr)
	req.Header = r.Header
	req.Header.Add("Authorization", api.fileServer.Pass)

	// Stream with exact size enforcement
	go func() {
		defer pw.Close()
		written, _ := io.CopyN(pw, r.Body, maxSize+1)

		if written > maxSize {
			pw.CloseWithError(http.ErrBodyNotAllowed)
			return
		}
	}()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if err == http.ErrBodyNotAllowed {
			return api.HandleErrorString("file too large").WithCode(http.StatusRequestEntityTooLarge)
		} else {
			return api.HandleErrorString(fmt.Sprintf("cannot send image post request: %v", err))
		}

	}

	// Close body to not to log image sent
	defer r.Body.Close()

	if res.StatusCode == 200 || res.StatusCode == 201 {
		return api.Respond(r, w, res.StatusCode, nil)
	} else {
		resBody, _ := io.ReadAll(res.Body)
		return api.HandleErrorString(fmt.Sprintf("file server error: %s", string(resBody)))
	}
}
