package api

import "personae-fasti/internal/handler"

// TODO: Sort and move endpoints to appropriate handlers.
func (api *APIServer) SetHandlers(
	authHandler *handler.AuthHandler,
	gameHandler *handler.GameHandler,
	playerHandler *handler.PlayerHandler,
	recordHandler *handler.RecordHandler,
	entitiesHandler *handler.EntitiesHandler,
	questHandler *handler.QuestHandler,
	appHandler *handler.AppHandler,
	imageHandler *handler.ImageHandler,
) {
	/* LOGIN */
	api.router.HandleFunc("GET /login", AuthAdapt(api.auth, authHandler.LoginByToken))
	api.router.HandleFunc("GET /login/{username}", Adapt(authHandler.LoginByUsername))
	api.router.HandleFunc("POST /login", Adapt(authHandler.Login))
	api.router.HandleFunc("POST /signup", Adapt(authHandler.SignUp))

	/* RECORDS */
	api.router.HandleFunc("GET /records", AuthAdapt(api.auth, recordHandler.GetPlayerCurrentGameRecords))
	api.router.HandleFunc("POST /record", AuthAdapt(api.auth, recordHandler.PostPlayerCurrentGameRecord))
	api.router.HandleFunc("PUT /record", AuthAdapt(api.auth, recordHandler.EditPlayerCurrentGameRecord))
	api.router.HandleFunc("DELETE /record/{id}", AuthAdapt(api.auth, recordHandler.DeletePlayerCurrentGameRecord))

	/* CHARS */
	api.router.HandleFunc("GET /chars", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameChars))
	api.router.HandleFunc("GET /char/{ext}", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameCharByExt))
	api.router.HandleFunc("POST /char", AuthAdapt(api.auth, entitiesHandler.PostPlayerCurrentGameChar))
	api.router.HandleFunc("PUT /char", AuthAdapt(api.auth, entitiesHandler.EditPlayerCurrentGameChar))
	// TODO: Add DELETE /char/{ext}.

	/* NPCS */
	api.router.HandleFunc("GET /npcs", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameNPCs))
	api.router.HandleFunc("GET /npc/{ext}", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameNPCByExt))
	api.router.HandleFunc("POST /npc", AuthAdapt(api.auth, entitiesHandler.PostPlayerCurrentGameNPC))
	api.router.HandleFunc("PUT /npc", AuthAdapt(api.auth, entitiesHandler.EditPlayerCurrentGameNPC))
	// TODO: Add DELETE /npc/{ext}.

	/* LOCATIONS */
	api.router.HandleFunc("GET /locations", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameLocations))
	api.router.HandleFunc("GET /location/{ext}", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameLocationByExt))
	api.router.HandleFunc("POST /location", AuthAdapt(api.auth, entitiesHandler.PostPlayerCurrentGameLocation))
	api.router.HandleFunc("PUT /location", AuthAdapt(api.auth, entitiesHandler.EditPlayerCurrentGameLocation))
	// TODO: Add DELETE /location/{ext}.

	/* SUGGESTIONS */
	api.router.HandleFunc("GET /suggestions", AuthAdapt(api.auth, entitiesHandler.GetPlayerCurrentGameSuggestions))

	/* QUESTS */
	api.router.HandleFunc("GET /quests", AuthAdapt(api.auth, questHandler.GetPlayerCurrentGameQuests))
	api.router.HandleFunc("GET /quest/{ext}", AuthAdapt(api.auth, questHandler.GetPlayerCurrentGameQuestByExt))
	api.router.HandleFunc("POST /quest", AuthAdapt(api.auth, questHandler.PostPlayerCurrentGameQuest))
	api.router.HandleFunc("PUT /quest", AuthAdapt(api.auth, questHandler.EditPlayerCurrentGameQuest))
	// TODO: Add DELETE /quest/{ext}.
	api.router.HandleFunc("PATCH /quest/{ext}/complete", AuthAdapt(api.auth, questHandler.CompletePlayerCurrentGameQuest))
	api.router.HandleFunc("PATCH /quest/{ext}/fail", AuthAdapt(api.auth, questHandler.FailPlayerCurrentGameQuest))
	api.router.HandleFunc("PATCH /quest/{ext}/reset", AuthAdapt(api.auth, questHandler.ResetPlayerCurrentGameQuest))
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
	api.router.HandleFunc("PUT /player/currentGame/{gameExt}", AuthAdapt(api.auth, playerHandler.ChangePlayerCurrentGame))
	api.router.HandleFunc("POST /player/invite/accept/{inviteCode}", AuthAdapt(api.auth, playerHandler.AcceptGameInvite))
	api.router.HandleFunc("POST /player/invite/refuse/{inviteCode}", AuthAdapt(api.auth, playerHandler.RefuseGameInvite))
	api.router.HandleFunc("GET /player/settings", AuthAdapt(api.auth, playerHandler.GetPlayerSettings))
	api.router.HandleFunc("GET /notes", AuthAdapt(api.auth, playerHandler.GetPersonalNote))
	api.router.HandleFunc("PUT /notes", AuthAdapt(api.auth, playerHandler.UpdatePersonalNote))

	/* IMAGES */
	api.router.HandleFunc("GET /entities/{type}/{ext}/images", AuthAdapt(api.auth, imageHandler.List))
	api.router.HandleFunc("POST /entities/{type}/{ext}/images/external", AuthAdapt(api.auth, imageHandler.CreateExternal))
	api.router.HandleFunc("POST /entities/{type}/{ext}/images", ImageAdapt(api.auth, 32<<20, api.imageUploadTimeout, imageHandler.Upload))
	api.router.HandleFunc("DELETE /images/{imageExt}", AuthAdapt(api.auth, imageHandler.Delete))
	api.router.HandleFunc("PATCH /images/{imageExt}/main", AuthAdapt(api.auth, imageHandler.SetMain))
	api.router.HandleFunc("GET /game/storage/quota", AuthAdapt(api.auth, imageHandler.GetQuota))

	/* APPLICATION */
	api.router.HandleFunc("POST /feedback", AuthAdapt(api.auth, appHandler.LoginByToken))
}
