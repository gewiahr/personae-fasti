package handler

import (
	"net/http"

	"personae-fasti/internal/dto"
	"personae-fasti/internal/mapper"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"
)

type RecordHandler struct {
	svc *service.RecordService
}

func NewRecordHandler(svc *service.RecordService) *RecordHandler {
	return &RecordHandler{svc: svc}
}

// GetPlayerCurrentGameRecords handles GET /records (protected).
func (h *RecordHandler) GetPlayerCurrentGameRecords(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	records, err := h.svc.GetPlayerCurrentGameRecords(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	resp := mapper.RecordToRecordFullArray(records, req.Player.CurrentGame.Players, req.Player.CurrentGame.Quests, req.Player.CurrentGame.ExtID)
	return httputils.Response{Status: http.StatusOK, Body: struct {
		Records []dto.RecordFull `json:"records"`
	}{
		Records: resp,
	}}
}

// PostPlayerCurrentGameRecord handles POST /record (protected).
func (h *RecordHandler) PostPlayerCurrentGameRecord(req httputils.RequestData[dto.RecordInsert]) httputils.Responder {
	record, err := h.svc.PostPlayerCurrentGameRecord(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	resp := mapper.RecordToRecordFull(*record, req.Player.CurrentGame.Players, req.Player.CurrentGame.Quests, req.Player.CurrentGame.ExtID)
	return httputils.Response{Status: http.StatusCreated, Body: resp}
}

// EditPlayerCurrentGameRecord handles PUT /record (protected).
func (h *RecordHandler) EditPlayerCurrentGameRecord(req httputils.RequestData[dto.RecordUpdate]) httputils.Responder {
	record, err := h.svc.EditPlayerCurrentGameRecord(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	resp := mapper.RecordToRecordFull(*record, req.Player.CurrentGame.Players, req.Player.CurrentGame.Quests, req.Player.CurrentGame.ExtID)
	return httputils.Response{Status: http.StatusCreated, Body: resp}
}

// DeletePlayerCurrentGameRecord handles DELETE /record/{ext} (protected).
func (h *RecordHandler) DeletePlayerCurrentGameRecord(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	recordExt := req.Request.PathValue("ext")
	if recordExt == "" {
		return e.NewApiError(http.StatusBadRequest, "record ext is invalid")
	}

	err := h.svc.DeletePlayerCurrentGameRecord(req.Context, req.Player, recordExt)
	if err != nil {
		return e.ErrToApiError(err)
	}

	return httputils.Response{Status: http.StatusCreated, Body: nil}
}
