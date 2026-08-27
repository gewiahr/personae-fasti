package handler

import (
	"net/http"
	"personae-fasti/internal/dto"
	"personae-fasti/internal/mapper"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"
)

type QuestHandler struct {
	svc *service.QuestService
}

func NewQuestHandler(svc *service.QuestService) *QuestHandler {
	return &QuestHandler{svc: svc}
}

// GetPlayerCurrentGameQuests handles GET /quests (protected).
func (h *QuestHandler) GetPlayerCurrentGameQuests(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	quests, err := h.svc.GetPlayerCurrentGameQuests(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}

	return httputils.Response{Status: http.StatusOK, Body: struct {
		Quests []dto.QuestBrief `json:"quests"`
	}{
		Quests: mapper.QuestToQuestBriefArray(quests, req.Player.CurrentGame.ExtID),
	}}
}

// GetPlayerCurrentGameQuestByExt handles GET /quest/{ext} (protected).
func (h *QuestHandler) GetPlayerCurrentGameQuestByExt(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	quest, err := h.svc.GetPlayerCurrentGameQuestByExt(req.Context, req.Player, req.Request.PathValue("ext"))
	if err != nil {
		return e.ErrToApiError(err)
	}

	questFull := mapper.QuestToQuestFullInfo(quest, req.Player.CurrentGame.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: dto.QuestPage{
		Quest:   *questFull,
		Tasks:   mapper.TaskToTaskFullInfoArray(quest.Tasks, req.Player.CurrentGame.ExtID),
		Records: mapper.RecordToRecordFullArray(quest.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// PostPlayerCurrentGameQuest handles POST /quest (protected).
func (h *QuestHandler) PostPlayerCurrentGameQuest(req httputils.RequestData[dto.QuestCreateData]) httputils.Responder {
	quest, err := h.svc.PostPlayerCurrentGameQuest(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	questFull := mapper.QuestToQuestFullInfo(quest, req.Player.CurrentGame.ExtID)
	return httputils.Response{Status: http.StatusCreated, Body: dto.QuestPage{
		Quest:   *questFull,
		Tasks:   mapper.TaskToTaskFullInfoArray(quest.Tasks, req.Player.CurrentGame.ExtID),
		Records: mapper.RecordToRecordFullArray(quest.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// EditPlayerCurrentGameQuest handles PUT /quest (protected).
func (h *QuestHandler) EditPlayerCurrentGameQuest(req httputils.RequestData[dto.QuestUpdateData]) httputils.Responder {
	quest, err := h.svc.EditPlayerCurrentGameQuest(req.Context, req.Player, &req.Body)
	if err != nil {
		return e.ErrToApiError(err)
	}

	questFull := mapper.QuestToQuestFullInfo(quest, req.Player.CurrentGame.ExtID)
	return httputils.Response{Status: http.StatusCreated, Body: dto.QuestPage{
		Quest:   *questFull,
		Tasks:   mapper.TaskToTaskFullInfoArray(quest.Tasks, req.Player.CurrentGame.ExtID),
		Records: mapper.RecordToRecordFullArray(quest.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// CompletePlayerCurrentGameQuest handles PATCH /quest/{id}/complete (protected).
func (h *QuestHandler) CompletePlayerCurrentGameQuest(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	quest, err := h.svc.FinishPlayerCurrentGameQuest(req.Context, req.Player, req.Request.PathValue("ext"), true)
	if err != nil {
		return e.ErrToApiError(err)
	}

	questFull := mapper.QuestToQuestFullInfo(quest, req.Player.CurrentGame.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: dto.QuestPage{
		Quest:   *questFull,
		Tasks:   mapper.TaskToTaskFullInfoArray(quest.Tasks, req.Player.CurrentGame.ExtID),
		Records: mapper.RecordToRecordFullArray(quest.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// FailPlayerCurrentGameQuest handles PATCH /quest/{id}/fail (protected).
func (h *QuestHandler) FailPlayerCurrentGameQuest(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	quest, err := h.svc.FinishPlayerCurrentGameQuest(req.Context, req.Player, req.Request.PathValue("ext"), false)
	if err != nil {
		return e.ErrToApiError(err)
	}

	questFull := mapper.QuestToQuestFullInfo(quest, req.Player.CurrentGame.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: dto.QuestPage{
		Quest:   *questFull,
		Tasks:   mapper.TaskToTaskFullInfoArray(quest.Tasks, req.Player.CurrentGame.ExtID),
		Records: mapper.RecordToRecordFullArray(quest.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// ResetPlayerCurrentGameQuest handles PATCH /quest/{id}/reset (protected).
func (h *QuestHandler) ResetPlayerCurrentGameQuest(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	quest, err := h.svc.ResetPlayerCurrentGameQuest(req.Context, req.Player, req.Request.PathValue("ext"))
	if err != nil {
		return e.ErrToApiError(err)
	}

	questFull := mapper.QuestToQuestFullInfo(quest, req.Player.CurrentGame.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: dto.QuestPage{
		Quest:   *questFull,
		Tasks:   mapper.TaskToTaskFullInfoArray(quest.Tasks, req.Player.CurrentGame.ExtID),
		Records: mapper.RecordToRecordFullArray(quest.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// UpdatePlayerCurrentGameQuestTasks handles PATCH /quest/tasks (protected).
func (h *QuestHandler) UpdatePlayerCurrentGameQuestTasks(req httputils.RequestData[dto.QuestTasksPatch]) httputils.Responder {
	if req.Body.QuestExtID == "" {
		return e.NewApiError(http.StatusBadRequest, "quest ext is required")
	}

	quest, err := h.svc.GetPlayerCurrentGameQuestByExt(req.Context, req.Player, req.Body.QuestExtID)
	if err != nil {
		return e.ErrToApiError(err)
	}

	tasks, err := h.svc.UpdatePlayerCurrentGameQuestTasks(req.Context, req.Body.Tasks, quest)
	if err != nil {
		return e.ErrToApiError(err)
	}

	tasksArrayFull := mapper.TaskToTaskFullInfoArray(tasks, req.Player.CurrentGame.ExtID)
	return httputils.Response{Status: http.StatusOK, Body: tasksArrayFull}
}
