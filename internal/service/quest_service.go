package service

import (
	"context"
	"fmt"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/repo"
)

type QuestService struct {
	repo       repo.QuestRepository
	recordRepo repo.RecordRepository
	// playerRepo repo.PlayerRepository
	// gameRepo   repo.GameRepository
}

func NewQuestService(
	repo repo.QuestRepository,
	// playerRepo repo.PlayerRepository,
	// gameRepo repo.GameRepository,
	recordRepo repo.RecordRepository,
) *QuestService {
	return &QuestService{
		repo:       repo,
		recordRepo: recordRepo,
		// playerRepo: playerRepo,
		// gameRepo:   gameRepo,
	}
}

func (s *QuestService) GetPlayerCurrentGameQuests(ctx context.Context, player *domain.Player) ([]domain.Quest, error) {
	quests, err := s.repo.GetCurrentGameQuestList(ctx, player.CurrentGameID, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения квестов", err)
	}
	return quests, nil
}

func (s *QuestService) GetPlayerCurrentGameQuestByID(ctx context.Context, player *domain.Player, questID int) (*domain.Quest, error) {
	quest, err := s.repo.GetPlayerCurrentGameQuestByID(ctx, player.CurrentGameID, questID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных квеста", err)
	} else if quest == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no quest with id %d", questID))
	} else if quest.GameID != player.CurrentGameID {
		return nil, e.NewForbiddenError(fmt.Sprintf("quest %d is not allowed to request for the game %d", quest.ID, player.CurrentGameID))
	} else if err := ensureHiddenContentEditable(quest.HiddenBy, player.ID); err != nil {
		return nil, err
	}

	if len(quest.Tasks) > 0 {
		quest.Tasks, err = s.repo.FilterAllowedTasks(ctx, quest.Tasks, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения задач квеста", err)
		}
	}

	if len(quest.Records) > 0 {
		quest.Records, err = s.recordRepo.FilterAllowedRecords(ctx, quest.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей квеста", err)
		}
	}

	return quest, nil
}

func (s *QuestService) PostPlayerCurrentGameQuest(ctx context.Context, player *domain.Player, questCreateData *dto.QuestCreateData) (*domain.Quest, error) {
	quest, err := s.repo.CreatePlayerCurrentGameQuest(ctx, questCreateData, player.ID, player.CurrentGameID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка создания квеста", err)
	} else if quest == nil {
		return nil, e.NewNotFoundError("quest not created")
	}

	return quest, nil
}

func (s *QuestService) EditPlayerCurrentGameQuest(ctx context.Context, player *domain.Player, questUpdateData *dto.QuestUpdateData) (*domain.Quest, error) {
	existing, err := s.GetPlayerCurrentGameQuestByID(ctx, player, questUpdateData.Quest.ID)
	if err != nil {
		return nil, err
	}
	if err := validateQuestTaskEdits(questUpdateData.Tasks, existing.Tasks); err != nil {
		return nil, err
	}

	quest, err := s.repo.EditPlayerCurrentGameQuest(ctx, questUpdateData, player.ID, player.CurrentGameID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных квеста", err)
	} else if quest == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no quest with id %d", questUpdateData.Quest.ID))
	}

	if len(quest.Records) > 0 {
		quest.Records, err = s.recordRepo.FilterAllowedRecords(ctx, quest.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей персонажа", err)
		}
	}

	return quest, nil
}

func (s *QuestService) FinishPlayerCurrentGameQuest(ctx context.Context, player *domain.Player, questID int, successful bool) (*domain.Quest, error) {
	if _, err := s.GetPlayerCurrentGameQuestByID(ctx, player, questID); err != nil {
		return nil, err
	}

	quest, err := s.repo.FinishPlayerCurrentGameQuest(ctx, questID, player.CurrentGameID, player.ID, successful)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных квеста", err)
	}
	/* TODO: add Tasks and Records on demand */
	return quest, nil
}

func (s *QuestService) ResetPlayerCurrentGameQuest(ctx context.Context, player *domain.Player, questID int) (*domain.Quest, error) {
	if _, err := s.GetPlayerCurrentGameQuestByID(ctx, player, questID); err != nil {
		return nil, err
	}

	quest, err := s.repo.ResetPlayerCurrentGameQuest(ctx, questID, player.CurrentGameID, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных квеста", err)
	}
	/* TODO: add Tasks and Records on demand */
	return quest, nil
}

func (s *QuestService) GetPlayerCurrentGameQuestTasks(ctx context.Context, quest *domain.Quest) ([]domain.QuestTask, error) {
	tasks, err := s.repo.GetTasksByQuest(ctx, quest)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения задач квеста", err)
	}
	return tasks, nil
}

func validateQuestTaskEdits(updates []dto.TaskUpdate, editable []domain.QuestTask) error {
	editableIDs := make(map[int]struct{}, len(editable))
	for _, task := range editable {
		editableIDs[task.ID] = struct{}{}
	}

	seen := make(map[int]struct{}, len(updates))
	for _, task := range updates {
		if task.ID == 0 {
			continue
		}
		if _, duplicate := seen[task.ID]; duplicate {
			return e.NewValidationError(fmt.Sprintf("duplicate quest task id %d", task.ID))
		}
		seen[task.ID] = struct{}{}
		if _, allowed := editableIDs[task.ID]; !allowed {
			return e.NewForbiddenError("quest task is not editable by this player")
		}
	}

	return nil
}

func (s *QuestService) UpdatePlayerCurrentGameQuestTasks(ctx context.Context, tasksPatch []dto.TaskPatch, quest *domain.Quest) ([]domain.QuestTask, error) {
	tasks, err := s.repo.UpdateTasks(ctx, tasksPatch, quest)
	if err != nil {
		return nil, e.NewInternalError("Ошибка обновления данных задач квеста", err)
	}
	return tasks, nil
}
