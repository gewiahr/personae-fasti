package data

import (
	"context"
	"database/sql"
	"fmt"
	"personae-fasti/api/models/reqData"
)

func (s *Storage) GetPlayerByAccessKey(accesskey string) (*Player, error) {
	var player Player

	err := s.db.NewSelect().Model(&player).Where("accesskey = ?", accesskey).Relation("CurrentGame").Scan(context.Background(), &player)
	if err != nil {
		return nil, err
	}

	return &player, nil
}

func (s *Storage) GetCurrentGamePlayers(game *Game) ([]Player, error) {
	err := s.db.NewSelect().Model(game).WherePK().Relation("Players").Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return game.Players, nil
}

func (s *Storage) GetCurrentGameRecords(game *Game) ([]Record, error) {
	err := s.db.NewSelect().Model(game).WherePK().Relation("Records").Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return []Record{}, nil
	}

	return game.Records, nil
}

func (s *Storage) InsertNewRecord(recordInsert *reqData.RecordInsert, p *Player) error {
	record := Record{
		Text:     recordInsert.Text,
		PlayerID: p.ID,
		GameID:   p.CurrentGameID,
	}

	result, err := s.db.NewInsert().Model(&record).Exec(context.Background())
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("error")
	}
	return nil
}

func (s *Storage) GetCurrentGameChars(game *Game) ([]Char, error) {
	err := s.db.NewSelect().Model(game).WherePK().Relation("Chars").Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows || game.Chars == nil {
		return []Char{}, nil
	}

	return game.Chars, nil
}

func (s *Storage) GetCharByID(charID int) (*Char, error) {
	char := Char{
		ID: charID,
	}

	err := s.db.NewSelect().Model(&char).WherePK().Relation("Records").Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &char, nil
}

func (s *Storage) UpdateChar(charUpdate *reqData.CharUpdate, char *Char) (*Char, error) {
	_, err := s.db.NewUpdate().Model(char).WherePK().
		Set("name = ?", charUpdate.Name).
		Set("title = ?", charUpdate.Title).
		Set("description = ?", charUpdate.Description).
		Returning("*").Exec(context.Background())
	return char, err
}
