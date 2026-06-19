package bunrepo

import (
	"context"
	"database/sql"
	"errors"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/repo"

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

// func (r *GameRepo) Create(ctx context.Context, game *domain.Game) error {
// 	_, err := r.db.NewInsert().Model(game).Exec(ctx)
// 	return err
// }

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
