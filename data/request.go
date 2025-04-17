package data

import (
	"context"
	"fmt"
)

func (s *Storage) GetPlayerByAccessKey(accesskey string) (*Player, error) {
	var player Player

	err := s.db.NewSelect().Model(&player).Where("accesskey = ?", accesskey).Relation("CurrentGame").Scan(context.Background(), &player)
	if err != nil {
		return nil, err
	}

	return &player, nil
}

func (s *Storage) GetCurrentGameRecords(game *Game) ([]Record, error) {
	err := s.db.NewSelect().Model(game).WherePK().Relation("Records").Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return game.Records, nil
}

func (s *Storage) GetCurrentGamePlayers(game *Game) ([]Player, error) {
	err := s.db.NewSelect().Model(game).WherePK().Relation("Players").Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return game.Players, nil
}

func (s *Storage) InsertNewRecord(record *Record, player *Player, game *Game) error {
	result, err := s.db.NewInsert().Model(record).Exec(context.Background())
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("error")
	}
	return nil
}
