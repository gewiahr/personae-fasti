package bunrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/repo"
	gewiutils "personae-fasti/pkg/gewi-utils"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

var _ repo.QuestRepository = (*QuestRepo)(nil)

type QuestRepo struct {
	db *bun.DB
}

func NewQuestRepo(db *bun.DB) *QuestRepo {
	return &QuestRepo{db: db}
}

func (r *QuestRepo) GetCurrentGameQuestList(ctx context.Context, gameID, playerID int) ([]domain.Quest, error) {
	questList := []domain.Quest{}
	err := r.db.NewSelect().
		Model(&questList).
		Where("game_id = ? AND (hidden_by = 0 OR hidden_by = ?)", gameID, playerID).
		Where("deleted IS NULL").
		// WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
		// 	return q.Where("hidden_by = 0").WhereOr("hidden_by = ?", playerID)
		// }).
		Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, e.ErrInternal
	}
	return questList, err
}

func (r *QuestRepo) GetPlayerCurrentGameQuestByID(ctx context.Context, questID int) (*domain.Quest, error) {
	quest := domain.Quest{
		ID: questID,
	}

	if err := r.db.NewSelect().
		Model(&quest).
		WherePK().
		Relation("Records").
		Relation("Tasks").
		Scan(context.Background()); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &quest, nil
}

func (r *QuestRepo) CreatePlayerCurrentGameQuest(ctx context.Context, questCreateData *dto.QuestCreateData, playerID, gameID int) (*domain.Quest, error) {
	var quest *domain.Quest

	if err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		quest = &domain.Quest{
			Name:        questCreateData.Quest.Name,
			Title:       questCreateData.Quest.Title,
			Description: questCreateData.Quest.Description,
			GameID:      gameID,
			ParentID:    questCreateData.Quest.ParentID,
			ChildID:     questCreateData.Quest.ChildID,
			HeadID:      questCreateData.Quest.HeadID,
			Successful:  questCreateData.Quest.Successful,
			HiddenBy:    gewiutils.TernaryInt(questCreateData.Quest.Hidden, playerID, 0),
		}

		_, err := tx.NewInsert().Model(quest).
			Column("name", "title", "description", "game_id", "parent_id", "child_id", "head_id", "successful", "hidden_by").
			Returning("*").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to insert quest: %w", err)
		}

		if len(questCreateData.Tasks) > 0 {
			questTasks := make([]*domain.QuestTask, len(questCreateData.Tasks))
			for i, taskCreate := range questCreateData.Tasks {
				questTasks[i] = &domain.QuestTask{
					GameID:      gameID,
					QuestID:     quest.ID,
					Name:        taskCreate.Name,
					Description: taskCreate.Description,
					Type:        domain.QuestTaskType(taskCreate.Type),
					Capacity:    taskCreate.Capacity,
					HiddenBy:    gewiutils.TernaryInt(taskCreate.Hidden, playerID, 0),
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
	}); err != nil {
		return nil, err
	}

	return quest, nil
}

func (r *QuestRepo) EditPlayerCurrentGameQuest(ctx context.Context, questUpdateData *dto.QuestUpdateData, playerID int) (*domain.Quest, error) {
	quest := &domain.Quest{ID: questUpdateData.Quest.ID}

	if err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(quest).WherePK().
			Set("name = ?", questUpdateData.Quest.Name).
			Set("title = ?", questUpdateData.Quest.Title).
			Set("description = ?", questUpdateData.Quest.Description).
			Set("hidden_by = ?", gewiutils.TernaryInt(questUpdateData.Quest.Hidden, playerID, 0)).
			Returning("*").Exec(context.Background()); err != nil {
			return fmt.Errorf("failed to update quest: %w", err)
		}

		if len(questUpdateData.Tasks) == 0 {

			if _, err := tx.NewDelete().
				Model((*domain.QuestTask)(nil)).
				Where("quest_id = ?", quest.ID).
				Exec(ctx); err != nil {
				return fmt.Errorf("failed to update quest: %w", err)
			}

		} else {

			var values []any
			var valuePlaceholders []string

			for _, task := range questUpdateData.Tasks {
				hiddenBy := gewiutils.TernaryInt(task.Hidden, playerID, 0)
				values = append(values,
					task.ID,
					task.Name,
					task.Description,
					task.Type,
					task.Capacity,
					hiddenBy,
				)
			}

			for range questUpdateData.Tasks {
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
	}); err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	if err := r.db.NewSelect().Model(quest).WherePK().Relation("Records").Relation("Tasks").Scan(context.Background()); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return quest, nil
}

