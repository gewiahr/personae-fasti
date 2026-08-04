package api

import (
	"personae-fasti/internal/handler"
)

// TODO: Sort and move endpoints to appropriate handlers
func (api *APIServer) SetHandlers(
	authHandler *handler.AuthHandler,
	gameHandler *handler.GameHandler,
	playerHandler *handler.PlayerHandler,
	recordHandler *handler.RecordHandler,
	entitiesHandler *handler.EntitiesHandler,
	questHandler *handler.QuestHandler,
	appHandler *handler.AppHandler,
) {
	/* LOGIN */
	api.router.HandleFunc("GET /login", AuthAdapt(api.auth, authHandler.LoginByToken))
	api.router.HandleFunc("GET /login/{username}", Adapt(authHandler.LoginByUsername))
	// 	router.HandleFunc("POST /login/tg", api.HTTPWrapper(api.handleLoginTG))
	api.router.HandleFunc("POST /login", Adapt(authHandler.Login))
	api.router.HandleFunc("POST /signup", Adapt(authHandler.SignUp))

	/* RECORDS */
	api.router.HandleFunc("GET /records", AuthAdapt(api.auth, recordHandler.GetPlayerCurrentGameRecords))
	api.router.HandleFunc("POST /record", AuthAdapt(api.auth, recordHandler.PostPlayerCurrentGameRecord))
	api.router.HandleFunc("PUT /record", AuthAdapt(api.auth, recordHandler.EditPlayerCurrentGameRecord))
	api.router.HandleFunc("DELETE /record/{id}", AuthAdapt(api.auth, recordHandler.DeletePlayerCurrentGameRecord))

	/* CHARS */
	api.router.HandleFunc("GET /chars", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameChars))
	api.router.HandleFunc("GET /char/{id}", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameCharByID))
	api.router.HandleFunc("POST /char", AuthAdapt(api.auth, entitiesHandler.PostPlayerCurrentGameChar))
	api.router.HandleFunc("PUT /char", AuthAdapt(api.auth, entitiesHandler.EditPlayerCurrentGameChar))
	//api.router.HandleFunc("DELETE /char/{id}", AuthAdapt(api.auth, recordHandler.DeletePlayerCurrentGameChar))

	/* NPCS */
	api.router.HandleFunc("GET /npcs", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameNPCs))
	api.router.HandleFunc("GET /npc/{id}", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameNPCByID))
	api.router.HandleFunc("POST /npc", AuthAdapt(api.auth, entitiesHandler.PostPlayerCurrentGameNPC))
	api.router.HandleFunc("PUT /npc", AuthAdapt(api.auth, entitiesHandler.EditPlayerCurrentGameNPC))
	//api.router.HandleFunc("DELETE /npc/{id}", AuthAdapt(api.auth, recordHandler.DeletePlayerCurrentGameNPC))

	/* LOCATIONS */
	api.router.HandleFunc("GET /locations", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameLocations))
	api.router.HandleFunc("GET /location/{id}", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameLocationByID))
	api.router.HandleFunc("POST /location", AuthAdapt(api.auth, entitiesHandler.PostPlayerCurrentGameLocation))
	api.router.HandleFunc("PUT /location", AuthAdapt(api.auth, entitiesHandler.EditPlayerCurrentGameLocation))
	//api.router.HandleFunc("DELETE /location/{id}", AuthAdapt(api.auth, recordHandler.DeletePlayerCurrentGameLocation))

	/* SUGGESTIONS */
	api.router.HandleFunc("GET /suggestions", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameSuggestions))

	/* QUESTS */
	api.router.HandleFunc("GET /quests", AuthAdapt(api.auth, questHandler.GetPlayerCurrentGameQuests))
	api.router.HandleFunc("GET /quest/{id}", AuthAdapt(api.auth, questHandler.GetPlayerCurrentGameQuestByID))
	api.router.HandleFunc("POST /quest", AuthAdapt(api.auth, questHandler.PostPlayerCurrentGameQuest))
	api.router.HandleFunc("PUT /quest", AuthAdapt(api.auth, questHandler.EditPlayerCurrentGameQuest))
	// api.router.HandleFunc("DELETE /quest/{id}", AuthAdapt(api.auth, questHandler.DeletePlayerCurrentGameQuest))
	api.router.HandleFunc("PATCH /quest/{id}/complete", AuthAdapt(api.auth, questHandler.CompletePlayerCurrentGameQuest))
	api.router.HandleFunc("PATCH /quest/{id}/fail", AuthAdapt(api.auth, questHandler.FailPlayerCurrentGameQuest))
	api.router.HandleFunc("PATCH /quest/{id}/reset", AuthAdapt(api.auth, questHandler.ResetPlayerCurrentGameQuest))

	api.router.HandleFunc("PATCH /quest/tasks", AuthAdapt(api.auth, questHandler.UpdatePlayerCurrentGameQuestTasks))

	/* GAME */
	api.router.HandleFunc("GET /game/{ext}", AuthAdapt(api.auth, gameHandler.GetPlayerGameByExt))
	api.router.HandleFunc("POST /game", AuthAdapt(api.auth, gameHandler.CreateNewGame))
	api.router.HandleFunc("PUT /game", AuthAdapt(api.auth, gameHandler.EditGame))

	api.router.HandleFunc("POST /game/session/new", AuthAdapt(api.auth, gameHandler.StartNewGameSession))
	api.router.HandleFunc("PATCH /game/session", AuthAdapt(api.auth, gameHandler.EditGameSession))
	api.router.HandleFunc("DELETE /game/session/remove", AuthAdapt(api.auth, gameHandler.RemoveLastGameSession))

	api.router.HandleFunc("POST /game/invite/{username}", AuthAdapt(api.auth, gameHandler.InvitePlayer))
	api.router.HandleFunc("DELETE /game/invite/{username}", AuthAdapt(api.auth, gameHandler.RemovePlayerInvite))

	api.router.HandleFunc("PUT /game/settings", AuthAdapt(api.auth, gameHandler.PutGameSettings))

	/* PLAYER */
	api.router.HandleFunc("GET /player/currentGame", AuthAdapt(api.auth, playerHandler.GetPlayerCurrentGame))
	api.router.HandleFunc("PUT /player/currentGame/{id}", AuthAdapt(api.auth, playerHandler.ChangePlayerCurrentGame))

	api.router.HandleFunc("POST /player/invite/accept/{gameID}", AuthAdapt(api.auth, playerHandler.AcceptGameInvite))
	api.router.HandleFunc("POST /player/invite/refuse/{gameID}", AuthAdapt(api.auth, playerHandler.RefuseGameInvite))

	// 	router.HandleFunc("GET /player/username/checkAvailability/{username}", api.HTTPWrapper(api.PlayerWrapper(api.handleCheckUsernameAvailability)))
	// 	router.HandleFunc("PATCH /player/username", api.HTTPWrapper(api.PlayerWrapper(api.handleChangePlayerUsername)))

	api.router.HandleFunc("GET /player/settings", AuthAdapt(api.auth, playerHandler.GetPlayerSettings))

	/* APPLICATION */
	api.router.HandleFunc("POST /feedback", AuthAdapt(api.auth, appHandler.LoginByToken))
}

