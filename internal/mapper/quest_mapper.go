package mapper

import (
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
)

func QuestToQuestBriefArray(quests []domain.Quest, gameExt string) []dto.QuestBrief {
	questBriefArray := []dto.QuestBrief{}
	for _, quest := range quests {
		questBriefArray = append(questBriefArray, QuestToQuestBrief(quest, gameExt))
	}

	return questBriefArray
}

func QuestToQuestBrief(quest domain.Quest, gameExt string) dto.QuestBrief {
	return dto.QuestBrief{
		ID:         quest.ID,
		Name:       quest.Name,
		Title:      quest.Title,
		GameExtID:  gameExt,
		Successful: quest.Successful,
		Hidden:     quest.HiddenBy != 0,
		Finished:   quest.Finished != nil,
	}
}

func QuestToQuestFullInfo(quest *domain.Quest, gameExt string) *dto.QuestFull {
	return &dto.QuestFull{
		ID:          quest.ID,
		Name:        quest.Name,
		Title:       quest.Title,
		Description: quest.Description,
		ParentID:    quest.ParentID,
		ChildID:     quest.ChildID,
		HeadID:      quest.HeadID,
		GameExtID:   gameExt,
		Successful:  quest.Successful,
		Hidden:      quest.HiddenBy != 0,
		Finished:    quest.Finished != nil,
	}
}

func TaskToTaskFullInfoArray(tasks []domain.QuestTask, gameExt string) []dto.QuestTaskFull {
	taskInfoArray := []dto.QuestTaskFull{}
	for _, task := range tasks {
		finishedTask := false
		if task.Finished != nil {
			finishedTask = true
		}
		taskInfoArray = append(taskInfoArray, dto.QuestTaskFull{
			ID:          task.ID,
			QuestID:     task.QuestID,
			Name:        task.Name,
			Description: task.Description,
			GameExtID:   gameExt,
			Type:        int(task.Type),
			Capacity:    task.Capacity,
			Current:     task.Current,
			Hidden:      task.HiddenBy != 0,
			Finished:    finishedTask,
		})
	}

	return taskInfoArray
}
