package data

// import (
// 	"context"
// 	"crypto/sha256"
// 	"database/sql"
// 	"encoding/hex"
// 	"errors"
// 	"fmt"
// 	"personae-fasti/api/models/reqData"
// 	gu "personae-fasti/pkg/gewi-utils"
// 	"strings"
// 	"time"

// 	tgInitData "github.com/telegram-mini-apps/init-data-golang"
// 	"github.com/uptrace/bun"
// )

// func (s *Storage) GetPlayerByAccessKey(accesskey string) (*Player, error) {
// 	var player Player

// 	if accesskey == "" {
// 		return nil, fmt.Errorf("accesskey cannot be empty")
// 	}

// 	err := s.db.NewSelect().Model(&player).Where("accesskey = ?", accesskey).Relation("CurrentGame.Settings").Relation("CurrentGame.Sessions").Scan(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &player, nil
// }

// func (s *Storage) GetPlayerByUsername(username string) (*Player, error) {
// 	if username == "" {
// 		return nil, fmt.Errorf("username cannot be empty")
// 	}

// 	var player Player

// 	if err := s.db.NewSelect().Model(&player).Where("username = ?", username).Relation("RegData").Relation("CurrentGame.Settings").Relation("CurrentGame.Sessions").Scan(context.Background(), &player); err == sql.ErrNoRows {
// 		return nil, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return &player, nil
// }

// func (s *Storage) GetPlayerByTGToken(tokenString string) (*Player, error) {
// 	//var player Player

// 	if tokenString == "" {
// 		return nil, fmt.Errorf("token cannot be empty")
// 	}

// 	tokenHash := sha256.Sum256([]byte(tokenString))
// 	tokenHashHex := hex.EncodeToString(tokenHash[:])

// 	var token Token
// 	err := s.db.NewSelect().Model(&token).Where("token_hash = ?", tokenHashHex).Relation("Player").Relation("Player.RegData").Relation("Player.CurrentGame.Settings").Relation("Player.CurrentGame.Sessions").Scan(context.Background())

// 	//err := s.db.NewSelect().Model(&player).Where("accesskey = ?", accesskey).Relation("CurrentGame.Settings").Relation("CurrentGame.Sessions").Scan(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}

// 	return token.Player, nil
// }

// func (s *Storage) CreateAuthToken(player *Player, jwtSecret string, jwtTime time.Duration) (string, error) {
// 	expirationTime := time.Now().Add(jwtTime)

// 	tokenString := s.GeneratePlayerToken()
// 	// tokenString, err := s.GeneratePlayerJWTToken(player, jwtSecret, expirationTime)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	tokenHash := sha256.Sum256([]byte(tokenString))
// 	tokenHashHex := hex.EncodeToString(tokenHash[:])

// 	dbToken := &Token{
// 		PlayerID:  player.ID,
// 		TokenHash: tokenHashHex,
// 		ExpiresAt: expirationTime,
// 		Revoked:   false,
// 	}

// 	_, err := s.db.NewInsert().Model(dbToken).Returning("*").Exec(context.Background(), dbToken)
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// func (s *Storage) GetTelegramPlayer(tgID int64) (*Player, error) {
// 	var player Player

// 	err := s.db.NewSelect().Model(&player).Where("telegram_id = ?", tgID).Relation("Telegram").Relation("RegData").Relation("CurrentGame.Settings").Relation("CurrentGame.Sessions").Scan(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &player, nil
// }

// func (s *Storage) CreateTelegramPlayer(data tgInitData.InitData) (*Player, error) {
// 	var player *Player
// 	ctx := context.Background()

// 	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
// 		telegram := &Telegram{
// 			ID:       data.User.ID,
// 			Username: data.User.Username,
// 			Lang:     data.User.LanguageCode,
// 			PicURL:   data.User.PhotoURL,
// 		}
// 		_, err := tx.NewInsert().Model(telegram).Exec(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		player = &Player{
// 			Username:   fmt.Sprintf("tguser_%d", data.AuthDate().Unix()),
// 			TelegramID: telegram.ID,
// 		}
// 		_, err = tx.NewInsert().Model(player).Exec(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		playerRegData := &PlayerRegData{
// 			PlayerID:    player.ID,
// 			UsernameSet: false,
// 		}
// 		_, err = tx.NewInsert().Model(playerRegData).Exec(ctx)

