package bunrepo

import (
	"context"
	"database/sql"
	"fmt"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	"personae-fasti/internal/repo"
	gewiutils "personae-fasti/pkg/gewi-utils"
	"regexp"
	"time"

	"github.com/uptrace/bun"
)

var _ repo.RecordRepository = (*RecordRepo)(nil)

type RecordRepo struct {
	db *bun.DB
}

func NewRecordRepo(db *bun.DB) *RecordRepo { return &RecordRepo{db: db} }

func (r *RecordRepo) GetCurrentGameRecordList(ctx context.Context, gameID, playerID int) ([]domain.Record, error) {
	records := []domain.Record{}
	err := r.db.NewSelect().Model(&records).Where("\"record\".player_id = ? AND \"record\".game_id = ? AND \"record\".deleted IS NULL", playerID, gameID).
		Relation("Quest").
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("\"record\".hidden_by = 0").WhereOr("\"record\".hidden_by = ?", playerID)
		}).
		Scan(context.Background(), &records)
	if err == sql.ErrNoRows {
		return records, nil
	} else if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *RecordRepo) GetRecord(ctx context.Context, playerID int, recordID int) (*domain.Record, error) {
	record := &domain.Record{ID: recordID}
	err := r.db.NewSelect().Model(record).WherePK().Scan(context.Background(), record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r *RecordRepo) PostRecord(ctx context.Context, record *domain.Record) (*domain.Record, error) {
	if err := r.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		result, err := tx.NewInsert().Model(record).Returning("*").Exec(context.Background())
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("empty insert")
		}
		if err := r.InsertMentionsForRecord(ctx, tx, record); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return record, nil
}

func (r *RecordRepo) EditRecord(ctx context.Context, recordUpdate *dto.RecordUpdate, playerID int) (*domain.Record, error) {
	record := &domain.Record{ID: recordUpdate.ID}

	if err := r.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		result, err := tx.NewUpdate().
			Model(record).
			WherePK().
			Set("text = ?", recordUpdate.Text).
			Set("quest_id = ?", recordUpdate.QuestID).
			Set("hidden_by = ?", gewiutils.TernaryInt(recordUpdate.Hidden, playerID, 0)).
			Set("updated = ?", time.Now().UTC()).
			Returning("*").
			Exec(context.Background())
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("empty insert")
		}
		if err := r.DeleteMentionsForRecord(ctx, tx, record); err != nil {
			return err
		}
		if err := r.InsertMentionsForRecord(ctx, tx, record); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return record, nil
}

func (r *RecordRepo) SoftDeleteRecord(ctx context.Context, playerID int, recordID int) error {
	now := time.Now().UTC()
	record := domain.Record{
		ID:      recordID,
		Deleted: &now,
	}

	// Soft Delete Record
	result, err := r.db.NewUpdate().Model(&record).Column("deleted").WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("empty delete")
	}

	return nil
}

func (r *RecordRepo) FilterAllowedRecords(ctx context.Context, records []domain.Record, playerID int) ([]domain.Record, error) {
	err := r.db.NewSelect().Model(&records).WherePK().Where("\"record\".deleted IS NULL").Relation("Quest").
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("\"record\".hidden_by = 0").WhereOr("\"record\".hidden_by = ?", playerID)
		}).
		Scan(context.Background(), &records)
	if err == sql.ErrNoRows {
		return []domain.Record{}, nil
	} else if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *RecordRepo) InsertMentionsForRecord(ctx context.Context, tx bun.Tx, record *domain.Record) error {
	re, err := regexp.Compile(`@(?P<type>\w+):(?P<ext>[\w-]+)` + "`(?P<name>[^`]+)`")
	if err != nil {
		return err
	}

	matches := re.FindAllStringSubmatch(record.Text, -1)
	for _, match := range matches {
		var id int
		switch match[1] {
		case "char":
			err = tx.NewSelect().TableExpr(`"char"`).Column("id").Where("ext = ? AND game_id = ?", match[2], record.GameID).Scan(ctx, &id)
			if err != nil {
				break
			}
			_, err = tx.NewInsert().Model(&domain.RecordChar{RecordID: record.ID, CharID: id}).On("CONFLICT DO NOTHING").Exec(ctx)
		case "npc":
			err = tx.NewSelect().Table("npc").Column("id").Where("ext = ? AND game_id = ?", match[2], record.GameID).Scan(ctx, &id)
			if err != nil {
				break
			}
			_, err = tx.NewInsert().Model(&domain.RecordNPC{RecordID: record.ID, NPCID: id}).On("CONFLICT DO NOTHING").Exec(ctx)
		case "location":
			err = tx.NewSelect().Table("location").Column("id").Where("ext = ? AND game_id = ?", match[2], record.GameID).Scan(ctx, &id)
			if err != nil {
				break
			}
			_, err = tx.NewInsert().Model(&domain.RecordLocation{RecordID: record.ID, LocationID: id}).On("CONFLICT DO NOTHING").Exec(ctx)
		default:
			fmt.Printf("error during record mention extracting: mention %s is incorrect in record %d", match[0], record.ID)
			// ++ add error logger ++ //
		}
		if err == sql.ErrNoRows {
			err = nil
			continue
		}
		if err != nil {
			return err
		}
	}

	return err
}

func (r *RecordRepo) DeleteMentionsForRecord(ctx context.Context, tx bun.Tx, record *domain.Record) error {
	_, err := tx.NewDelete().Model(&domain.RecordChar{}).Where("record_id = ?", record.ID).Exec(ctx)
	if err != nil {
		return err
	}
	_, err = tx.NewDelete().Model(&domain.RecordNPC{}).Where("record_id = ?", record.ID).Exec(ctx)
	if err != nil {
		return err
	}
	_, err = tx.NewDelete().Model(&domain.RecordLocation{}).Where("record_id = ?", record.ID).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
