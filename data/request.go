package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"personae-fasti/api/models/reqData"
	gu "personae-fasti/gewi-utils"
	"strings"
	"time"

	tgInitData "github.com/telegram-mini-apps/init-data-golang"
	"github.com/uptrace/bun"
)

func (s *Storage) GetPlayerByAccessKey(accesskey string) (*Player, error) {
	var player Player

	if accesskey == "" {
		return nil, fmt.Errorf("accesskey cannot be empty")
	}

	err := s.db.NewSelect().Model(&player).Where("accesskey = ?", accesskey).Relation("CurrentGame.Settings").Relation("CurrentGame.Sessions").Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return &player, nil
}

func (s *Storage) GetPlayerByTGToken(tokenString string) (*Player, error) {
	//var player Player

	if tokenString == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	tokenHash := sha256.Sum256([]byte(tokenString))
	tokenHashHex := hex.EncodeToString(tokenHash[:])

	var token Token
	err := s.db.NewSelect().Model(&token).Where("token_hash = ?", tokenHashHex).Relation("Player").Relation("Player.RegData").Relation("Player.CurrentGame.Settings").Relation("Player.CurrentGame.Sessions").Scan(context.Background())

	//err := s.db.NewSelect().Model(&player).Where("accesskey = ?", accesskey).Relation("CurrentGame.Settings").Relation("CurrentGame.Sessions").Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return token.Player, nil
}

func (s *Storage) CreateAuthToken(player *Player, jwtSecret string, jwtTime time.Duration) (string, error) {
	expirationTime := time.Now().Add(jwtTime)

	tokenString := s.GeneratePlayerToken(player, expirationTime)
	// tokenString, err := s.GeneratePlayerJWTToken(player, jwtSecret, expirationTime)
	// if err != nil {
	// 	return nil, err
	// }

	tokenHash := sha256.Sum256([]byte(tokenString))
	tokenHashHex := hex.EncodeToString(tokenHash[:])

	dbToken := &Token{
		PlayerID:  player.ID,
		TokenHash: tokenHashHex,
		ExpiresAt: expirationTime,
		Revoked:   false,
	}

	_, err := s.db.NewInsert().Model(dbToken).Returning("*").Exec(context.Background(), dbToken)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Storage) GetTelegramPlayer(tgID int64) (*Player, error) {
	var player Player

	err := s.db.NewSelect().Model(&player).Where("telegram_id = ?", tgID).Relation("Telegram").Relation("RegData").Relation("CurrentGame.Settings").Relation("CurrentGame.Sessions").Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return &player, nil
}

func (s *Storage) CreateTelegramPlayer(data tgInitData.InitData) (*Player, error) {
	var player *Player
	ctx := context.Background()

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		telegram := &Telegram{
			ID:       data.User.ID,
			Username: data.User.Username,
			Lang:     data.User.LanguageCode,
			PicURL:   data.User.PhotoURL,
		}
		_, err := s.db.NewInsert().Model(telegram).Exec(ctx)
		if err != nil {
			return err
		}

		player = &Player{
			Username:   fmt.Sprintf("tguser_%d", data.AuthDate().Unix()),
			AccessKey:  "",
			TelegramID: telegram.ID,
		}
		_, err = s.db.NewInsert().Model(player).Exec(ctx)
		if err != nil {
			return err
		}

		playerRegData := &PlayerRegData{
			PlayerID:    player.ID,
			UsernameSet: false,
		}
		_, err = s.db.NewInsert().Model(playerRegData).Exec(ctx)

		return nil
	})

	return player, err
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
			return q.Where("record.game_id = ?", game.ID)
		}).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("record.hidden_by = 0").WhereOr("record.hidden_by = ?", player.ID)
		}).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("record.deleted IS NULL")
		}).
		Relation("Quest").
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

