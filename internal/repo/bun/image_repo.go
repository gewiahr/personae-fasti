package bunrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"personae-fasti/internal/domain"
	"personae-fasti/internal/repo"

	"github.com/uptrace/bun"
)

var _ repo.ImageRepository = (*ImageRepo)(nil)

type ImageRepo struct {
	db *bun.DB
}

func NewImageRepo(db *bun.DB) *ImageRepo {
	return &ImageRepo{db: db}
}

func (r *ImageRepo) ListByEntity(ctx context.Context, gameID int, entityType string, entityID int) ([]domain.Image, error) {
	images := []domain.Image{}
	err := r.db.NewSelect().
		Model(&images).
		Where("game_id = ? AND entity_type = ? AND entity_id = ?", gameID, entityType, entityID).
		Where("status = ? AND deleted IS NULL", domain.ImageStatusComplete).
		OrderExpr("is_main DESC, created ASC").
		Scan(ctx)
	return images, err
}

func (r *ImageRepo) GetByExt(ctx context.Context, gameID int, imageExt string) (*domain.Image, error) {
	image := new(domain.Image)
	err := r.db.NewSelect().
		Model(image).
		Where("game_id = ? AND ext = ? AND status = ? AND deleted IS NULL", gameID, imageExt, domain.ImageStatusComplete).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return image, err
}

func (r *ImageRepo) CreateExternal(ctx context.Context, image *domain.Image) (*domain.Image, error) {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := lockImageEntity(ctx, tx, image.EntityType, image.EntityID, image.GameID); err != nil {
			return err
		}

		var hasMain bool
		if err := tx.NewSelect().
			ColumnExpr("EXISTS (SELECT 1 FROM image WHERE entity_type = ? AND entity_id = ? AND status = ? AND is_main AND deleted IS NULL)", image.EntityType, image.EntityID, domain.ImageStatusComplete).
			Scan(ctx, &hasMain); err != nil {
			return err
		}
		image.IsMain = !hasMain

		_, err := tx.NewInsert().Model(image).Returning("*").Exec(ctx)
		return err
	})
	return image, err
}

func (r *ImageRepo) CreatePendingUpload(ctx context.Context, image *domain.Image) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := lockImageEntity(ctx, tx, image.EntityType, image.EntityID, image.GameID); err != nil {
			return err
		}

		quota := new(domain.GameImageQuota)
		if err := tx.NewSelect().Model(quota).Where("game_id = ?", image.GameID).For("UPDATE").Scan(ctx); errors.Is(err, sql.ErrNoRows) {
			return repo.ErrUploadDisabled
		} else if err != nil {
			return err
		}
		if quota.MaxBytes <= 0 {
			return repo.ErrUploadDisabled
		}
		if image.ByteSize > quota.MaxBytes-quota.UsedBytes-quota.ReservedBytes {
			return repo.ErrQuotaExceeded
		}
		if quota.MaxImages > 0 {
			count, err := tx.NewSelect().Model((*domain.Image)(nil)).
				Where("game_id = ? AND source_type = ? AND status IN (?, ?) AND deleted IS NULL", image.GameID, domain.ImageSourceUploaded, domain.ImageStatusPending, domain.ImageStatusComplete).
				Count(ctx)
			if err != nil {
				return err
			}
			if count >= quota.MaxImages {
				return repo.ErrImageLimit
			}
		}

		if _, err := tx.NewUpdate().Model(quota).
			Set("reserved_bytes = reserved_bytes + ?", image.ByteSize).
			WherePK().Exec(ctx); err != nil {
			return err
		}
		_, err := tx.NewInsert().Model(image).Returning("*").Exec(ctx)
		return err
	})
}

func (r *ImageRepo) CompleteUpload(ctx context.Context, image *domain.Image) (*domain.Image, error) {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := lockImageEntity(ctx, tx, image.EntityType, image.EntityID, image.GameID); err != nil {
			return err
		}
		pending := new(domain.Image)
		if err := tx.NewSelect().Model(pending).
			Where("id = ? AND status = ? AND deleted IS NULL", image.ID, domain.ImageStatusPending).
			For("UPDATE").Scan(ctx); errors.Is(err, sql.ErrNoRows) {
			return repo.ErrNotFound
		} else if err != nil {
			return err
		}

		quota := new(domain.GameImageQuota)
		if err := tx.NewSelect().Model(quota).Where("game_id = ?", image.GameID).For("UPDATE").Scan(ctx); err != nil {
			return err
		}
		var hasMain bool
		if err := tx.NewSelect().
			ColumnExpr("EXISTS (SELECT 1 FROM image WHERE entity_type = ? AND entity_id = ? AND status = ? AND is_main AND deleted IS NULL)", image.EntityType, image.EntityID, domain.ImageStatusComplete).
			Scan(ctx, &hasMain); err != nil {
			return err
		}
		if _, err := tx.NewUpdate().Model(quota).
			Set("reserved_bytes = reserved_bytes - ?", pending.ByteSize).
			Set("used_bytes = used_bytes + ?", pending.ByteSize).
			WherePK().Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewUpdate().Model(pending).
			Set("status = ?", domain.ImageStatusComplete).
			Set("is_main = ?", !hasMain).
			WherePK().Returning("*").Exec(ctx); err != nil {
			return err
		}
		pending.Status = domain.ImageStatusComplete
		pending.IsMain = !hasMain
		*image = *pending
		return nil
	})
	return image, err
}

