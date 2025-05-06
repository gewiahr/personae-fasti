package data

import (
	"context"
	"database/sql"
	"fmt"
	"personae-fasti/api/models/reqData"
	"time"

	"github.com/uptrace/bun"
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

func (s *Storage) GetCurrentGameRecordsForPlayer(game *Game, player *Player) ([]Record, error) {
	var records []Record
	err := s.db.NewSelect().Model(&records).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("game_id = ?", game.ID)
		}).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("hidden_by = 0").WhereOr("hidden_by = ?", player.ID)
		}).
		Scan(context.Background(), &records)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return []Record{}, nil
	}

	return records, nil

	// === Old implementation without hidden records === //
	//
	// err := s.db.NewSelect().Model(game).WherePK().Relation("Records").Scan(context.Background())
	// if err != nil {
	// 	return nil, err
	// } else if err == sql.ErrNoRows || game.Records == nil {
	// 	return []Record{}, nil
	// }
	//
	// return game.Records, nil
}

func (s *Storage) GetCurrentGameSessions(game *Game) ([]Session, error) {
	err := s.db.NewSelect().Model(game).WherePK().Relation("Sessions").Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows || game.Sessions == nil {
		return []Session{}, nil
	}

	return game.Sessions, nil
}

func (s *Storage) InsertNewRecord(recordInsert *reqData.RecordInsert, p *Player) error {
	record := Record{
		Text:     recordInsert.Text,
		PlayerID: p.ID,
		GameID:   p.CurrentGameID,
	}

	if recordInsert.Hidden {
		record.HiddenBy = p.ID
	}

	err := s.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		// Insert Record
		result, err := s.db.NewInsert().Model(&record).Exec(context.Background())
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("empty insert")
		}
		// Insert Mentions
		if err := s.InsertMentionsForRecord(&record); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateRecord(recordUpdate *reqData.RecordUpdate, p *Player) error {
	now := time.Now().UTC()
	record := Record{
		ID:      recordUpdate.ID,
		Text:    recordUpdate.Text,
		Updated: &now,
	}

	if recordUpdate.Hidden {
		record.HiddenBy = p.ID
	} else {
		record.HiddenBy = 0
	}

	err := s.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		// Update Record
		result, err := s.db.NewUpdate().Model(&record).Column("text", "updated", "hidden_by").WherePK().Exec(context.Background())
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("empty insert")
		}

		// Delete Old Mentions
		if err := s.DeleteMentionsForRecord(&record); err != nil {
			return err
		}

		// Insert Mentions
		if err := s.InsertMentionsForRecord(&record); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
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

func (s *Storage) CreateChar(charCreate *reqData.CharCreate, player *Player) (*Char, error) {
	char := Char{
		Name:        charCreate.Name,
		Title:       charCreate.Title,
		Description: charCreate.Description,
		PlayerID:    player.ID,
		GameID:      player.CurrentGameID,
	}

	_, err := s.db.NewInsert().Model(&char).
		Column("name", "title", "description", "player_id", "game_id").
		Returning("*").Exec(context.Background(), &char)
	//Exec(context.Background())

	return &char, err
}

func (s *Storage) UpdateChar(charUpdate *reqData.CharUpdate, char *Char) (*Char, error) {
	_, err := s.db.NewUpdate().Model(char).WherePK().
		Set("name = ?", charUpdate.Name).
		Set("title = ?", charUpdate.Title).
		Set("description = ?", charUpdate.Description).
		Returning("*").Exec(context.Background())
	return char, err
}

func (s *Storage) GetCurrentGameNPCs(game *Game) ([]NPC, error) {
	err := s.db.NewSelect().Model(game).WherePK().Relation("NPCs").Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows || game.NPCs == nil {
		return []NPC{}, nil
	}

	return game.NPCs, nil
}

func (s *Storage) GetNPCByID(npcID int) (*NPC, error) {
	npc := NPC{
		ID: npcID,
	}

	err := s.db.NewSelect().Model(&npc).WherePK().Relation("Records").Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &npc, nil
}

