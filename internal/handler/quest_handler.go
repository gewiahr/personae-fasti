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

// GetPlayerCurrentGameQuestByID handles GET /quest/{id} (protected).
func (h *QuestHandler) GetPlayerCurrentGameQuestByID(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	questID := httputils.GetPathValueInt(req.Request, "id")
	if questID <= 0 {
		return e.NewApiError(http.StatusBadRequest, "error parsing id: quest id is invalid")
	}

	quest, err := h.svc.GetPlayerCurrentGameQuestByID(req.Context, req.Player, questID)
	if err != nil {
		return e.ErrToApiError(err)
	}

	questFull := mapper.QuestToQuestFullInfo(quest, req.Player.ExtID)

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
	quest, err := h.svc.GetPlayerCurrentGameQuestByID(req.Context, req.Player, req.Body.Quest.ID)
	if err != nil {
		return e.ErrToApiError(err)
	}

	quest, err = h.svc.EditPlayerCurrentGameQuest(req.Context, req.Player, &req.Body)
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

// RemovePlayerCurrentGameQuest handles DELETE /quest/{id} (protected).
// func (h *QuestHandler) RemovePlayerCurrentGameQuest(req httputils.RequestData[dto.NoBody]) httputils.Responder {
// 	questID := httputils.GetPathValueInt(req.Request, "id")
// 	if questID <= 0 {
// 		return e.NewApiError(http.StatusBadRequest, "error parsing id: quest id is invalid")
// 	}
//
// 	err := api.storage.DeleteQuest(questID, p)
// 	if err != nil {
// 		return api.HandleError(err)
// 	}
//
// 	return api.Respond(r, w, http.StatusOK, nil)
// }

// CompletePlayerCurrentGameQuest handles PATCH /quest/{id}/complete (protected).
func (h *QuestHandler) CompletePlayerCurrentGameQuest(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	questID := httputils.GetPathValueInt(req.Request, "id")
	if questID <= 0 {
		return e.NewApiError(http.StatusBadRequest, "error parsing id: quest id is invalid")
	}

	quest, err := h.svc.FinishPlayerCurrentGameQuest(req.Context, req.Player, questID, true)
	if err != nil {
		return e.ErrToApiError(err)
	}

	questFull := mapper.QuestToQuestFullInfo(quest, req.Player.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: dto.QuestPage{
		Quest:   *questFull,
		Tasks:   mapper.TaskToTaskFullInfoArray(quest.Tasks, req.Player.CurrentGame.ExtID),
		Records: mapper.RecordToRecordFullArray(quest.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// FailPlayerCurrentGameQuest handles PATCH /quest/{id}/fail (protected).
func (h *QuestHandler) FailPlayerCurrentGameQuest(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	questID := httputils.GetPathValueInt(req.Request, "id")
	if questID <= 0 {
		return e.NewApiError(http.StatusBadRequest, "error parsing id: quest id is invalid")
	}

	quest, err := h.svc.FinishPlayerCurrentGameQuest(req.Context, req.Player, questID, false)
	if err != nil {
		return e.ErrToApiError(err)
	}

	questFull := mapper.QuestToQuestFullInfo(quest, req.Player.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: dto.QuestPage{
		Quest:   *questFull,
		Tasks:   mapper.TaskToTaskFullInfoArray(quest.Tasks, req.Player.CurrentGame.ExtID),
		Records: mapper.RecordToRecordFullArray(quest.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// ResetPlayerCurrentGameQuest handles PATCH /quest/{id}/reset (protected).
func (h *QuestHandler) ResetPlayerCurrentGameQuest(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	questID := httputils.GetPathValueInt(req.Request, "id")
	if questID <= 0 {
		return e.NewApiError(http.StatusBadRequest, "error parsing id: quest id is invalid")
	}

	quest, err := h.svc.ResetPlayerCurrentGameQuest(req.Context, req.Player, questID)
	if err != nil {
		return e.ErrToApiError(err)
	}

	questFull := mapper.QuestToQuestFullInfo(quest, req.Player.ExtID)

	return httputils.Response{Status: http.StatusOK, Body: dto.QuestPage{
		Quest:   *questFull,
		Tasks:   mapper.TaskToTaskFullInfoArray(quest.Tasks, req.Player.CurrentGame.ExtID),
		Records: mapper.RecordToRecordFullArray(quest.Records, req.Player.ExtID, req.Player.CurrentGame.ExtID),
	}}
}

// UpdatePlayerCurrentGameQuestTasks handles PATCH /quest/tasks (protected).
func (h *QuestHandler) UpdatePlayerCurrentGameQuestTasks(req httputils.RequestData[dto.QuestTasksPatch]) httputils.Responder {
	if req.Body.QuestID <= 0 {
		return e.NewApiError(http.StatusBadRequest, "error parsing id: quest id is invalid")
	}

	quest, err := h.svc.GetPlayerCurrentGameQuestByID(req.Context, req.Player, req.Body.QuestID)
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