// func (api *APIServer) SetHandlers(router *http.ServeMux) {

// 	//router.HandleFunc("DELETE /game/player/{username}", api.HTTPWrapper(api.PlayerWrapper(api.handleRemovePlayerFromGame)))

// 	router.HandleFunc("GET /image/{type}/{id}", api.HTTPWrapper(api.handleGetImage))
// 	router.HandleFunc("POST /image/{type}/{id}", api.HTTPWrapper(api.handlePostImage))
// }

// // func (api *APIServer) handleHome(w http.ResponseWriter, r *http.Request) *APIError {
// // 	return api.Respond(r, w, http.StatusOK, nil)
// // }

// // POST /login/TG
// func (api *APIServer) handleLoginTG(w http.ResponseWriter, r *http.Request) *APIError {
// 	var loginTG reqData.LoginTGRequest
// 	err := ReadJsonBody(r, &loginTG)
// 	if err != nil {
// 		return api.HandleError(fmt.Errorf("cannot read token: %v", err)).WithCode(http.StatusUnauthorized)
// 	}

// 	err = tgInitData.Validate(loginTG.InitDataRaw, api.auth.BotToken, 2400000*time.Hour)
// 	if err != nil {
// 		return api.HandleError(fmt.Errorf("token is invalid: %v", err)).WithCode(http.StatusUnauthorized)
// 	}