func (s *Storage) GetCurrentGameSession(game *Game) (*Session, error) {
	var currentSession Session

	err := s.db.NewSelect().Model(&currentSession).Where("game_id = ? AND end_time IS NULL", game.ID).Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &currentSession, nil
}

func (s *Storage) StartNewGameSession(game *Game) (*Session, error) {
	currentSession, err := s.GetCurrentGameSession(game)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	err = s.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		sessionNumber := 0
		currentTime := time.Now().UTC()

		// Start new session
		if currentSession != nil {
			currentSession.EndTime = &currentTime

			_, err := s.db.NewUpdate().Model(currentSession).Column("end_time").WherePK().Exec(context.Background())
			if err != nil {
				return err
			} else if err == sql.ErrNoRows {
				return fmt.Errorf("cannot update previous session row")
			}

			sessionNumber = currentSession.Number + 1

			// Start first session
		} else {
			sessionZero := &Session{
				GameID:  game.ID,
				Number:  sessionNumber,
				EndTime: &currentTime,
			}

			_, err = s.db.NewInsert().Model(sessionZero).Exec(context.Background())
			if err != nil {
				return err
			}

			sessionNumber++
		}

		newSession := &Session{
			GameID: game.ID,
			Number: sessionNumber,
		}

		_, err = s.db.NewInsert().Model(newSession).Exec(context.Background())
		if err != nil {
			return err
		} else if err == sql.ErrNoRows {
			return fmt.Errorf("cannot create new session row")
		}

		return nil
	})

	return currentSession, nil
}