func (s *Storage) CreateNPC(npcCreate *reqData.NPCCreate, player *Player) (*NPC, error) {
	npc := NPC{
		Name:        npcCreate.Name,
		Title:       npcCreate.Title,
		Description: npcCreate.Description,
		CreatedByID: player.ID,
		GameID:      player.CurrentGameID,
	}

	_, err := s.db.NewInsert().Model(&npc).
		Column("name", "title", "description", "created_by_id", "game_id").
		Returning("*").Exec(context.Background(), &npc)

	return &npc, err
}

func (s *Storage) UpdateNPC(npcUpdate *reqData.NPCUpdate, npc *NPC) (*NPC, error) {
	_, err := s.db.NewUpdate().Model(npc).WherePK().
		Set("name = ?", npcUpdate.Name).
		Set("title = ?", npcUpdate.Title).
		Set("description = ?", npcUpdate.Description).
		Returning("*").Exec(context.Background())
	return npc, err
}

func (s *Storage) GetCurrentGameLocations(game *Game) ([]Location, error) {
	err := s.db.NewSelect().Model(game).WherePK().Relation("Locations").Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows || game.Locations == nil {
		return []Location{}, nil
	}

	return game.Locations, nil
}

func (s *Storage) GetLocationByID(locationID int) (*Location, error) {
	location := Location{
		ID: locationID,
	}

	err := s.db.NewSelect().Model(&location).WherePK().Relation("Records").Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &location, nil
}

func (s *Storage) CreateLocation(locationCreate *reqData.LocationCreate, player *Player) (*Location, error) {
	location := Location{
		Name:        locationCreate.Name,
		Title:       locationCreate.Title,
		Description: locationCreate.Description,
		CreatedByID: player.ID,
		GameID:      player.CurrentGameID,
	}

	_, err := s.db.NewInsert().Model(&location).
		Column("name", "title", "description", "created_by_id", "game_id").
		Returning("*").Exec(context.Background(), &location)

	return &location, err
}

func (s *Storage) UpdateLocation(locationUpdate *reqData.LocationUpdate, location *Location) (*Location, error) {
	_, err := s.db.NewUpdate().Model(location).WherePK().
		Set("name = ?", locationUpdate.Name).
		Set("title = ?", locationUpdate.Title).
		Set("description = ?", locationUpdate.Description).
		Returning("*").Exec(context.Background())
	return location, err
}

func (s *Storage) GetSuggestions(player *Player) ([]Suggestion, error) {
	var suggestions []Suggestion

	err := s.db.NewRaw(
		`SELECT 
			id,
			CONCAT('char:', id) as sid,
			'char' as type,
			name
		FROM char
		WHERE game_id = ?

		UNION ALL

		SELECT
			id,
			CONCAT('npc:', id) as sid,
			'npc' as type,
			name
		FROM npc
		WHERE game_id = ?

		UNION ALL

		SELECT
			id,
			CONCAT('location:', id) as sid,
			'location' as type,
			name
		FROM location
		WHERE game_id = ?`,
		player.CurrentGameID, player.CurrentGameID, player.CurrentGameID,
	).Scan(context.Background(), &suggestions)

	if suggestions == nil {
		suggestions = []Suggestion{}
	}

	return suggestions, err
}

func (s *Storage) GetPlayerGames(player *Player) ([]Game, error) {
	err := s.db.NewSelect().Model(player).WherePK().Relation("Games").Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return player.Games, nil
}

func (s *Storage) ChangeCurrentGame(player *Player, gameID int) (*Game, error) {
	player.CurrentGameID = gameID
	_, err := s.db.NewUpdate().Model(player).Column("current_game_id").WherePK().Returning("*").Exec(context.Background())
	if err != nil {
		return nil, err
	}
	// ** Get to know why RETURNING is not working here properly ** //
	err = s.db.NewSelect().Model(player).WherePK().Relation("CurrentGame").Scan(context.Background(), player)
	if err != nil {
		return nil, err
	}
	// ** Get to know why RETURNING is not working here properly ** //
	return player.CurrentGame, nil
}
