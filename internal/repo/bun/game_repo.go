package bunrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	"personae-fasti/internal/repo"
	"time"

	"github.com/uptrace/bun"
)

var _ repo.GameRepository = (*GameRepo)(nil)

type GameRepo struct {
	db *bun.DB
}

func NewGameRepo(db *bun.DB) *GameRepo {
	return &GameRepo{db: db}
}

func (r *GameRepo) GetCurrentGame(ctx context.Context, playerCurrentGameID int) (*domain.Game, error) {
	game := new(domain.Game)
	err := r.db.NewSelect().
		Model(game).
		Where("game.id = ?", playerCurrentGameID).
		Relation("GM").
		Relation("Players").
		Relation("Sessions").
		Relation("Settings").
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrNotFound
	}
	return game, err
}

func (r *GameRepo) GetByExt(ctx context.Context, gameExt string) (*domain.Game, error) {
	game := new(domain.Game)
	err := r.db.NewSelect().
		Model(game).
		Where("game.ext = ?", gameExt).
		Relation("GM").
		Relation("Players").
		Relation("Sessions").
		Relation("Settings").
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrNotFound
	}
	return game, err
}

func (r *GameRepo) Create(ctx context.Context, game *domain.Game) (*domain.Game, error) {
	_, err := r.db.NewInsert().
		Model(game).
		Returning("*").
		Exec(ctx)
	return game, err
}

func (r *GameRepo) UpdateByExt(ctx context.Context, game *domain.Game) (*domain.Game, error) {
	_, err := r.db.NewUpdate().
		Model(game).
		Where("ext = ?", game.ExtID).
		Set("name = ?", game.Name).
		Returning("*").
		Exec(ctx)
	return game, err
}

// func (r *GameRepo) GetByID(ctx context.Context, id int) (*domain.Game, error) {
// 	game := new(domain.Game)
// 	err := r.db.NewSelect().
// 		Model(game).
// 		Where("id = ?", id).
// 		Scan(ctx)
// 	if errors.Is(err, sql.ErrNoRows) {
// 		return nil, repository.ErrNotFound
// 	}
// 	return game, err
// }

// func (s *Storage) GetCurrentGameSessions(game *Game) ([]Session, error) {
// 	if err := s.db.NewSelect().Model(game).WherePK().Relation("Sessions").Scan(context.Background()); err == sql.ErrNoRows || game.Sessions == nil {
// 		return []Session{}, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return game.Sessions, nil
// }

func (r *GameRepo) GetCurrentGameSession(ctx context.Context, gameID int) (*domain.Session, error) {
	var currentSession domain.Session

	if err := r.db.NewSelect().Model(&currentSession).Where("game_id = ? AND end_time IS NULL", gameID).Scan(context.Background()); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &currentSession, nil
}

func (r *GameRepo) GetGameSessionByNumber(ctx context.Context, gameID int, sessionNumber int) (*domain.Session, error) {
	var currentSession domain.Session

	if err := r.db.NewSelect().Model(&currentSession).Where("game_id = ? AND number = ?", gameID, sessionNumber).Scan(context.Background()); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &currentSession, nil
}

func (r *GameRepo) CreateNewSession(ctx context.Context, game *domain.Game) (*domain.Session, error) {
	var newSession *domain.Session
	currentSession, err := r.GetCurrentGameSession(ctx, game.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	err = r.db.RunInTx(context.Background(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		sessionNumber := 0
		currentTime := time.Now().UTC()

		// Start new session
		if currentSession != nil {
			currentSession.EndTime = &currentTime

			if _, err := tx.NewUpdate().Model(currentSession).Column("end_time").WherePK().Exec(context.Background()); err == sql.ErrNoRows {
				return fmt.Errorf("cannot update previous session row")
			} else if err != nil {
				return err
			}

			sessionNumber = currentSession.Number + 1

			// Start first session
		} else {
			sessionZero := &domain.Session{
				GameID:  game.ID,
				Number:  sessionNumber,
				EndTime: &currentTime,
			}

			_, err = tx.NewInsert().Model(sessionZero).Exec(context.Background())
			if err != nil {
				return err
			}

			sessionNumber++
		}

		newSession = &domain.Session{
			GameID: game.ID,
			Number: sessionNumber,
		}

		if _, err = tx.NewInsert().Model(newSession).Exec(context.Background()); err == sql.ErrNoRows {
			return fmt.Errorf("cannot create new session row")
		} else if err != nil {
			return err
		}

		return nil
	})

	return newSession, nil
}

func (r *GameRepo) EditSession(ctx context.Context, game *domain.Game, sessionUpdate *dto.SessionUpdate) (*domain.Session, error) {
	sessionToUpdate, err := r.GetGameSessionByNumber(ctx, game.ID, sessionUpdate.Number)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	err = r.db.RunInTx(context.Background(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		sessionToUpdate.Name = sessionUpdate.Name

		if _, err := tx.NewUpdate().Model(sessionToUpdate).Column("name").WherePK().Returning("*").Exec(context.Background()); err == sql.ErrNoRows {
			return fmt.Errorf("cannot update session")
		} else if err != nil {
			return err
		}

		if sessionToUpdate.Number > 0 {
			previousSession, err := r.GetGameSessionByNumber(ctx, game.ID, sessionUpdate.Number-1)
			if err != nil {
				return err
			}
			previousSession.EndTime = &sessionUpdate.StartTime
			if _, err := tx.NewUpdate().Model(previousSession).Column("end_time").WherePK().Returning("*").Exec(context.Background()); err == sql.ErrNoRows {
				return fmt.Errorf("cannot update previous session end time")
			} else if err != nil {
				return err
			}
		}

		return nil
	})

	return sessionToUpdate, nil
}

func (r *GameRepo) RemoveLastSession(ctx context.Context, game *domain.Game) error {
	currentSession, err := r.GetCurrentGameSession(ctx, game.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err := r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var previousSession domain.Session
		if currentSession.Number > 0 {
			if err := tx.NewSelect().Model(&previousSession).Where("game_id = ? AND number = ?", game.ID, currentSession.Number-1).Scan(ctx, &previousSession); err != nil {
				return err
			}
		}
		if _, err := tx.NewUpdate().Model((*domain.Session)(nil)).Set("end_time = NULL").Where("id = ?", previousSession.ID).Returning("*").Exec(ctx, &previousSession); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model(currentSession).WherePK().Returning("*").Exec(ctx, currentSession); err != nil {
			return err
		}
		return nil

	}); err != nil {
		return err
	}

	return nil
}