// 		return nil
// 	})

// 	return player, err
// }

// func (s *Storage) GetCurrentGamePlayers(game *Game) ([]Player, error) {
// 	err := s.db.NewSelect().Model(game).WherePK().Relation("Players").Scan(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}

// 	return game.Players, nil
// }

// func (s *Storage) GetCurrentGameRecordsForPlayer(game *Game, player *Player) ([]Record, error) {
// 	if game == nil {
// 		return []Record{}, nil
// 	}

// 	var records []Record
// 	if err := s.db.NewSelect().Model(&records).
// 		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
// 			return q.Where("record.game_id = ?", game.ID)
// 		}).
// 		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
// 			return q.Where("record.hidden_by = 0").WhereOr("record.hidden_by = ?", player.ID)
// 		}).
// 		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
// 			return q.Where("record.deleted IS NULL")
// 		}).
// 		Relation("Quest").
// 		Scan(context.Background(), &records); err == sql.ErrNoRows {
// 		return []Record{}, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return records, nil
// }

// func (s *Storage) InsertNewRecord(recordInsert *reqData.RecordInsert, p *Player) error {
// 	record := Record{
// 		Text:     recordInsert.Text,
// 		PlayerID: p.ID,
// 		GameID:   p.CurrentGameID,
// 		QuestID:  recordInsert.QuestID,
// 		HiddenBy: gu.TernaryInt(recordInsert.Hidden, p.ID, 0),
// 	}

// 	err := s.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
// 		// Insert Record
// 		result, err := tx.NewInsert().Model(&record).Exec(context.Background())
// 		if err != nil {
// 			return err
// 		}
// 		if result == nil {
// 			return fmt.Errorf("empty insert")
// 		}
// 		// Insert Mentions
// 		if err := s.InsertMentionsForRecord(tx, &record); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *Storage) UpdateRecord(recordUpdate *reqData.RecordUpdate, p *Player) error {
// 	var oldRecord = Record{ID: recordUpdate.ID}
// 	err := s.db.NewSelect().Model(&oldRecord).WherePK().Scan(context.Background(), &oldRecord)
// 	if err != nil {
// 		return err
// 	}

// 	if p.ID != oldRecord.PlayerID && p.ID != p.CurrentGame.GMID {
// 		if !p.CurrentGame.Settings.AllowAllEditRecords {
// 			return fmt.Errorf("player %s cannot edit other players' records", p.Username)
// 		}
// 	}

// 	now := time.Now().UTC()
// 	record := Record{
// 		ID:       recordUpdate.ID,
// 		Text:     recordUpdate.Text,
// 		Updated:  &now,
// 		QuestID:  recordUpdate.QuestID,
// 		HiddenBy: gu.TernaryInt(recordUpdate.Hidden, p.ID, 0),
// 	}

// 	err = s.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
// 		// Update Record
// 		result, err := tx.NewUpdate().Model(&record).Column("text", "updated", "hidden_by", "quest_id").WherePK().Exec(context.Background())
// 		if err != nil {
// 			return err
// 		}
// 		if result == nil {
// 			return fmt.Errorf("empty insert")
// 		}

// 		// Delete Old Mentions
// 		if err := s.DeleteMentionsForRecord(tx, &record); err != nil {
// 			return err
// 		}

// 		// Insert Mentions
// 		if err := s.InsertMentionsForRecord(tx, &record); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *Storage) DeleteRecord(recordID int, p *Player) error {
// 	var oldRecord = Record{ID: recordID}
// 	err := s.db.NewSelect().Model(&oldRecord).WherePK().Scan(context.Background(), &oldRecord)
// 	if err != nil {
// 		return err
// 	}

