package data

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bun"
)

func (s *Storage) GeneratePlayerToken(player *Player, expTime time.Time) string {
	data := make([]byte, 16)
	rand.Read(data[0:16])

	return hex.EncodeToString(data)
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *Storage) GeneratePlayerJWTToken(player *Player, jwtSecret string, expTime time.Time) (string, error) {
	claims := &Claims{
		UserID:   int64(player.ID),
		Username: player.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   player.Username,
		},
	}

	var JWTBytes = []byte(jwtSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTBytes)
}

func (s *Storage) InsertMentionsForRecord(record *Record) error {
	re, err := regexp.Compile(`@(?P<type>\w+):(?P<id>\d+)` + "`(?P<name>[^`]+)`")
	if err != nil {
		return err
	}

	matches := re.FindAllStringSubmatch(record.Text, -1)
	for _, match := range matches {
		// Parse mention ID
		id, err := strconv.Atoi(match[2]) //ParseInt(match[2], 10, 64)
		if err != nil {
			continue
		}
		// Insert to a correct type
		switch match[1] {
		case "char":
			_, err = s.db.NewInsert().Model(&RecordChar{RecordID: record.ID, CharID: id}).Exec(context.Background())
		case "npc":
			_, err = s.db.NewInsert().Model(&RecordNPC{RecordID: record.ID, NPCID: id}).Exec(context.Background())
		case "location":
			_, err = s.db.NewInsert().Model(&RecordLocation{RecordID: record.ID, LocationID: id}).Exec(context.Background())
		default:
			fmt.Printf("error during record mention extracting: mention %s is incorrect in record %d", match[0], record.ID)
			// ++ add error logger ++ //
		}
		// Return on Insert Error
		if err != nil {
			return err
		}
	}

	return err
}

func (s *Storage) DeleteMentionsForRecord(record *Record) error {
	_, err := s.db.NewDelete().Model(&RecordChar{}).Where("record_id = ?", record.ID).Exec(context.Background())
	if err != nil {
		return err
	}
	_, err = s.db.NewDelete().Model(&RecordNPC{}).Where("record_id = ?", record.ID).Exec(context.Background())
	if err != nil {
		return err
	}
	_, err = s.db.NewDelete().Model(&RecordLocation{}).Where("record_id = ?", record.ID).Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetAllowedRecords(records []Record, playerID int) ([]Record, error) {
	err := s.db.NewSelect().Model(&records).WherePK().
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("hidden_by = 0").WhereOr("hidden_by = ?", playerID)
		}).
		Scan(context.Background(), &records)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return []Record{}, nil
	}

	return records, nil
}

func (s *Storage) GetAllowedNPCs(npcs []NPC, playerID int) ([]NPC, error) {
	err := s.db.NewSelect().Model(&npcs).WherePK().
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("hidden_by = 0").WhereOr("hidden_by = ?", playerID)
		}).
		Scan(context.Background(), &npcs)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return []NPC{}, nil
	}

	return npcs, nil
}

func (s *Storage) GetAllowedLocations(locations []Location, playerID int) ([]Location, error) {
	err := s.db.NewSelect().Model(&locations).WherePK().
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("hidden_by = 0").WhereOr("hidden_by = ?", playerID)
		}).
		Scan(context.Background(), &locations)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return []Location{}, nil
	}

	return locations, nil
}

func (s *Storage) GetAllowedQuests(quests []Quest, playerID int) ([]Quest, error) {
	err := s.db.NewSelect().Model(&quests).WherePK().
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("hidden_by = 0").WhereOr("hidden_by = ?", playerID)
		}).
		Scan(context.Background(), &quests)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return []Quest{}, nil
	}

	return quests, nil
}
