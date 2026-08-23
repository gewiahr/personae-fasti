package service

import (
	"context"
	"testing"

	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/repo"
)

type questRepoStub struct {
	repo.QuestRepository
	quest        *domain.Quest
	editCalled   bool
	finishCalled bool
	resetCalled  bool
	writeGameID  int
}

func (r *questRepoStub) GetPlayerCurrentGameQuestByID(_ context.Context, gameID, questID int) (*domain.Quest, error) {
	if r.quest == nil || r.quest.ID != questID || r.quest.GameID != gameID {
		return nil, nil
	}
	copy := *r.quest
	copy.Tasks = append([]domain.QuestTask(nil), r.quest.Tasks...)
	return &copy, nil
}

func (r *questRepoStub) GetPlayerCurrentGameQuestByExt(_ context.Context, gameID int, ext string) (*domain.Quest, error) {
	if r.quest == nil || r.quest.ExtID != ext || r.quest.GameID != gameID {
		return nil, nil
	}
	copy := *r.quest
	copy.Tasks = append([]domain.QuestTask(nil), r.quest.Tasks...)
	return &copy, nil
}

func (r *questRepoStub) FilterAllowedTasks(_ context.Context, tasks []domain.QuestTask, playerID int) ([]domain.QuestTask, error) {
	allowed := make([]domain.QuestTask, 0, len(tasks))
	for _, task := range tasks {
		if task.HiddenBy == 0 || task.HiddenBy == playerID {
			allowed = append(allowed, task)
		}
	}
	return allowed, nil
}

func (r *questRepoStub) EditPlayerCurrentGameQuest(_ context.Context, update *dto.QuestUpdateData, playerID, gameID int) (*domain.Quest, error) {
	r.editCalled = true
	r.writeGameID = gameID
	return &domain.Quest{ID: update.Quest.ID, GameID: gameID, HiddenBy: hiddenBy(update.Quest.Hidden, playerID)}, nil
}

func (r *questRepoStub) FinishPlayerCurrentGameQuest(_ context.Context, questID, gameID, playerID int, successful bool) (*domain.Quest, error) {
	r.finishCalled = true
	r.writeGameID = gameID
	return &domain.Quest{ID: questID, GameID: gameID, Successful: successful}, nil
}

func (r *questRepoStub) ResetPlayerCurrentGameQuest(_ context.Context, questID, gameID, playerID int) (*domain.Quest, error) {
	r.resetCalled = true
	r.writeGameID = gameID
	return &domain.Quest{ID: questID, GameID: gameID}, nil
}

func TestEditQuestAllowsVisibleAndOwnHiddenSameGameQuest(t *testing.T) {
	for _, test := range []struct {
		name     string
		hiddenBy int
	}{
		{name: "visible", hiddenBy: 0},
		{name: "own hidden", hiddenBy: 12},
	} {
		t.Run(test.name, func(t *testing.T) {
			repository := &questRepoStub{quest: &domain.Quest{ID: 5, ExtID: "quest-ext", GameID: 8, HiddenBy: test.hiddenBy}}
			service := NewQuestService(repository, nil)
			player := &domain.Player{ID: 12, CurrentGameID: 8}

			_, err := service.EditPlayerCurrentGameQuest(context.Background(), player, &dto.QuestUpdateData{Quest: dto.QuestUpdate{ExtID: "quest-ext"}})
			if err != nil {
				t.Fatalf("edit returned an unexpected error: %v", err)
			}
			if !repository.editCalled {
				t.Fatal("expected repository edit to be called")
			}
			if repository.writeGameID != player.CurrentGameID {
				t.Fatalf("edit used game %d, want %d", repository.writeGameID, player.CurrentGameID)
			}
		})
	}
}

func TestQuestMutationsRejectQuestHiddenByAnotherPlayerBeforeWrite(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*QuestService, *domain.Player) error
		called func(*questRepoStub) bool
	}{
		{
			name: "edit",
			mutate: func(s *QuestService, p *domain.Player) error {
				_, err := s.EditPlayerCurrentGameQuest(context.Background(), p, &dto.QuestUpdateData{Quest: dto.QuestUpdate{ExtID: "quest-ext"}})
				return err
			},
			called: func(r *questRepoStub) bool { return r.editCalled },
		},
		{
			name: "finish",
			mutate: func(s *QuestService, p *domain.Player) error {
				_, err := s.FinishPlayerCurrentGameQuest(context.Background(), p, "quest-ext", true)
				return err
			},
			called: func(r *questRepoStub) bool { return r.finishCalled },
		},
		{
			name: "reset",
			mutate: func(s *QuestService, p *domain.Player) error {
				_, err := s.ResetPlayerCurrentGameQuest(context.Background(), p, "quest-ext")
				return err
			},
			called: func(r *questRepoStub) bool { return r.resetCalled },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &questRepoStub{quest: &domain.Quest{ID: 5, ExtID: "quest-ext", GameID: 8, HiddenBy: 99}}
			player := &domain.Player{ID: 12, CurrentGameID: 8}
			err := test.mutate(NewQuestService(repository, nil), player)
			assertAppErrorType(t, err, e.ErrForbidden)
			if test.called(repository) {
				t.Fatal("repository mutation was called before authorization")
			}
		})
	}
}

func TestEditQuestRejectsCrossGameIDBeforeWrite(t *testing.T) {
	repository := &questRepoStub{quest: &domain.Quest{ID: 5, ExtID: "quest-ext", GameID: 9}}
	service := NewQuestService(repository, nil)
	player := &domain.Player{ID: 12, CurrentGameID: 8}

	_, err := service.EditPlayerCurrentGameQuest(context.Background(), player, &dto.QuestUpdateData{Quest: dto.QuestUpdate{ExtID: "quest-ext"}})
	assertAppErrorType(t, err, e.ErrNotFound)
	if repository.editCalled {
		t.Fatal("repository edit was called for a cross-game quest")
	}
}

func TestEditQuestRejectsHiddenTaskIDBeforeWrite(t *testing.T) {
	repository := &questRepoStub{quest: &domain.Quest{
		ID:     5,
		ExtID:  "quest-ext",
		GameID: 8,
		Tasks: []domain.QuestTask{
			{ID: 10, GameID: 8, QuestID: 5},
			{ID: 11, GameID: 8, QuestID: 5, HiddenBy: 99},
		},
	}}
	service := NewQuestService(repository, nil)
	player := &domain.Player{ID: 12, CurrentGameID: 8}
	update := &dto.QuestUpdateData{
		Quest: dto.QuestUpdate{ExtID: "quest-ext"},
		Tasks: []dto.TaskUpdate{{ID: 10}, {ID: 11}},
	}

	_, err := service.EditPlayerCurrentGameQuest(context.Background(), player, update)
	assertAppErrorType(t, err, e.ErrForbidden)
	if repository.editCalled {
		t.Fatal("repository edit was called with another player's hidden task")
	}
}
