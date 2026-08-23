package handler

import (
	"net/http"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	"personae-fasti/internal/mapper"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"
)

type EntitiesHandler struct {
	svc *service.EntitiesService
}

func NewEntitiesHandler(svc *service.EntitiesService) *EntitiesHandler {
	return &EntitiesHandler{svc: svc}
}

// GetPlayerCurrentGameChars handles GET /chars (protected)
func (h *EntitiesHandler) GetPlayerCurrentGameChars(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	chars, err := h.svc.GetPlayerCurrentGameChars(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	charBriefArray := mapper.CharToCharBriefArray(chars, req.Player.CurrentGame.ExtID, req.Player.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: struct {
		Chars []dto.CharBrief `json:"chars"`
	}{
		Chars: charBriefArray,
	}}
}

// GetPlayerCurrentGameCharByExt handles GET /char/{ext} (protected)
func (h *EntitiesHandler) GetPlayerCurrentGameCharByExt(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	charExt := req.Request.PathValue("ext")
	char, err := h.svc.GetPlayerCurrentGameCharByExt(req.Context, req.Player, charExt)
	if err != nil {
		return e.ErrToApiError(err)
	}

	charFull := mapper.CharToCharFullInfo(char, req.Player.CurrentGame.ExtID, req.Player.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: dto.CharPage{
		Char:    *charFull,
		Records: mapper.RecordToRecordFullArray(char.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// PostPlayerCurrentGameChar handles POST /char (protected)
func (h *EntitiesHandler) PostPlayerCurrentGameChar(req httputils.RequestData[dto.CharCreate]) httputils.Responder {
	char, err := h.svc.PostPlayerCurrentGameChar(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	charFull := mapper.CharToCharFullInfo(char, req.Player.CurrentGame.ExtID, req.Player.ExtID)
	return httputils.Response{Status: http.StatusCreated, Body: dto.CharPage{
		Char:    *charFull,
		Records: mapper.RecordToRecordFullArray(char.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// // POST /char
// func (api *APIServer) handleCreateChar(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
// 	var charCreate reqData.CharCreate
// 	err := ReadJsonBody(r, &charCreate)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	char, err := api.storage.CreateChar(&charCreate, p)
// 	if err != nil {
// 		return api.HandleError(err).WithCode(http.StatusBadRequest)
// 	}

// 	charFullInfo := respData.CharToCharFullInfo(char)
// 	return api.Respond(r, w, http.StatusCreated, charFullInfo)
// }

// EditPlayerCurrentGameChar handles PUT /char (protected)
func (h *EntitiesHandler) EditPlayerCurrentGameChar(req httputils.RequestData[dto.CharUpdate]) httputils.Responder {
	char, err := h.svc.EditPlayerCurrentGameChar(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	charFull := mapper.CharToCharFullInfo(char, req.Player.CurrentGame.ExtID, req.Player.ExtID)
	return httputils.Response{Status: http.StatusOK, Body: dto.CharPage{
		Char:    *charFull,
		Records: mapper.RecordToRecordFullArray(char.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// // PUT /char
// func (api *APIServer) handleUpdateChar(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
// 	var charUpdate reqData.CharUpdate
// 	err := ReadJsonBody(r, &charUpdate)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	char, err := api.storage.GetCharByID(charUpdate.ID)
// 	if err != nil {
// 		return api.HandleError(err)
// 	} else if char == nil {
// 		return api.HandleErrorString(fmt.Sprintf("no character with id %d", charUpdate.ID)).WithCode(http.StatusNotFound)
// 	} else if char.GameID != p.CurrentGameID {
// 		return api.HandleErrorString(fmt.Sprintf("char %d is not allowed to request for the game %d", char.ID, p.CurrentGameID)).WithCode(http.StatusUnprocessableEntity)
// 	}
// 	// ++ Add char check ++//

// 	char, err = api.storage.UpdateChar(&charUpdate, char, p)
// 	if err != nil {
// 		return api.HandleError(err).WithCode(http.StatusBadRequest)
// 	}

// 	charFullInfo := respData.CharToCharFullInfo(char)
// 	return api.Respond(r, w, http.StatusOK, charFullInfo)
// }

// GetPlayerCurrentGameNPCs handles GET /npcs (protected)
func (h *EntitiesHandler) GetPlayerCurrentGameNPCs(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	npcs, err := h.svc.GetPlayerCurrentGameNPCs(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	npcBriefArray := mapper.NPCToNPCBriefArray(npcs, req.Player.CurrentGame.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: struct {
		NPCs []dto.NPCBrief `json:"npcs"`
	}{
		NPCs: npcBriefArray,
	}}
}

// GetPlayerCurrentGameNPCByExt handles GET /npc/{ext} (protected)
func (h *EntitiesHandler) GetPlayerCurrentGameNPCByExt(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	npcExt := req.Request.PathValue("ext")
	npc, err := h.svc.GetPlayerCurrentGameNPCByExt(req.Context, req.Player, npcExt)
	if err != nil {
		return e.ErrToApiError(err)
	}

	npcFull := mapper.NPCToNPCFullInfo(npc, req.Player.CurrentGame.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: dto.NPCPage{
		NPC:     *npcFull,
		Records: mapper.RecordToRecordFullArray(npc.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// // GET /npc/{id}
// func (api *APIServer) handleGetNPCByID(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
// 	npcID := getPathValueInt(r, "id")
// 	if npcID < 0 {
// 		return api.HandleError(fmt.Errorf("error parsing id: npc id is invalid"))
// 	}

// 	npc, err := api.storage.GetNPCByID(npcID)
// 	if err != nil {
// 		return api.HandleError(err)
// 	} else if npc == nil {
// 		return api.HandleErrorString(fmt.Sprintf("no npc with id %d", npcID)).WithCode(http.StatusNotFound)
// 	} else if npc.GameID != p.CurrentGameID {
// 		return api.HandleErrorString(fmt.Sprintf("npc %d is not allowed to request for the game %d", npc.ID, p.CurrentGameID)).WithCode(http.StatusUnprocessableEntity)
// 	} else if npc.HiddenBy != 0 && npc.HiddenBy != p.ID {
// 		return api.HandleErrorString(fmt.Sprintf("npc %d is not allowed to request for the player %d", npc.ID, p.ID)).WithCode(http.StatusForbidden)
// 	}

// 	records := []data.Record{}
// 	if len(npc.Records) > 0 {
// 		records, err = api.storage.GetAllowedRecords(npc.Records, p.ID)
// 	}

// 	npcPage := respData.NPCPage{
// 		NPC:     *respData.NPCToNPCFullInfo(npc),
// 		Records: records, // ** change to mention API type ** //
// 	}

// 	return api.Respond(r, w, http.StatusOK, npcPage)
// }

// PostPlayerCurrentGameNPC handles POST /npc (protected)
func (h *EntitiesHandler) PostPlayerCurrentGameNPC(req httputils.RequestData[dto.NPCCreate]) httputils.Responder {
	npc, err := h.svc.PostPlayerCurrentGameNPC(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	npcFull := mapper.NPCToNPCFullInfo(npc, req.Player.CurrentGame.ExtID)
	return httputils.Response{Status: http.StatusCreated, Body: dto.NPCPage{
		NPC:     *npcFull,
		Records: mapper.RecordToRecordFullArray(npc.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// // POST /npc
// func (api *APIServer) handleCreateNPC(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
// 	var npcCreate reqData.NPCCreate
// 	err := ReadJsonBody(r, &npcCreate)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	npc, err := api.storage.CreateNPC(&npcCreate, p)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	npcFullInfo := respData.NPCToNPCFullInfo(npc)
// 	return api.Respond(r, w, http.StatusCreated, npcFullInfo)
// }

// EditPlayerCurrentGameNPC handles PUT /npc (protected)
func (h *EntitiesHandler) EditPlayerCurrentGameNPC(req httputils.RequestData[dto.NPCUpdate]) httputils.Responder {
	npc, err := h.svc.EditPlayerCurrentGameNPC(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	npcFull := mapper.NPCToNPCFullInfo(npc, req.Player.CurrentGame.ExtID)
	return httputils.Response{Status: http.StatusOK, Body: dto.NPCPage{
		NPC:     *npcFull,
		Records: mapper.RecordToRecordFullArray(npc.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// // PUT /npc
// func (api *APIServer) handleUpdateNPC(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
// 	var npcUpdate reqData.NPCUpdate
// 	err := ReadJsonBody(r, &npcUpdate)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	npc, err := api.storage.GetNPCByID(npcUpdate.ID)
// 	if err != nil {
// 		return api.HandleError(err)
// 	} else if npc == nil {
// 		return api.HandleErrorString(fmt.Sprintf("no npc with id %d", npcUpdate.ID)).WithCode(http.StatusNotFound)
// 	} else if npc.GameID != p.CurrentGameID {
// 		return api.HandleErrorString(fmt.Sprintf("npc %d is not allowed to request for the game %d", npc.ID, p.CurrentGameID)).WithCode(http.StatusUnprocessableEntity)
// 	}
// 	// ++ Add char check ++//

// 	npc, err = api.storage.UpdateNPC(&npcUpdate, npc, p)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	npcFullInfo := respData.NPCToNPCFullInfo(npc)
// 	return api.Respond(r, w, http.StatusOK, npcFullInfo)
// }

// GetPlayerCurrentGameLocations handles GET /locations (protected)
func (h *EntitiesHandler) GetPlayerCurrentGameLocations(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	locations, err := h.svc.GetPlayerCurrentGameLocations(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	locationBriefArray := mapper.LocationToLocationBriefArray(locations, req.Player.CurrentGame.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: struct {
		Locations []dto.LocationBrief `json:"locations"`
	}{
		Locations: locationBriefArray,
	}}
}

// GetPlayerCurrentGameLocationByExt handles GET /location/{ext} (protected)
func (h *EntitiesHandler) GetPlayerCurrentGameLocationByExt(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	locationExt := req.Request.PathValue("ext")
	location, err := h.svc.GetPlayerCurrentGameLocationByExt(req.Context, req.Player, locationExt)
	if err != nil {
		return e.ErrToApiError(err)
	}

	// ** Cut hidden locations ** //
	var locationChildren []domain.Location
	locationChildren, err = h.svc.GetPlayerCurrentGameLocationChildrenByID(req.Context, req.Player, location.ID)
	if err != nil {
		return e.ErrToApiError(err)
	}
	// ** Cut hidden locations ** //

	locationPage := dto.LocationPage{
		Location: *mapper.LocationToLocationFullInfo(location, req.Player.CurrentGame.ExtID),
		Records:  mapper.RecordToRecordFullArray(location.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
		Includes: mapper.LocationToLocationBriefArray(locationChildren, req.Player.CurrentGame.ExtID),
	}

	if location.ParentID != 0 {
		locationParent, err := h.svc.GetPlayerCurrentGameLocationByID(req.Context, req.Player, location.ParentID)
		if err != nil {
			return e.ErrToApiError(err)
		}

		if err == nil {
			locationPage.Parent = mapper.LocationToLocationBrief(*locationParent, req.Player.CurrentGame.ExtID)
		} else if err == e.ErrNotFound || err == e.ErrForbidden {
			//locationPage.Parent = nil
		} else {
			return e.ErrToApiError(err)
		}
	}

	return httputils.Response{Status: http.StatusOK, Body: locationPage}
}

// PostPlayerCurrentGameLocation handles POST /location (protected)
func (h *EntitiesHandler) PostPlayerCurrentGameLocation(req httputils.RequestData[dto.LocationCreate]) httputils.Responder {
	location, err := h.svc.PostPlayerCurrentGameLocation(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	// ** Cut hidden locations ** //
	var locationChildren []domain.Location
	locationChildren, err = h.svc.GetPlayerCurrentGameLocationChildrenByID(req.Context, req.Player, location.ID)
	if err != nil {
		return e.ErrToApiError(err)
	}
	// ** Cut hidden locations ** //

	locationPage := dto.LocationPage{
		Location: *mapper.LocationToLocationFullInfo(location, req.Player.CurrentGame.ExtID),
		Records:  mapper.RecordToRecordFullArray(location.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
		Includes: mapper.LocationToLocationBriefArray(locationChildren, req.Player.CurrentGame.ExtID),
	}

	if location.ParentID != 0 {
		locationParent, err := h.svc.GetPlayerCurrentGameLocationByID(req.Context, req.Player, location.ParentID)
		if err != nil {
			return e.ErrToApiError(err)
		}

		if err == nil {
			locationPage.Parent = mapper.LocationToLocationBrief(*locationParent, req.Player.CurrentGame.ExtID)
		} else if err == e.ErrNotFound || err == e.ErrForbidden {
			//locationPage.Parent = nil
		} else {
			return e.ErrToApiError(err)
		}
	}

	return httputils.Response{Status: http.StatusCreated, Body: locationPage}
}

// // POST /location
// func (api *APIServer) handleCreateLocation(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
// 	var locationCreate reqData.LocationCreate
// 	err := ReadJsonBody(r, &locationCreate)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	location, err := api.storage.CreateLocation(&locationCreate, p)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	locationFullInfo := respData.LocationToLocationFullInfo(location)
// 	return api.Respond(r, w, http.StatusCreated, locationFullInfo)
// }

// EditPlayerCurrentGameLocation handles PUT /location (protected)
func (h *EntitiesHandler) EditPlayerCurrentGameLocation(req httputils.RequestData[dto.LocationUpdate]) httputils.Responder {
	location, err := h.svc.EditPlayerCurrentGameLocation(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	// ** Cut hidden locations ** //
	var locationChildren []domain.Location
	locationChildren, err = h.svc.GetPlayerCurrentGameLocationChildrenByID(req.Context, req.Player, location.ID)
	if err != nil {
		return e.ErrToApiError(err)
	}
	// ** Cut hidden locations ** //

	locationPage := dto.LocationPage{
		Location: *mapper.LocationToLocationFullInfo(location, req.Player.CurrentGame.ExtID),
		Records:  mapper.RecordToRecordFullArray(location.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
		Includes: mapper.LocationToLocationBriefArray(locationChildren, req.Player.CurrentGame.ExtID),
	}

	if location.ParentID != 0 {
		locationParent, err := h.svc.GetPlayerCurrentGameLocationByID(req.Context, req.Player, location.ParentID)
		if err != nil {
			return e.ErrToApiError(err)
		}

		if err == nil {
			locationPage.Parent = mapper.LocationToLocationBrief(*locationParent, req.Player.CurrentGame.ExtID)
		} else if err == e.ErrNotFound || err == e.ErrForbidden {
			//locationPage.Parent = nil
		} else {
			return e.ErrToApiError(err)
		}
	}

	return httputils.Response{Status: http.StatusOK, Body: locationPage}
}

// // PUT /location
// func (api *APIServer) handleUpdateLocation(w http.ResponseWriter, r *http.Request, p *data.Player) *APIError {
// 	var locationUpdate reqData.LocationUpdate
// 	err := ReadJsonBody(r, &locationUpdate)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	location, err := api.storage.GetLocationByID(locationUpdate.ID)
// 	if err != nil {
// 		return api.HandleError(err)
// 	} else if location == nil {
// 		return api.HandleErrorString(fmt.Sprintf("no location with id %d", locationUpdate.ID)).WithCode(http.StatusNotFound)
// 	} else if location.GameID != p.CurrentGameID {
// 		return api.HandleErrorString(fmt.Sprintf("location %d is not allowed to request for the game %d", location.ID, p.CurrentGameID)).WithCode(http.StatusUnprocessableEntity)
// 	}
// 	// ++ Add char check ++//

// 	location, err = api.storage.UpdateLocation(&locationUpdate, location, p)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}

// 	locationFullInfo := respData.LocationToLocationFullInfo(location)
// 	return api.Respond(r, w, http.StatusOK, locationFullInfo)
// }

// GetPlayerCurrentGameSuggestions handles GET /suggestions (protected)
func (h *EntitiesHandler) GetPlayerCurrentGameSuggestions(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	suggestions, err := h.svc.GetPlayerCurrentGameSuggestions(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	return httputils.Response{Status: http.StatusOK, Body: struct {
		Suggestions []dto.Suggestion `json:"entities"`
	}{
		Suggestions: suggestions,
	}}
}