// 	initData, err := tgInitData.Parse(loginTG.InitDataRaw)
// 	if err != nil {
// 		return api.HandleError(fmt.Errorf("can't parse userdata: %v", err))
// 	}

// 	player, err := api.storage.GetTelegramPlayer(initData.User.ID)
// 	if err == sql.ErrNoRows {
// 		return api.HandleErrorString("haven't registered yet").WithCode(http.StatusNotFound)
// 		// player, err = api.storage.CreateTelegramPlayer(initData)
// 		// if err != nil {
// 		// 	return api.HandleError(err)
// 		// }
// 	} else if err != nil {
// 		return api.HandleError(err)
// 	}

// 	// isMember, err := api.checkTGUserChatMembership(player.Telegram.ID)
// 	// if err != nil {
// 	// 	return api.HandleError(err)
// 	// } else if !isMember {
// 	// 	return api.HandleErrorString("you are not member of the group").WithCode(http.StatusForbidden)
// 	// }

// 	token, err := api.storage.CreateAuthToken(player, api.auth.JWTSecret, time.Duration(api.auth.JWTTokenLifetimeHours)*time.Hour)
// 	if err != nil {
// 		return api.HandleError(err).WithMessage("error creating token")
// 	}

// 	loginInfo := respData.LoginInfo{
// 		Authorization: fmt.Sprintf("TG %s", token),
// 		Player: respData.LoginPlayerInfo{
// 			ID:       player.ID,
// 			Username: player.Username,
// 			Settings: nil,
// 		},
// 		CurrentGame: nil,
// 	}

// 	if player.RegData.UsernameSet {
// 		loginInfo.Player.Settings = &respData.LoginPlayerInfoSettings{
// 			CouldChangeUsername: false,
// 		}
// 	}

// 	if player.CurrentGame != nil {
// 		loginInfo.CurrentGame = respData.GameToGameFullInfo(player.CurrentGame)
// 	}

// 	// response := AuthResponse{
// 	// 	Token: tokenString,
// 	// }

// 	// check subscription

// 	// player, err := api.storage.GetPlayerByAccessKey(accesskey)
// 	// if err != nil {
// 	// 	if err == sql.ErrNoRows {
// 	// 		return api.HandleError(fmt.Errorf("login failed: no user info for the passkey %s", strings.ToLower(accesskey))).WithCode(http.StatusUnauthorized)
// 	// 	} else {
// 	// 		return api.HandleError(err)
// 	// 	}
// 	// }

// 	// loginInfo := respData.LoginInfo{
// 	// 	AccessKey: player.AccessKey,
// 	// 	Player: respData.PlayerInfo{
// 	// 		ID:       player.ID,
// 	// 		Username: player.Username,
// 	// 	},
// 	// 	CurrentGame: *respData.GameToGameFullInfo(player.CurrentGame),
// 	// }

// 	return api.Respond(r, w, http.StatusOK, loginInfo)
// }

// // // DELETE /game/player/{username}
// // func (api *APIServer) handleRemovePlayerFromGame(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
// // 	if p.CurrentGame.GMID != p.ID {
// // 		return api.HandleErrorString("only GM may remove players").WithCode(http.StatusForbidden)
// // 	}

// // 	usernameToRemove := r.PathValue("username")
// // 	if usernameToRemove == "" {
// // 		return api.HandleErrorString("error parsing username: username is invalid").WithCode(http.StatusBadRequest)
// // 	}

// // 	playerToRemove, err := api.storage.GetPlayerByUsername(usernameToRemove)
// // 	if err != nil {
// // 		return api.HandleError(err)
// // 	}
// // 	if playerToRemove == nil {
// // 		return api.HandleErrorString("no player with such username").WithCode(http.StatusNotFound)
// // 	}