// 	if p.ID != oldRecord.PlayerID && p.ID != p.CurrentGame.GMID {
// 		if !p.CurrentGame.Settings.AllowAllEditRecords {
// 			return fmt.Errorf("player %s cannot delete other players' records", p.Username)
// 		}
// 	}

// 	now := time.Now().UTC()
// 	record := Record{
// 		ID:      recordID,
// 		Deleted: &now,
// 	}

// 	// Delete Record
// 	result, err := s.db.NewUpdate().Model(&record).Column("deleted").WherePK().Exec(context.Background())
// 	if err != nil {
// 		return err
// 	}
// 	if result == nil {
// 		return fmt.Errorf("empty delete")
// 	}

// 	return nil
// }

// func (s *Storage) GetCurrentGameNPCs(game *Game) ([]NPC, error) {
// 	if err := s.db.NewSelect().Model(game).WherePK().Relation("NPCs").Scan(context.Background()); err == sql.ErrNoRows || game.NPCs == nil {
// 		return []NPC{}, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return game.NPCs, nil
// }

// func (s *Storage) GetNPCByID(npcID int) (*NPC, error) {
// 	npc := NPC{
// 		ID: npcID,
// 	}

// 	if err := s.db.NewSelect().Model(&npc).WherePK().Relation("Records").Scan(context.Background()); err == sql.ErrNoRows {
// 		return nil, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return &npc, nil
// }

// func (s *Storage) CreateNPC(npcCreate *reqData.NPCCreate, player *Player) (*NPC, error) {
// 	if npcCreate.Name == "" {
// 		return nil, fmt.Errorf("name cannot be empty")
// 	}

// 	var hiddenBy = 0
// 	if npcCreate.Hidden {
// 		hiddenBy = player.ID
// 	}

// 	npc := NPC{
// 		Name:        npcCreate.Name,
// 		Title:       npcCreate.Title,
// 		Description: npcCreate.Description,
// 		HiddenBy:    hiddenBy,
// 		CreatedByID: player.ID,
// 		GameID:      player.CurrentGameID,
// 	}

// 	_, err := s.db.NewInsert().Model(&npc).
// 		Column("name", "title", "description", "hidden_by", "created_by_id", "game_id").
// 		Returning("*").Exec(context.Background(), &npc)

// 	return &npc, err
// }

// func (s *Storage) UpdateNPC(npcUpdate *reqData.NPCUpdate, npc *NPC, player *Player) (*NPC, error) {
// 	if npcUpdate.Name == "" {
// 		return nil, fmt.Errorf("name cannot be empty")
// 	}

// 	var hiddenBy = 0
// 	if npcUpdate.Hidden {
// 		hiddenBy = player.ID
// 	}

// 	_, err := s.db.NewUpdate().Model(npc).WherePK().
// 		Set("name = ?", npcUpdate.Name).
// 		Set("title = ?", npcUpdate.Title).
// 		Set("description = ?", npcUpdate.Description).
// 		Set("hidden_by = ?", hiddenBy).
// 		Returning("*").Exec(context.Background())
// 	return npc, err
// }

// func (s *Storage) GetCurrentGameLocations(game *Game) ([]Location, error) {
// 	if err := s.db.NewSelect().Model(game).WherePK().Relation("Locations").Scan(context.Background()); err == sql.ErrNoRows || game.Locations == nil {
// 		return []Location{}, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return game.Locations, nil
// }

// func (s *Storage) GetLocationChildren(location *Location) ([]Location, error) {
// 	var locations []Location

// 	if err := s.db.NewSelect().Model(&locations).Where("game_id = ? AND pid = ?", location.GameID, location.ID).Scan(context.Background()); err == sql.ErrNoRows || locations == nil {
// 		return []Location{}, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return locations, nil
// }

// func (s *Storage) GetLocationByID(locationID int) (*Location, error) {
// 	location := Location{
// 		ID: locationID,
// 	}

// 	if err := s.db.NewSelect().Model(&location).WherePK().Relation("Records").Relation("Parent").Scan(context.Background()); err == sql.ErrNoRows {
// 		return nil, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return &location, nil
// }

