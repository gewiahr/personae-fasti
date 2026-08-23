package service

import (
	"context"
	"database/sql"
	"fmt"
	"personae-fasti/internal/domain"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/repo"
)

type PlayerService struct {
	playerRepo repo.PlayerRepository
	gameRepo   repo.GameRepository
}

func NewPlayerService(playerRepo repo.PlayerRepository, gameRepo repo.GameRepository) *PlayerService {
	return &PlayerService{
		playerRepo: playerRepo,
		gameRepo:   gameRepo,
	}
}

func (s *PlayerService) GetPlayerCurrentGame(ctx context.Context, player *domain.Player) (*domain.Game, error) {
	game, err := s.gameRepo.GetCurrentGame(ctx, player.CurrentGameID)
	if err == sql.ErrNoRows {
		return nil, e.NewNotFoundError("Нет текущей игры")
	} else if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных текущей игры", err)
	}
	return game, nil
}

func (s *PlayerService) ChangePlayerCurrentGame(ctx context.Context, player *domain.Player, gameExt string) (*domain.Game, error) {
	p, err := s.playerRepo.GetPlayerWithGames(ctx, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных игрока", err)
	}

	var game *domain.Game
	for _, playerGame := range p.Games {
		if playerGame.ExtID == gameExt {
			player, err = s.playerRepo.ChangeCurrentGame(ctx, player.ID, playerGame.ID)
			if err != nil {
				return nil, e.NewInternalError("Ошибка изменения текущей игры", err)
			}

			game, err = s.GetPlayerCurrentGame(ctx, player)
			if err != nil {
				return nil, e.NewInternalError("Ошибка получения данных игрока", err)
			}

			break
		}
	}

	if game == nil {
		return nil, e.NewForbiddenError(fmt.Sprintf("Игрок не состоит в игре %s", gameExt))
	}

	// TODO: transaction
	game, err = s.GetPlayerCurrentGame(ctx, player)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных текущей игры", err)
	}

	return game, nil
}

func (s *PlayerService) GetPlayerSettings(ctx context.Context, player *domain.Player) (*domain.Player, error) {
	p, err := s.playerRepo.GetPlayerWithGames(ctx, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных игрока", err)
	}

	return p, nil
}

func (s *PlayerService) GetPlayerInvites(ctx context.Context, player *domain.Player) ([]domain.GameInvite, error) {
	invites, err := s.gameRepo.GetPlayerInvites(ctx, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных игрока", err)
	}

	return invites, nil
}

func (s *PlayerService) AcceptGameInvite(ctx context.Context, player *domain.Player, inviteCode string) error {
	invite, err := s.playerRepo.GetInvite(ctx, player.ID, inviteCode)
	if err != nil {
		return e.NewInternalError("unable to get invite", err)
	}
	if invite == nil {
		return e.NewNotFoundError("invite does not exist")
	}

	if err := s.playerRepo.AddPlayerToGame(ctx, invite.PlayerID, invite.GameID); err != nil {
		return e.NewInternalError("cannot add player to game", err)
	}

	return nil
}

func (s *PlayerService) RefuseGameInvite(ctx context.Context, player *domain.Player, inviteCode string) error {
	invite, err := s.playerRepo.GetInvite(ctx, player.ID, inviteCode)
	if err != nil {
		return e.NewInternalError("unable to get invite", err)
	}
	if invite == nil {
		return e.NewNotFoundError("invite does not exist")
	}

	if err := s.gameRepo.DeleteInvite(ctx, invite); err != nil {
		return e.NewInternalError("unable to delete invite", err)
	}

	return nil
}
