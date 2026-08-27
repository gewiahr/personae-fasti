package bunrepo

import (
	"context"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/repo"
	"time"

	"github.com/uptrace/bun"
)

var _ repo.LogRepository = (*LogRepo)(nil)

type LogRepo struct {
	db *bun.DB
}

func (r *LogRepo) Prune(ctx context.Context, successBefore, errorBefore *time.Time) error {
	_, err := r.db.NewDelete().Model((*domain.ApiLog)(nil)).
		Where("(? IS NOT NULL AND code < 400 AND created < ?)", successBefore, successBefore).
		WhereOr("(? IS NOT NULL AND code >= 400 AND created < ?)", errorBefore, errorBefore).
		Exec(ctx)
	return err
}

func NewLogRepo(db *bun.DB) *LogRepo { return &LogRepo{db: db} }

func (r *LogRepo) Insert(ctx context.Context, log *domain.ApiLog) error {
	_, err := r.db.NewInsert().Model(log).Exec(ctx)
	return err
}
