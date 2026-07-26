package service

import (
	"context"
	"database/sql"
	"personae-fasti/internal/domain"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/repo"
)

type GameService struct {
	recordRepo repo.RecordRepository
	playerRepo repo.PlayerRepository
	gameRepo   repo.GameRepository
}

func NewGameService(playerRepo repo.PlayerRepository, gameRepo repo.GameRepository, recordRepo repo.RecordRepository) *GameService {
	return &GameService{
		recordRepo: recordRepo,
		playerRepo: playerRepo,
		gameRepo:   gameRepo,
	}
}

func (s *GameService) GetPlayerCurrentGame(ctx context.Context, player *domain.Player) (*domain.Game, error) {
	game, err := s.gameRepo.GetCurrentGame(ctx, player.CurrentGameID)
	if err == sql.ErrNoRows {
		return nil, e.NewNotFoundError("Нет текущей игры")
	} else if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных текущей игры", err)
	}
	return game, nil
}

func (s *GameService) GetPlayerSettings(ctx context.Context, player *domain.Player) (*domain.Player, error) {
	p, err := s.playerRepo.GetPlayerWithGames(ctx, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных игрока", err)
	}

	return p, nil
}
