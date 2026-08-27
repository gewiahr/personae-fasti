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

	var locationChildren []domain.Location
	locationChildren, err = h.svc.GetPlayerCurrentGameLocationChildrenByID(req.Context, req.Player, location.ID)
	if err != nil {
		return e.ErrToApiError(err)
	}

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

		locationPage.Parent = mapper.LocationToLocationBrief(*locationParent, req.Player.CurrentGame.ExtID)
	}

	return httputils.Response{Status: http.StatusOK, Body: locationPage}
}

// PostPlayerCurrentGameLocation handles POST /location (protected)
func (h *EntitiesHandler) PostPlayerCurrentGameLocation(req httputils.RequestData[dto.LocationCreate]) httputils.Responder {
	location, err := h.svc.PostPlayerCurrentGameLocation(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	var locationChildren []domain.Location
	locationChildren, err = h.svc.GetPlayerCurrentGameLocationChildrenByID(req.Context, req.Player, location.ID)
	if err != nil {
		return e.ErrToApiError(err)
	}

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

		locationPage.Parent = mapper.LocationToLocationBrief(*locationParent, req.Player.CurrentGame.ExtID)
	}

	return httputils.Response{Status: http.StatusCreated, Body: locationPage}
}

// EditPlayerCurrentGameLocation handles PUT /location (protected)
func (h *EntitiesHandler) EditPlayerCurrentGameLocation(req httputils.RequestData[dto.LocationUpdate]) httputils.Responder {
	location, err := h.svc.EditPlayerCurrentGameLocation(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	var locationChildren []domain.Location
	locationChildren, err = h.svc.GetPlayerCurrentGameLocationChildrenByID(req.Context, req.Player, location.ID)
	if err != nil {
		return e.ErrToApiError(err)
	}

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

		locationPage.Parent = mapper.LocationToLocationBrief(*locationParent, req.Player.CurrentGame.ExtID)
	}

	return httputils.Response{Status: http.StatusOK, Body: locationPage}
}

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