func (s *Storage) InsertNewRecord(recordInsert *reqData.RecordInsert, p *Player) error {
	record := Record{
		Text:     recordInsert.Text,
		PlayerID: p.ID,
		GameID:   p.CurrentGameID,
		QuestID:  recordInsert.QuestID,
		HiddenBy: gu.TernaryInt(recordInsert.Hidden, p.ID, 0),
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
	var oldRecord = Record{ID: recordUpdate.ID}
	err := s.db.NewSelect().Model(&oldRecord).WherePK().Scan(context.Background(), &oldRecord)
	if err != nil {
		return err
	}

	if p.ID != oldRecord.PlayerID && p.ID != p.CurrentGame.GMID {
		if !p.CurrentGame.Settings.AllowAllEditRecords {
			return fmt.Errorf("player %s cannot edit other players' records", p.Username)
		}
	}

	now := time.Now().UTC()
	record := Record{
		ID:       recordUpdate.ID,
		Text:     recordUpdate.Text,
		Updated:  &now,
		QuestID:  recordUpdate.QuestID,
		HiddenBy: gu.TernaryInt(recordUpdate.Hidden, p.ID, 0),
	}

	err = s.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		// Update Record
		result, err := s.db.NewUpdate().Model(&record).Column("text", "updated", "hidden_by", "quest_id").WherePK().Exec(context.Background())
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

func (s *Storage) DeleteRecord(recordID int, p *Player) error {
	var oldRecord = Record{ID: recordID}
	err := s.db.NewSelect().Model(&oldRecord).WherePK().Scan(context.Background(), &oldRecord)
	if err != nil {
		return err
	}

	if p.ID != oldRecord.PlayerID && p.ID != p.CurrentGame.GMID {
		if !p.CurrentGame.Settings.AllowAllEditRecords {
			return fmt.Errorf("player %s cannot delete other players' records", p.Username)
		}
	}

	now := time.Now().UTC()
	record := Record{
		ID:      recordID,
		Deleted: &now,
	}

	// Delete Record
	result, err := s.db.NewUpdate().Model(&record).Column("deleted").WherePK().Exec(context.Background())
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("empty delete")
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
	if charCreate.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	var hiddenBy = 0
	if charCreate.Hidden {
		hiddenBy = player.ID
	}

	char := Char{
		Name:        charCreate.Name,
		Title:       charCreate.Title,
		Description: charCreate.Description,
		HiddenBy:    hiddenBy,
		PlayerID:    player.ID,
		GameID:      player.CurrentGameID,
	}

	_, err := s.db.NewInsert().Model(&char).
		Column("name", "title", "description", "hidden_by", "player_id", "game_id").
		Returning("*").Exec(context.Background(), &char)
	//Exec(context.Background())

	return &char, err
}

func (s *Storage) UpdateChar(charUpdate *reqData.CharUpdate, char *Char, player *Player) (*Char, error) {
	if charUpdate.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	var hiddenBy = 0
	if charUpdate.Hidden {
		hiddenBy = player.ID
	}

	_, err := s.db.NewUpdate().Model(char).WherePK().
		Set("name = ?", charUpdate.Name).
		Set("title = ?", charUpdate.Title).
		Set("description = ?", charUpdate.Description).
		Set("hidden_by = ?", hiddenBy).
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
	if npcCreate.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	var hiddenBy = 0
	if npcCreate.Hidden {
		hiddenBy = player.ID
	}

	npc := NPC{
		Name:        npcCreate.Name,
		Title:       npcCreate.Title,
		Description: npcCreate.Description,
		HiddenBy:    hiddenBy,
		CreatedByID: player.ID,
		GameID:      player.CurrentGameID,
	}

	_, err := s.db.NewInsert().Model(&npc).
		Column("name", "title", "description", "hidden_by", "created_by_id", "game_id").
		Returning("*").Exec(context.Background(), &npc)

	return &npc, err
}

func (s *Storage) UpdateNPC(npcUpdate *reqData.NPCUpdate, npc *NPC, player *Player) (*NPC, error) {
	if npcUpdate.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	var hiddenBy = 0
	if npcUpdate.Hidden {
		hiddenBy = player.ID
	}

	_, err := s.db.NewUpdate().Model(npc).WherePK().
		Set("name = ?", npcUpdate.Name).
		Set("title = ?", npcUpdate.Title).
		Set("description = ?", npcUpdate.Description).
		Set("hidden_by = ?", hiddenBy).
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

func (s *Storage) GetLocationChildren(location *Location) ([]Location, error) {
	var locations []Location

	err := s.db.NewSelect().Model(&locations).Where("game_id = ? AND pid = ?", location.GameID, location.ID).Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows || locations == nil {
		return []Location{}, nil
	}

	return locations, nil
}

func (s *Storage) GetLocationByID(locationID int) (*Location, error) {
	location := Location{
		ID: locationID,
	}

	err := s.db.NewSelect().Model(&location).WherePK().Relation("Records").Relation("Parent").Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &location, nil
}

func (s *Storage) CreateLocation(locationCreate *reqData.LocationCreate, player *Player) (*Location, error) {
	if locationCreate.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	location := Location{
		Name:        locationCreate.Name,
		Title:       locationCreate.Title,
		Description: locationCreate.Description,
		ParentID:    locationCreate.ParentID,
		HiddenBy:    gu.TernaryInt(locationCreate.Hidden, player.ID, 0),
		CreatedByID: player.ID,
		GameID:      player.CurrentGameID,
	}

	_, err := s.db.NewInsert().Model(&location).
		Column("name", "title", "description", "pid", "hidden_by", "created_by_id", "game_id").
		Returning("*").Exec(context.Background(), &location)

	return &location, err
}

func (s *Storage) UpdateLocation(locationUpdate *reqData.LocationUpdate, location *Location, player *Player) (*Location, error) {
	if locationUpdate.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	_, err := s.db.NewUpdate().Model(location).WherePK().
		Set("name = ?", locationUpdate.Name).
		Set("title = ?", locationUpdate.Title).
		Set("description = ?", locationUpdate.Description).
		Set("pid = ?", locationUpdate.ParentID).
		Set("hidden_by = ?", gu.TernaryInt(locationUpdate.Hidden, player.ID, 0)).
		Returning("*").Exec(context.Background())
	return location, err
}

func (s *Storage) GetCurrentGameQuests(game *Game) ([]Quest, error) {
	var quests []Quest
	err := s.db.NewSelect().Model(&quests).Where("game_id = ? AND deleted is NULL", game.ID).Scan(context.Background())
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows || quests == nil {
		return quests, nil
	}

	return quests, nil
}

func (s *Storage) GetQuestByID(questID int) (*Quest, error) {
	quest := Quest{
		ID: questID,
	}

	err := s.db.NewSelect().Model(&quest).WherePK().Relation("Records").Relation("Tasks").Scan(context.Background())
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &quest, nil
}

func (s *Storage) CreateQuest(questCreate *reqData.QuestCreate, tasksCreate []reqData.TaskCreate, player *Player) (*Quest, error) {
	if questCreate.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	var quest *Quest
	ctx := context.Background()

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		quest = &Quest{
			Name:        questCreate.Name,
			Title:       questCreate.Title,
			Description: questCreate.Description,
			GameID:      player.CurrentGameID,
			ParentID:    questCreate.ParentID,
			ChildID:     questCreate.ChildID,
			HeadID:      questCreate.HeadID,
			Successful:  questCreate.Successful,
			HiddenBy:    gu.TernaryInt(questCreate.Hidden, player.ID, 0),
		}

		_, err := tx.NewInsert().Model(quest).
			Column("name", "title", "description", "game_id", "parent_id", "child_id", "head_id", "successful", "hidden_by").
			Returning("*").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to insert quest: %w", err)
		}

		if len(tasksCreate) > 0 {
			questTasks := make([]*QuestTask, len(tasksCreate))
			for i, taskCreate := range tasksCreate {
				questTasks[i] = &QuestTask{
					GameID:      player.CurrentGameID,
					QuestID:     quest.ID,
					Name:        taskCreate.Name,
					Description: taskCreate.Description,
					Type:        QuestTaskType(taskCreate.Type),
					Capacity:    taskCreate.Capacity,
					HiddenBy:    gu.TernaryInt(taskCreate.Hidden, player.ID, 0),
				}
			}

			_, err = tx.NewInsert().Model(&questTasks).Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to insert tasks: %w", err)
			}
		}

		err = tx.NewSelect().
			Model(quest).
			Relation("Tasks").
			Where("id = ?", quest.ID).
			Scan(ctx)
		if err != nil {
			return fmt.Errorf("failed to load quest with tasks: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return quest, nil
}

func (s *Storage) UpdateQuest(questUpdate *reqData.QuestUpdate, tasksUpdate []reqData.TaskUpdate, quest *Quest, player *Player) (*Quest, error) {
	if questUpdate.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	ctx := context.Background()
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := s.db.NewUpdate().Model(quest).WherePK().
			Set("name = ?", questUpdate.Name).
			Set("title = ?", questUpdate.Title).
			Set("description = ?", questUpdate.Description).
			Set("hidden_by = ?", gu.TernaryInt(questUpdate.Hidden, player.ID, 0)).
			Returning("*").Exec(context.Background()); err != nil {
			return fmt.Errorf("failed to update quest: %w", err)
		}

		if len(tasksUpdate) == 0 {

			if _, err := tx.NewDelete().
				Model((*QuestTask)(nil)).
				Where("quest_id = ?", quest.ID).
				Exec(ctx); err != nil {
				return fmt.Errorf("failed to update quest: %w", err)
			}

		} else {

			var values []any
			var valuePlaceholders []string

			for _, task := range tasksUpdate {
				hiddenBy := gu.TernaryInt(task.Hidden, player.ID, 0)
				values = append(values,
					task.ID,
					task.Name,
					task.Description,
					task.Type,
					task.Capacity,
					hiddenBy,
				)
			}

			for range tasksUpdate {
				valuePlaceholders = append(valuePlaceholders, "(?,?,?,?,?,?)")
			}

			query := fmt.Sprintf(`
				WITH input_data(id, name, description, type, capacity, hidden_by) AS (
					VALUES %s
				),
				updated AS (
					UPDATE quest_task t SET
						name = i.name,
						description = i.description,
						type = i.type,
						capacity = i.capacity,
						hidden_by = i.hidden_by
					FROM input_data i
					WHERE t.id = i.id AND t.quest_id = ?
					RETURNING t.id
				),
				inserted AS (
					INSERT INTO quest_task
						(quest_id, game_id, name, description, type, capacity, hidden_by)
					SELECT
						?, q.game_id, i.name, i.description, i.type, i.capacity, i.hidden_by
					FROM input_data i
					JOIN quest q ON q.id = ?
					WHERE i.id = 0
					RETURNING id
				),
				deleted AS (
					DELETE FROM quest_task
					WHERE quest_id = ?
					AND id NOT IN (SELECT id FROM input_data WHERE id != 0)
					RETURNING id
				)
				SELECT
					(SELECT COUNT(*) FROM updated) AS updated_count,
					(SELECT COUNT(*) FROM inserted) AS inserted_count,
					(SELECT COUNT(*) FROM deleted) AS deleted_count
			`,
				strings.Join(valuePlaceholders, ","))

			values = append(values, quest.ID, quest.ID, quest.ID, quest.ID)

			if _, err := tx.Exec(query, values...); err != nil {
				return fmt.Errorf("bulk task update failed: %w", err)
			}

		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	err = s.db.NewSelect().Model(quest).WherePK().Relation("Records").Relation("Tasks").Scan(context.Background())
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return quest, nil
}

func (s *Storage) DeleteQuest(questID int, p *Player) error {
	now := time.Now().UTC()
	quest := Quest{
		ID:      questID,
		Deleted: &now,
	}

	// Delete Quest
	result, err := s.db.NewUpdate().Model(&quest).Column("deleted").WherePK().Exec(context.Background())
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("empty delete")
	}

	return nil
}

func (s *Storage) GetTasksByQuest(quest *Quest) ([]QuestTask, error) {
	tasks := []QuestTask{}

	err := s.db.NewSelect().Model(&tasks).Where("quest_id = ?", quest.ID).Scan(context.Background())
	if err == sql.ErrNoRows {
		return tasks, nil
	} else if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) UpdateQuestTasks(tasksUpdate []reqData.TaskPatch, quest *Quest, player *Player) ([]QuestTask, error) {
	if len(tasksUpdate) == 0 || len(quest.Tasks) == 0 {
		return nil, errors.New("empty tasks on update or quest itself")
	}

	var tasks = quest.Tasks
	var finishTime = time.Now().UTC()
	for i := range tasks {
		for _, task := range tasksUpdate {
			if tasks[i].ID == task.ID {
				tasks[i].Current = task.Current
				switch tasks[i].Type {
				case Binary:
					if tasks[i].Current > 0 {
						tasks[i].Finished = &finishTime
					} else {
						tasks[i].Finished = nil
					}
				case Decimal:
					if tasks[i].Current >= tasks[i].Capacity {
						tasks[i].Finished = &finishTime
					} else {
						tasks[i].Finished = nil
					}
				}
			}
		}
	}

	_, err := s.db.NewUpdate().Model(&tasks).Column("current", "finished").Bulk().Returning("*").Exec(context.Background())
	return tasks, err
}

func (s *Storage) GetSuggestions(player *Player) ([]Suggestion, error) {
	var suggestions []Suggestion

	err := s.db.NewRaw(
		`SELECT 
			id,
			CONCAT('char:', id) as sid,
			'char' as type,
			name,
			CASE 
				WHEN hidden_by = 0 OR hidden_by = ? THEN false
				ELSE true
			END as hidden
		FROM char
		WHERE game_id = ?

		UNION ALL

		SELECT
			id,
			CONCAT('npc:', id) as sid,
			'npc' as type,
			name,
			CASE 
				WHEN hidden_by = 0 OR hidden_by = ? THEN false
				ELSE true
			END as hidden
		FROM npc
		WHERE game_id = ?

		UNION ALL

		SELECT
			id,
			CONCAT('location:', id) as sid,
			'location' as type,
			name,
			CASE 
				WHEN hidden_by = 0 OR hidden_by = ? THEN false
				ELSE true
			END as hidden
		FROM location
		WHERE game_id = ?`,
		player.ID, player.CurrentGameID, player.ID, player.CurrentGameID, player.ID, player.CurrentGameID,
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
	//err = s.db.NewSelect().Model(player).WherePK().Relation("CurrentGame").Scan(context.Background(), player)
	var currentGame Game
	err = s.db.NewSelect().Model(&currentGame).Where("id = ?", player.CurrentGameID).Relation("Settings").Scan(context.Background(), &currentGame)
	if err != nil {
		return nil, err
	}
	// ** Get to know why RETURNING is not working here properly ** //
	return &currentGame, nil
}

func (s *Storage) CreateGame(player *Player, newGameRequest *reqData.GameCreate) (*Game, error) {
	if newGameRequest.Name == "" {
		return nil, fmt.Errorf("game name cannot be empty")
	}
	ctx := context.Background()

	// ** Run in Transaction ** //
	newGame := Game{
		Name: newGameRequest.Name,
		GMID: player.ID,
	}
	_, err := s.db.NewInsert().Model(&newGame).ExcludeColumn("id").Returning("*").Exec(ctx, &newGame)
	if err != nil {
		return nil, err
	}
	newGameSettings := GameSettings{
		GameID: newGame.ID,
	}
	_, err = s.db.NewInsert().Model(&newGameSettings).Returning("*").Exec(ctx, &newGameSettings)
	if err != nil {
		return nil, err
	}
	// ** Run in Transaction ** //

	s.db.NewSelect().Model(&newGame).Relation("Settings").WherePK().Exec(ctx, &newGame)

	return &newGame, err
}

func (s *Storage) CheckUsernameAvailability(player *Player, usernameToCheck string) (bool, error) {
	count, err := s.db.NewSelect().Model(player).Where("username = ?", usernameToCheck).Count(context.Background())
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (s *Storage) ChangeUsername(player *Player, newUsername string) (*Player, error) {
	available, err := s.CheckUsernameAvailability(player, newUsername)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, fmt.Errorf("username %s is not available", newUsername)
	}

	// ## Wrap in transaction ## //
	player.Username = newUsername
	_, err = s.db.NewUpdate().Model(player).Column("username").WherePK().Returning("*").Exec(context.Background(), player)
	if err != nil {
		return nil, err
	}

	player.RegData.UsernameSet = true
	_, err = s.db.NewUpdate().Model(player.RegData).Column("username_set").WherePK().Exec(context.Background())
	if err != nil {
		return nil, err
	}
	// ## Wrap in transaction ## //

	return player, nil
}

func (s *Storage) UpdateGameSettings(gameSettingsUpdate *reqData.GameSettingsUpdate) (*Game, error) {
	gameSettings := GameSettings{
		GameID:              gameSettingsUpdate.GameID,
		AllowAllEditRecords: gameSettingsUpdate.AllowAllEditRecords,
	}

	_, err := s.db.NewUpdate().Model(&gameSettings).Column("allow_all_edit_records").WherePK().Returning("*").Exec(context.Background(), &gameSettings)
	if err != nil {
		return nil, err
	}

	var currentGame Game
	err = s.db.NewSelect().Model(&currentGame).Where("game_id = ?", gameSettings.GameID).Relation("Settings").Scan(context.Background(), &currentGame)
	if err != nil {
		return nil, err
	}

	return &currentGame, nil
}
