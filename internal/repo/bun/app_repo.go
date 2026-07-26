package bunrepo

import (
	"context"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/repo"

	"github.com/uptrace/bun"
)

var _ repo.AppRepository = (*AppRepo)(nil)

type AppRepo struct {
	db *bun.DB
}

func NewAppRepo(db *bun.DB) *AppRepo { return &AppRepo{db: db} }

func (r *AppRepo) InsertFeedback(ctx context.Context, feedback *domain.ServiceFeedback) (*domain.ServiceFeedback, error) {
	err := r.db.NewInsert().
		Model(feedback).
		Returning("*").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return feedback, nil
}
