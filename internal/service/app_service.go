package service

import (
	"context"
	"personae-fasti/internal/domain"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/repo"
)

type AppService struct {
	repo repo.AppRepository
}

func NewAppService(repo repo.AppRepository) *AppService {
	return &AppService{
		repo: repo,
	}
}

func (s *AppService) PostFeedback(ctx context.Context, playerID, gameID int, feedbackType, feedbackText string) (*domain.ServiceFeedback, error) {
	feedback := &domain.ServiceFeedback{
		PlayerID: playerID,
		GameID:   gameID,
		Type:     feedbackType,
		Text:     feedbackText,
	}

	feedback, err := s.repo.InsertFeedback(ctx, feedback)
	if err != nil {
		return nil, e.NewInternalError("feedback insertion failed", err)
	}

	return feedback, nil
}