func (r *ImageRepo) AbortUpload(ctx context.Context, image *domain.Image) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		pending := new(domain.Image)
		if err := tx.NewSelect().Model(pending).
			Where("id = ? AND status = ?", image.ID, domain.ImageStatusPending).
			For("UPDATE").Scan(ctx); errors.Is(err, sql.ErrNoRows) {
			return nil
		} else if err != nil {
			return err
		}
		quota := new(domain.GameImageQuota)
		if err := tx.NewSelect().Model(quota).Where("game_id = ?", image.GameID).For("UPDATE").Scan(ctx); err != nil {
			return err
		}
		if _, err := tx.NewUpdate().Model(quota).
			Set("reserved_bytes = GREATEST(0, reserved_bytes - ?)", pending.ByteSize).
			WherePK().Exec(ctx); err != nil {
			return err
		}
		_, err := tx.NewDelete().Model(pending).WherePK().ForceDelete().Exec(ctx)
		return err
	})
}

func (r *ImageRepo) SetMain(ctx context.Context, image *domain.Image) (*domain.Image, error) {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := lockImageEntity(ctx, tx, image.EntityType, image.EntityID, image.GameID); err != nil {
			return err
		}
		selected := new(domain.Image)
		if err := tx.NewSelect().
			Model(selected).
			Where("id = ? AND game_id = ? AND entity_type = ? AND entity_id = ?", image.ID, image.GameID, image.EntityType, image.EntityID).
			Where("status = ? AND deleted IS NULL", domain.ImageStatusComplete).
			For("UPDATE").
			Scan(ctx); errors.Is(err, sql.ErrNoRows) {
			return repo.ErrNotFound
		} else if err != nil {
			return err
		}
		if _, err := tx.NewUpdate().
			Model((*domain.Image)(nil)).
			Set("is_main = false").
			Where("game_id = ? AND entity_type = ? AND entity_id = ? AND deleted IS NULL", image.GameID, image.EntityType, image.EntityID).
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewUpdate().
			Model(selected).
			Set("is_main = true").
			WherePK().
			Returning("*").
			Exec(ctx); err != nil {
			return err
		}
		*image = *selected
		image.IsMain = true
		return nil
	})
	return image, err
}

func (r *ImageRepo) SoftDelete(ctx context.Context, image *domain.Image) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := lockImageEntity(ctx, tx, image.EntityType, image.EntityID, image.GameID); err != nil {
			return err
		}
		current := new(domain.Image)
		if err := tx.NewSelect().
			Model(current).
			Where("id = ? AND game_id = ? AND entity_type = ? AND entity_id = ?", image.ID, image.GameID, image.EntityType, image.EntityID).
			Where("status = ? AND deleted IS NULL", domain.ImageStatusComplete).
			For("UPDATE").
			Scan(ctx); errors.Is(err, sql.ErrNoRows) {
			return repo.ErrNotFound
		} else if err != nil {
			return err
		}

		wasMain := current.IsMain
		deletedAt := time.Now().UTC()
		result, err := tx.NewUpdate().
			Model((*domain.Image)(nil)).
			Set("status = ?", domain.ImageStatusDeleted).
			Set("deleted = ?", deletedAt).
			Set("is_main = false").
			Where("id = ? AND deleted IS NULL", current.ID).
			Exec(ctx)
		if err != nil {
			return err
		}
		if affected, err := result.RowsAffected(); err != nil || affected != 1 {
			if err != nil {
				return err
			}
			return repo.ErrNotFound
		}
		if !wasMain {
			return nil
		}

		next := new(domain.Image)
		err = tx.NewSelect().
			Model(next).
			Where("game_id = ? AND entity_type = ? AND entity_id = ?", image.GameID, image.EntityType, image.EntityID).
			Where("status = ? AND deleted IS NULL", domain.ImageStatusComplete).
			OrderExpr("created ASC").
			Limit(1).
			For("UPDATE").
			Scan(ctx)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		_, err = tx.NewUpdate().Model(next).Set("is_main = true").WherePK().Exec(ctx)
		return err
	})
}

func (r *ImageRepo) GetQuota(ctx context.Context, gameID int) (*domain.GameImageQuota, error) {
	quota := &domain.GameImageQuota{GameID: gameID}
	if _, err := r.db.NewInsert().Model(quota).On("CONFLICT (game_id) DO NOTHING").Exec(ctx); err != nil {
		return nil, err
	}
	err := r.db.NewSelect().Model(quota).Where("game_id = ?", gameID).Scan(ctx)
	return quota, err
}

func lockImageEntity(ctx context.Context, tx bun.Tx, entityType string, entityID, gameID int) error {
	table, err := imageEntityTable(entityType)
	if err != nil {
		return err
	}
	var id int
	err = tx.NewSelect().
		TableExpr(table).
		Column("id").
		Where("id = ? AND game_id = ?", entityID, gameID).
		For("UPDATE").
		Scan(ctx, &id)
	if errors.Is(err, sql.ErrNoRows) {
		return repo.ErrNotFound
	}
	return err
}

func imageEntityTable(entityType string) (string, error) {
	switch entityType {
	case "char":
		return `"char"`, nil
	case "npc":
		return "npc", nil
	case "location":
		return "location", nil
	default:
		return "", fmt.Errorf("unsupported image entity type %q", entityType)
	}
}