// func (s *Storage) CreateLocation(locationCreate *reqData.LocationCreate, player *Player) (*Location, error) {
// 	if locationCreate.Name == "" {
// 		return nil, fmt.Errorf("name cannot be empty")
// 	}

// 	location := Location{
// 		Name:        locationCreate.Name,
// 		Title:       locationCreate.Title,
// 		Description: locationCreate.Description,
// 		ParentID:    locationCreate.ParentID,
// 		HiddenBy:    gu.TernaryInt(locationCreate.Hidden, player.ID, 0),
// 		CreatedByID: player.ID,
// 		GameID:      player.CurrentGameID,
// 	}

// 	_, err := s.db.NewInsert().Model(&location).
// 		Column("name", "title", "description", "pid", "hidden_by", "created_by_id", "game_id").
// 		Returning("*").Exec(context.Background(), &location)

// 	return &location, err
// }

// func (s *Storage) UpdateLocation(locationUpdate *reqData.LocationUpdate, location *Location, player *Player) (*Location, error) {
// 	if locationUpdate.Name == "" {
// 		return nil, fmt.Errorf("name cannot be empty")
// 	}

// 	_, err := s.db.NewUpdate().Model(location).WherePK().
// 		Set("name = ?", locationUpdate.Name).
// 		Set("title = ?", locationUpdate.Title).
// 		Set("description = ?", locationUpdate.Description).
// 		Set("pid = ?", locationUpdate.ParentID).
// 		Set("hidden_by = ?", gu.TernaryInt(locationUpdate.Hidden, player.ID, 0)).
// 		Returning("*").Exec(context.Background())
// 	return location, err
// }

// func (s *Storage) GetPlayerGamesAndInvites(player *Player) (*Player, error) {
// 	err := s.db.NewSelect().Model(player).WherePK().Relation("Games").Relation("Invites").Scan(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}

// 	return player, nil
// }

// func (s *Storage) ChangeCurrentGame(player *Player, gameID int) (*Game, error) {
// 	player.CurrentGameID = gameID
// 	_, err := s.db.NewUpdate().Model(player).Column("current_game_id").WherePK().Returning("*").Exec(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}
// 	// ** Get to know why RETURNING is not working here properly ** //
// 	//err = s.db.NewSelect().Model(player).WherePK().Relation("CurrentGame").Scan(context.Background(), player)
// 	var currentGame Game
// 	err = s.db.NewSelect().Model(&currentGame).Where("id = ?", player.CurrentGameID).Relation("Settings").Scan(context.Background(), &currentGame)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// ** Get to know why RETURNING is not working here properly ** //
// 	return &currentGame, nil
// }

// func (s *Storage) GetGameByID(gameID int) (*Game, error) {
// 	game := Game{
// 		ID: gameID,
// 	}

// 	if err := s.db.NewSelect().Model(&game).WherePK().Relation("Sessions").Relation("Settings").Relation("Players").Relation("Invites").Scan(context.Background()); err == sql.ErrNoRows {
// 		return nil, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return &game, nil
// }

// func (s *Storage) GetPlayerGame(gameID, playerID int) (*Game, error) {
// 	game := Game{
// 		ID: gameID,
// 	}

// 	if err := s.db.NewSelect().Model(&game).Join("JOIN players_games AS pg ON pg.game_id = game.id").Where("game.id = ?", gameID).Where("pg.player_id = ?", playerID).Relation("Sessions").Relation("Settings").Relation("Players").Relation("Invites").Scan(context.Background()); err == sql.ErrNoRows {
// 		return nil, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return &game, nil
// }

// func (s *Storage) CreateGame(player *Player, newGameRequest *reqData.GameCreate) (*Game, error) {
// 	if newGameRequest.Title == "" {
// 		return nil, fmt.Errorf("game title cannot be empty")
// 	}
// 	var err error
// 	ctx := context.Background()
// 	newGame := Game{
// 		Name: newGameRequest.Title,
// 		GMID: player.ID,
// 	}

