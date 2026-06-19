package bunrepo

import (
	"context"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/repo"

	"github.com/uptrace/bun"
)

var _ repo.LogRepository = (*LogRepo)(nil)

type LogRepo struct {
	db *bun.DB
}

func NewLogRepo(db *bun.DB) *LogRepo { return &LogRepo{db: db} }

func (r *LogRepo) Insert(ctx context.Context, log *domain.ApiLog) error {
	_, err := r.db.NewInsert().Model(log).Exec(ctx)
	return err
}
