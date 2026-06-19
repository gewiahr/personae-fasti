package service

import (
	"context"
	"fmt"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/repo"
	gu "personae-fasti/pkg/gewi-utils"
)

type RecordService struct {
	recordRepo repo.RecordRepository
	playerRepo repo.PlayerRepository
	gameRepo   repo.GameRepository
}

func NewRecordService(playerRepo repo.PlayerRepository, gameRepo repo.GameRepository, recordRepo repo.RecordRepository) *RecordService {
	return &RecordService{
		recordRepo: recordRepo,
		playerRepo: playerRepo,
		gameRepo:   gameRepo,
	}
}

// Get records in current game available for player
func (s *RecordService) GetPlayerCurrentGameRecords(ctx context.Context, player *domain.Player) ([]domain.Record, error) {
	records, err := s.recordRepo.GetCurrentGameRecordList(ctx, player.CurrentGameID, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения записей", err)
	}
	return records, nil
}

// Post record in player current game
func (s *RecordService) PostPlayerCurrentGameRecord(ctx context.Context, player *domain.Player, recordInsert *dto.RecordInsert) (*domain.Record, error) {
	record := &domain.Record{
		Text:     recordInsert.Text,
		PlayerID: player.ID,
		GameID:   player.CurrentGameID,
		QuestID:  recordInsert.QuestID,
		HiddenBy: gu.TernaryInt(recordInsert.Hidden, player.ID, 0),
	}

	record, err := s.recordRepo.PostRecord(ctx, record)
	if err != nil {
		return nil, e.NewInternalError("Ошибка записи", err)
	}
	return record, nil
}

// Post record in player current game
func (s *RecordService) EditPlayerCurrentGameRecord(ctx context.Context, player *domain.Player, recordUpdate *dto.RecordUpdate) (*domain.Record, error) {
	record, err := s.recordRepo.EditRecord(ctx, recordUpdate, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка записи", err)
	}
	return record, nil
}

// Delete record in player current game
func (s *RecordService) DeletePlayerCurrentGameRecord(ctx context.Context, player *domain.Player, recordID int) error {
	record, err := s.recordRepo.GetRecord(ctx, player.ID, recordID)
	if err != nil {
		return e.NewInternalError("Ошибка получения записи", err)
	}

	if player.ID != record.PlayerID && player.ID != player.CurrentGame.GMID {
		if !player.CurrentGame.Settings.AllowAllEditRecords {
			return e.NewForbiddenError(fmt.Sprintf("player %s cannot delete other players' records", player.Username))
		}
	}

	err = s.recordRepo.SoftDeleteRecord(ctx, player.ID, recordID)
	if err != nil {
		return e.NewInternalError("Ошибка записи", err)
	}
	return nil
}
