package service

import (
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
