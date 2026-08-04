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

func (s *GameService) StartNewGameSession(ctx context.Context, player *domain.Player) (*domain.Session, error) {
	if player.CurrentGame.GMID != player.ID {
		return nil, e.NewForbiddenError("only GM may start new session")
	}

	newSession, err := s.gameRepo.CreateNewSession(ctx, player.CurrentGame)
	if err != nil {
		return nil, e.NewInternalError("fail to create new session", err)
	}

	return newSession, nil
}

func (s *GameService) EditGameSession(ctx context.Context, player *domain.Player, sessionUpdate *dto.SessionUpdate) (*domain.Session, error) {
	if player.CurrentGame.GMID != player.ID {
		return nil, e.NewForbiddenError("only GM may edit session")
	}

	session := &domain.Session{
		Name: sessionUpdate.Name,
	}

	if valid, message := session.ValidateSessionTitle(); !valid {
		return nil, e.NewValidationError(message)
	}

	game, err := s.gameRepo.EditSession(ctx, player.CurrentGame, sessionUpdate)
	if err != nil {
		return nil, e.NewInternalError("Невозможно создать новую игру", err)
	}

	return game, nil
}

func (s *GameService) RemoveLastGameSession(ctx context.Context, player *domain.Player) error {
	if player.CurrentGame.GMID != player.ID {
		return e.NewForbiddenError("only GM may delete session")
	}

	err := s.gameRepo.RemoveLastSession(ctx, player.CurrentGame)
	if err != nil {
		return e.NewInternalError("Невозможно удалить сессию", err)
	}

	return nil
}

func (s *GameService) InvitePlayer(ctx context.Context, player *domain.Player, invited string) error {
	if player.CurrentGame.GMID != player.ID {
		return e.NewForbiddenError("only GM may invite players")
	}

	playerInvited, err := s.playerRepo.GetByUsername(ctx, invited)
	if err != nil {
		return e.NewInternalError("error getting players", err)
	}
	if playerInvited == nil {
		return e.NewNotFoundError("Нет игрока с таким именем")
	}

	game, err := s.gameRepo.GetCurrentGame(ctx, player.CurrentGameID)
	if err != nil {
		return e.NewInternalError("error getting current game", err)
	}

	for _, inv := range game.Invites {
		if inv.ID == playerInvited.ID {
			return e.NewValidationError("Игрок уже приглашён")
		}
	}

	for _, pl := range game.Players {
		if pl.ID == playerInvited.ID {
			return e.NewValidationError("Игрок уже участвует в игре")
		}
	}

	err = s.gameRepo.InvitePlayer(ctx, &domain.GameInvite{GameID: player.CurrentGame.ID, PlayerID: playerInvited.ID})
	if err != nil {
		return e.NewInternalError("error inviting player", err)
	}

	return nil
}

func (s *GameService) RemoveInvite(ctx context.Context, player *domain.Player, invited string) error {
	if player.CurrentGame.GMID != player.ID {
		return e.NewForbiddenError("only GM may delete invites")
	}

	playerInvited, err := s.playerRepo.GetByUsername(ctx, invited)
	if err != nil {
		return e.NewInternalError("error getting players", err)
	}
	if playerInvited == nil {
		return e.NewNotFoundError("no player with such username")
	}

	err = s.gameRepo.DeleteInvite(ctx, &domain.GameInvite{GameID: player.CurrentGame.ID, PlayerID: playerInvited.ID})
	if err != nil {
		return e.NewInternalError("error removing invite", err)
	}

	return nil
}

func (s *GameService) UpdateGameSettings(ctx context.Context, settingsUpdate *dto.GameSettingsUpdate) (*domain.Game, error) {
	currentGame, err := s.gameRepo.GetByExt(ctx, settingsUpdate.GameExt)
	if err != nil {
		return nil, e.NewInternalError("error getting game settings", err)
	}

	gameSettings, err := s.gameRepo.UpdateGameSettings(ctx, currentGame.ID, settingsUpdate)
	if err != nil {
		return nil, e.NewInternalError("error updating settings", err)
	}

	currentGame.Settings = gameSettings
	return currentGame, nil
}
