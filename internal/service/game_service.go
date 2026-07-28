package service

import (
	"context"
	"database/sql"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
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

func (s *GameService) GetPlayerGameByExt(ctx context.Context, playerID int, gameExt string) (*domain.Game, error) {
	game, err := s.gameRepo.GetByExt(ctx, gameExt)
	if err == sql.ErrNoRows {
		return nil, e.NewNotFoundError("Нет текущей игры")
	} else if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных текущей игры", err)
	}

	playerParticipate := false
	for _, player := range game.Players {
		if player.ID == playerID {
			playerParticipate = true
			break
		}
	}

	if !playerParticipate {
		return nil, e.NewForbiddenError("Игрок не участвует в данной игре")
	}

	return game, nil
}

func (s *GameService) CreateGame(ctx context.Context, player *domain.Player, title string) (*domain.Game, error) {
	newGame := &domain.Game{
		Name: title,
		GMID: player.ID,
	}

	if valid, message := newGame.ValidateGameTitle(); !valid {
		return nil, e.NewValidationError(message)
	}

	game, err := s.gameRepo.Create(ctx, newGame)
	if err != nil {
		return nil, e.NewInternalError("Невозможно создать новую игру", err)
	}

	var currentGame = player.CurrentGame
	if currentGame == nil {
		// TODO: errors should not prevent endpoint to finish
		player, err = s.playerRepo.ChangeCurrentGame(ctx, player.ID, game.ID)
		if err != nil {
			return nil, e.NewInternalError("Невозможно сменить текущую игру", err)
		}
	}

	return game, nil
}

func (s *GameService) EditGame(ctx context.Context, player *domain.Player, gameEdit *dto.GameUpdate) (*domain.Game, error) {
	game := &domain.Game{
		Name:  gameEdit.Title,
		ExtID: gameEdit.Ext,
	}

	if valid, message := game.ValidateGameTitle(); !valid {
		return nil, e.NewValidationError(message)
	}

	game, err := s.gameRepo.UpdateByExt(ctx, game)
	if err != nil {
		return nil, e.NewInternalError("Невозможно обновить игру", err)
	}

	game, err = s.gameRepo.GetByExt(ctx, game.ExtID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных игры", err)
	}

	return game, nil
}