// func (s *Storage) DeleteQuest(questID int, p *Player) error {
// 	now := time.Now().UTC()
// 	quest := Quest{
// 		ID:      questID,
// 		Deleted: &now,
// 	}

// 	// Delete Quest
// 	result, err := s.db.NewUpdate().Model(&quest).Column("deleted").WherePK().Exec(context.Background())
// 	if err != nil {
// 		return err
// 	}
// 	if result == nil {
// 		return fmt.Errorf("empty delete")
// 	}

// 	return nil
// }

func (r *QuestRepo) FinishPlayerCurrentGameQuest(ctx context.Context, questID int, successful bool) (*domain.Quest, error) {
	quest := &domain.Quest{ID: questID}
	if _, err := r.db.NewUpdate().
		Model(quest).
		WherePK().
		Set("finished = ?", time.Now().UTC()).
		Set("successful = ?", successful).
		Returning("*").
		Exec(context.Background()); err != nil {
		return nil, err
	}
	return quest, nil
}

func (r *QuestRepo) ResetPlayerCurrentGameQuest(ctx context.Context, questID int) (*domain.Quest, error) {
	quest := &domain.Quest{ID: questID}
	if _, err := r.db.NewUpdate().
		Model(quest).
		WherePK().
		Set("finished = ?", nil).
		Set("successful = ?", false).
		Returning("*").
		Exec(context.Background()); err != nil {
		return nil, err
	}
	return quest, nil
}

func (r *QuestRepo) GetTasksByQuest(ctx context.Context, quest *domain.Quest) ([]domain.QuestTask, error) {
	tasks := []domain.QuestTask{}

	if err := r.db.NewSelect().Model(&tasks).Where("quest_id = ?", quest.ID).Scan(context.Background()); err == sql.ErrNoRows {
		return tasks, nil
	} else if err != nil {
		return nil, err
	}

	return tasks, nil
}

// func (r *QuestRepo) CreateTasks(ctx context.Context, tasks []domain.QuestTask) ([]domain.QuestTask, error) {

// }

// func (r *QuestRepo) EditTasks(ctx context.Context, tasksUpdate []dto.TaskUpdate, playerID int) ([]domain.QuestTask, error) {

// }

func (r *QuestRepo) UpdateTasks(ctx context.Context, tasksPatch []dto.TaskPatch, quest *domain.Quest) ([]domain.QuestTask, error) {
	if len(tasksPatch) == 0 || len(quest.Tasks) == 0 {
		return nil, errors.New("empty tasks on update or quest itself")
	}

	var tasks = quest.Tasks
	var finishTime = time.Now().UTC()
	for i := range tasks {
		for _, task := range tasksPatch {
			if tasks[i].ID == task.ID {
				tasks[i].Current = task.Current
				switch tasks[i].Type {
				case domain.Binary:
					if tasks[i].Current > 0 {
						tasks[i].Finished = &finishTime
					} else {
						tasks[i].Finished = nil
					}
				case domain.Decimal:
					if tasks[i].Current >= tasks[i].Capacity {
						tasks[i].Finished = &finishTime
					} else {
						tasks[i].Finished = nil
					}
				}
			}
		}
	}

	_, err := r.db.NewUpdate().Model(&tasks).Column("current", "finished").Bulk().Returning("*").Exec(context.Background())
	return tasks, err
}

func (r *QuestRepo) FilterAllowedTasks(ctx context.Context, tasks []domain.QuestTask, playerID int) ([]domain.QuestTask, error) {
	err := r.db.NewSelect().Model(&tasks).WherePK().
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("hidden_by = 0").WhereOr("hidden_by = ?", playerID)
		}).
		Scan(context.Background(), &tasks)
	if err == sql.ErrNoRows {
		return []domain.QuestTask{}, nil
	} else if err != nil {
		return nil, err
	}

	return tasks, nil
}