// // 	err = api.storage.RemovePlayerFromGame(p.CurrentGame, playerToRemove)
// // 	if err != nil {
// // 		return api.HandleError(err)
// // 	}

// // 	return api.Respond(r, w, http.StatusCreated, nil)
// // }

// // GET /image/{type}/{id}
// func (api *APIServer) handleGetImage(w http.ResponseWriter, r *http.Request) *APIError {
// 	// ++ add permissions by player ++ //
// 	imageType := r.PathValue("type")
// 	imageID := getPathValueInt(r, "id")
// 	if imageType == "" || imageID == 0 {
// 		return api.HandleErrorString("image type and id cannot be empty or 0").WithCode(http.StatusBadRequest)
// 	}

// 	params := fmt.Sprintf("%s_%d", imageType, imageID)
// 	uri := fmt.Sprintf("%s/file/%s/%s", api.fileServer.Addr, api.fileServer.Proj, params)

// 	req, _ := http.NewRequest(r.Method, uri, nil)
// 	req.Header.Add("Authorization", api.fileServer.Pass)

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return api.HandleErrorString(fmt.Sprintf("cannot send image get request: %v", err))
// 	}

// 	if res.StatusCode == 200 || res.StatusCode == 201 {
// 		resBody, _ := io.ReadAll(res.Body)
// 		return api.Respond(r, w, res.StatusCode, string(resBody))
// 	} else {
// 		resBody, _ := io.ReadAll(res.Body)
// 		return api.HandleErrorString(fmt.Sprintf("file server error: %s", string(resBody))).WithCode(res.StatusCode)
// 	}
// }

// // POST /image/{type}/{id}
// func (api *APIServer) handlePostImage(w http.ResponseWriter, r *http.Request) *APIError {
// 	imageType := r.PathValue("type")
// 	imageID := getPathValueInt(r, "id")

// 	if imageType == "" || imageID == 0 {
// 		return api.HandleErrorString("image type and id cannot be empty or 0").WithCode(http.StatusBadRequest)
// 	}

// 	// ++ add permissions by player ++ //

// 	params := fmt.Sprintf("%s_%d", imageType, imageID)
// 	uri := fmt.Sprintf("%s/file/%s/%s", api.fileServer.Addr, api.fileServer.Proj, params)

// 	maxSize := int64(4 * 1024 * 1024)
// 	// body, err := io.ReadAll(io.LimitReader(r.Body, maxSize))
// 	// if err != nil {
// 	// 	return api.HandleErrorString(fmt.Sprintf("error reading request body: %s", err))
// 	// }

// 	// req, err := http.NewRequest("POST", uri, bytes.NewReader(body))
// 	// if err != nil {
// 	// 	return api.HandleErrorString(fmt.Sprintf("error creating forward request: %s", err))
// 	// }

// 	// Create pipe for streaming
// 	pr, pw := io.Pipe()
// 	defer pr.Close()

// 	// Prepare the outgoing request
// 	req, _ := http.NewRequest(r.Method, uri, pr)
// 	req.Header = r.Header
// 	req.Header.Add("Authorization", api.fileServer.Pass)

// 	// Stream with exact size enforcement
// 	go func() {
// 		defer pw.Close()
// 		written, _ := io.CopyN(pw, r.Body, maxSize+1)

// 		if written > maxSize {
// 			pw.CloseWithError(http.ErrBodyNotAllowed)
// 			return
// 		}
// 	}()

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		if err == http.ErrBodyNotAllowed {
// 			return api.HandleErrorString("file too large").WithCode(http.StatusRequestEntityTooLarge)
// 		} else {
// 			return api.HandleErrorString(fmt.Sprintf("cannot send image post request: %v", err))
// 		}

// 	}

// 	// Close body to not to log image sent
// 	defer r.Body.Close()

// 	if res.StatusCode == 200 || res.StatusCode == 201 {
// 		return api.Respond(r, w, res.StatusCode, nil)
// 	} else {
// 		resBody, _ := io.ReadAll(res.Body)
// 		return api.HandleErrorString(fmt.Sprintf("file server error: %s", string(resBody)))
// 	}
// }