// 	if err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
// 		_, err := tx.NewInsert().Model(&newGame).ExcludeColumn("id").Returning("*").Exec(ctx, &newGame)
// 		if err != nil {
// 			return err
// 		}
// 		newGameSettings := GameSettings{
// 			GameID: newGame.ID,
// 		}
// 		_, err = tx.NewInsert().Model(&newGameSettings).Returning("*").Exec(ctx, &newGameSettings)
// 		if err != nil {
// 			return err
// 		}
// 		newPlayerGame := PlayerGame{
// 			PlayerID: player.ID,
// 			GameID:   newGame.ID,
// 		}
// 		_, err = tx.NewInsert().Model(&newPlayerGame).Returning("*").Exec(ctx, &newPlayerGame)
// 		if err != nil {
// 			return err
// 		}
// 		if player.CurrentGameID == 0 {
// 			_, err := s.ChangeCurrentGame(player, newGame.ID)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}); err != nil {
// 		return nil, err
// 	}

// 	s.db.NewSelect().Model(&newGame).Relation("Sessions").Relation("Settings").Relation("Players").Relation("Invites").WherePK().Exec(ctx, &newGame)

// 	return &newGame, err
// }

// func (s *Storage) UpdateGame(player *Player, updateGameRequest *reqData.GameUpdate) (*Game, error) {
// 	if updateGameRequest.Title == "" {
// 		return nil, fmt.Errorf("game title cannot be empty")
// 	}
// 	if updateGameRequest.GMID <= 0 {
// 		return nil, fmt.Errorf("gmID cannot be 0 or negative")
// 	}

// 	ctx := context.Background()

// 	updateGame := Game{
// 		ID: updateGameRequest.ID,
// 	}
// 	if _, err := s.db.NewSelect().Model(&updateGame).Relation("Sessions").Relation("Settings").Relation("Players").Relation("Invites").WherePK().Exec(ctx, &updateGame); err == sql.ErrNoRows {
// 		return nil, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	if updateGame.GMID != player.ID {
// 		return nil, fmt.Errorf("only GM may edit game")
// 	}

// 	updateGame.Name = updateGameRequest.Title
// 	updateGame.GMID = updateGameRequest.GMID

// 	if _, err := s.db.NewUpdate().Model(&updateGame).WherePK().
// 		Set("name = ?", updateGameRequest.Title).
// 		Set("gm_id = ?", updateGameRequest.GMID).
// 		Returning("*").Exec(context.Background()); err != nil {
// 		return nil, err
// 	}

// 	return &updateGame, nil
// }

// func (s *Storage) CheckUsernameAvailability(player *Player, usernameToCheck string) (bool, error) {
// 	count, err := s.db.NewSelect().Model(player).Where("username = ?", usernameToCheck).Count(context.Background())
// 	if err != nil {
// 		return false, err
// 	}

// 	return count == 0, nil
// }

// func (s *Storage) ChangeUsername(player *Player, newUsername string) (*Player, error) {
// 	available, err := s.CheckUsernameAvailability(player, newUsername)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !available {
// 		return nil, fmt.Errorf("username %s is not available", newUsername)
// 	}

// 	// ## Wrap in transaction ## //
// 	player.Username = newUsername
// 	_, err = s.db.NewUpdate().Model(player).Column("username").WherePK().Returning("*").Exec(context.Background(), player)
// 	if err != nil {
// 		return nil, err
// 	}

// 	player.RegData.UsernameSet = true
// 	_, err = s.db.NewUpdate().Model(player.RegData).Column("username_set").WherePK().Exec(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}
// 	// ## Wrap in transaction ## //

// 	return player, nil
// }

// func (s *Storage) AddServiceFeedback(p *Player, serviceFeedback *reqData.ServiceFeedback) error {
// 	feedback := ServiceFeedback{
// 		PlayerID: p.ID,
// 		GameID:   p.CurrentGameID,
// 		Type:     serviceFeedback.Type,
// 		Text:     serviceFeedback.Text,
// 	}

// 	_, err := s.db.NewInsert().Model(&feedback).Returning("*").Exec(context.Background(), &feedback)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
