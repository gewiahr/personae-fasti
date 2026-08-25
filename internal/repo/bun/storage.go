package bunrepo

import (
	"context"
	"database/sql"
	"fmt"
	"personae-fasti/configs"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/repo"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type BunStorage struct {
	db *bun.DB
}

func NewBunStorage(c *configs.DBConfig) (*BunStorage, error) {
	ctx := context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Name)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	if err := sqldb.Ping(); err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	storage := &BunStorage{
		db: db,
	}
	if err := storage.Migrate(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	return storage, nil
}

func (s *BunStorage) GameRepo() repo.GameRepository         { return NewGameRepo(s.db) }
func (s *BunStorage) PlayerRepo() repo.PlayerRepository     { return NewPlayerRepo(s.db) }
func (s *BunStorage) RecordRepo() repo.RecordRepository     { return NewRecordRepo(s.db) }
func (s *BunStorage) EntitiesRepo() repo.EntitiesRepository { return NewEntitiesRepo(s.db) }
func (s *BunStorage) QuestRepo() repo.QuestRepository       { return NewQuestRepo(s.db) }
func (s *BunStorage) LogRepo() repo.LogRepository           { return NewLogRepo(s.db) }
func (s *BunStorage) AppRepo() repo.AppRepository           { return NewAppRepo(s.db) }

func (s *BunStorage) Migrate(ctx context.Context) error {
	err := addNanoID(s.db)
	if err != nil {
		return fmt.Errorf("failed to create nanoid function: %w", err)
	}

	models := []any{
		(*domain.PlayerGame)(nil),
		(*domain.GameInvite)(nil),
		(*domain.RecordChar)(nil),
		(*domain.RecordNPC)(nil),
		(*domain.RecordLocation)(nil),
	}
	for _, m := range models {
		s.db.RegisterModel(m)
	}

	tables := []any{
		(*domain.Game)(nil),
		(*domain.GameSettings)(nil),
		(*domain.Player)(nil),
		(*domain.PlayerRegData)(nil),
		(*domain.Token)(nil),
		(*domain.Telegram)(nil),
		(*domain.Char)(nil),
		(*domain.NPC)(nil),
		(*domain.Location)(nil),
		(*domain.Record)(nil),
		(*domain.Session)(nil),
		(*domain.Quest)(nil),
		(*domain.QuestTask)(nil),
		(*domain.PlayerGame)(nil),
		(*domain.GameInvite)(nil),
		(*domain.RecordChar)(nil),
		(*domain.RecordNPC)(nil),
		(*domain.RecordLocation)(nil),
		(*domain.ApiLog)(nil),
		(*domain.ServiceFeedback)(nil),
	}

	for _, table := range tables {
		_, err := s.db.NewCreateTable().IfNotExists().Model(table).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create table for %T: %w", table, err)
		}
	}

	if err := s.applyPublicExtMigration(ctx); err != nil {
		return err
	}
	if err := s.applyPersonalNoteMigration(ctx); err != nil {
		return err
	}

	_, err = s.db.NewCreateIndex().
		IfNotExists().
		Model((*domain.Token)(nil)).
		Index("token_idx").
		Column("token_hash").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	return nil
}

func (s *BunStorage) Close() error {
	return s.db.Close()
}

func addNanoID(db *bun.DB) error {
	const nanoidFunction = `
        CREATE OR REPLACE FUNCTION nanoid(
            size INT DEFAULT 12,
            alphabet TEXT DEFAULT '346789ACDEFGHJKMNPQRTWXY'
        )
        RETURNS TEXT AS $$
        DECLARE
            result TEXT;
        BEGIN
            SELECT string_agg(
                substr(alphabet, (abs(hashtext(gen_random_uuid()::text)) % length(alphabet))::int + 1, 1), 
                ''
            )
            FROM generate_series(1, size)
            INTO result;
            
            RETURN result;
        END;
        $$ LANGUAGE plpgsql VOLATILE PARALLEL SAFE;
    `

	if _, err := db.ExecContext(context.Background(), `CREATE EXTENSION IF NOT EXISTS pgcrypto`); err != nil {
		return err
	}

	if _, err := db.ExecContext(context.Background(), nanoidFunction); err != nil {
		return err
	}

	return nil
}

const publicExtMigration = "20260824_public_entity_exts"

const personalNoteMigration = "20260825_player_personal_note"

func (s *BunStorage) applyPersonalNoteMigration(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migration (
			name TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
		)
	`); err != nil {
		return fmt.Errorf("failed to create schema migration table: %w", err)
	}

	var applied bool
	if err := s.db.NewSelect().
		ColumnExpr("EXISTS (SELECT 1 FROM schema_migration WHERE name = ?)", personalNoteMigration).
		Scan(ctx, &applied); err != nil {
		return fmt.Errorf("failed to check personal note migration: %w", err)
	}
	if applied {
		return nil
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.ExecContext(ctx, `
			ALTER TABLE player
			ADD COLUMN IF NOT EXISTS personal_note TEXT NOT NULL DEFAULT ''
		`); err != nil {
			return fmt.Errorf("personal note migration failed: %w", err)
		}

		if _, err := tx.ExecContext(
			ctx,
			"INSERT INTO schema_migration (name) VALUES (?)",
			personalNoteMigration,
		); err != nil {
			return fmt.Errorf("failed to record personal note migration: %w", err)
		}

		return nil
	})
}

func (s *BunStorage) applyPublicExtMigration(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migration (
			name TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
		)
	`); err != nil {
		return fmt.Errorf("failed to create schema migration table: %w", err)
	}

	var applied bool
	if err := s.db.NewSelect().
		ColumnExpr("EXISTS (SELECT 1 FROM schema_migration WHERE name = ?)", publicExtMigration).
		Scan(ctx, &applied); err != nil {
		return fmt.Errorf("failed to check public ext migration: %w", err)
	}
	if applied {
		return nil
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		statements := []string{
			`ALTER TABLE "char" ADD COLUMN IF NOT EXISTS ext VARCHAR(16)`,
			`ALTER TABLE npc ADD COLUMN IF NOT EXISTS ext VARCHAR(16)`,
			`ALTER TABLE location ADD COLUMN IF NOT EXISTS ext VARCHAR(16)`,
			`ALTER TABLE quest ADD COLUMN IF NOT EXISTS ext VARCHAR(16)`,
			`UPDATE "char" SET ext = nanoid(12) WHERE ext IS NULL OR ext = ''`,
			`UPDATE npc SET ext = nanoid(12) WHERE ext IS NULL OR ext = ''`,
			`UPDATE location SET ext = nanoid(12) WHERE ext IS NULL OR ext = ''`,
			`UPDATE quest SET ext = nanoid(12) WHERE ext IS NULL OR ext = ''`,
			`CREATE UNIQUE INDEX IF NOT EXISTS char_ext_idx ON "char" (ext)`,
			`CREATE UNIQUE INDEX IF NOT EXISTS npc_ext_idx ON npc (ext)`,
			`CREATE UNIQUE INDEX IF NOT EXISTS location_ext_idx ON location (ext)`,
			`CREATE UNIQUE INDEX IF NOT EXISTS quest_ext_idx ON quest (ext)`,
			`ALTER TABLE "char" ALTER COLUMN ext SET DEFAULT nanoid(12), ALTER COLUMN ext SET NOT NULL`,
			`ALTER TABLE npc ALTER COLUMN ext SET DEFAULT nanoid(12), ALTER COLUMN ext SET NOT NULL`,
			`ALTER TABLE location ALTER COLUMN ext SET DEFAULT nanoid(12), ALTER COLUMN ext SET NOT NULL`,
			`ALTER TABLE quest ALTER COLUMN ext SET DEFAULT nanoid(12), ALTER COLUMN ext SET NOT NULL`,
			migrateMentionSIDsSQL,
		}

		for _, statement := range statements {
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("public ext migration failed: %w", err)
			}
		}

		if _, err := tx.ExecContext(
			ctx,
			"INSERT INTO schema_migration (name) VALUES (?)",
			publicExtMigration,
		); err != nil {
			return fmt.Errorf("failed to record public ext migration: %w", err)
		}

		return nil
	})
}

const migrateMentionSIDsSQL = `
DO $$
DECLARE
	mention RECORD;
BEGIN
	FOR mention IN SELECT id, ext FROM "char" LOOP
		UPDATE "record" SET text = replace(text, '@char:' || mention.id::text || chr(96), '@char:' || mention.ext || chr(96));
		UPDATE "char" SET description = replace(description, '@char:' || mention.id::text || chr(96), '@char:' || mention.ext || chr(96));
		UPDATE npc SET description = replace(description, '@char:' || mention.id::text || chr(96), '@char:' || mention.ext || chr(96));
		UPDATE location SET description = replace(description, '@char:' || mention.id::text || chr(96), '@char:' || mention.ext || chr(96));
		UPDATE quest SET description = replace(description, '@char:' || mention.id::text || chr(96), '@char:' || mention.ext || chr(96));
		UPDATE quest_task SET description = replace(description, '@char:' || mention.id::text || chr(96), '@char:' || mention.ext || chr(96));
	END LOOP;

	FOR mention IN SELECT id, ext FROM npc LOOP
		UPDATE "record" SET text = replace(text, '@npc:' || mention.id::text || chr(96), '@npc:' || mention.ext || chr(96));
		UPDATE "char" SET description = replace(description, '@npc:' || mention.id::text || chr(96), '@npc:' || mention.ext || chr(96));
		UPDATE npc SET description = replace(description, '@npc:' || mention.id::text || chr(96), '@npc:' || mention.ext || chr(96));
		UPDATE location SET description = replace(description, '@npc:' || mention.id::text || chr(96), '@npc:' || mention.ext || chr(96));
		UPDATE quest SET description = replace(description, '@npc:' || mention.id::text || chr(96), '@npc:' || mention.ext || chr(96));
		UPDATE quest_task SET description = replace(description, '@npc:' || mention.id::text || chr(96), '@npc:' || mention.ext || chr(96));
	END LOOP;

	FOR mention IN SELECT id, ext FROM location LOOP
		UPDATE "record" SET text = replace(text, '@location:' || mention.id::text || chr(96), '@location:' || mention.ext || chr(96));
		UPDATE "char" SET description = replace(description, '@location:' || mention.id::text || chr(96), '@location:' || mention.ext || chr(96));
		UPDATE npc SET description = replace(description, '@location:' || mention.id::text || chr(96), '@location:' || mention.ext || chr(96));
		UPDATE location SET description = replace(description, '@location:' || mention.id::text || chr(96), '@location:' || mention.ext || chr(96));
		UPDATE quest SET description = replace(description, '@location:' || mention.id::text || chr(96), '@location:' || mention.ext || chr(96));
		UPDATE quest_task SET description = replace(description, '@location:' || mention.id::text || chr(96), '@location:' || mention.ext || chr(96));
	END LOOP;
END $$;
`

// func (s *BunStorage) InitTables() {
// 	ctx := context.Background()

// 	s.db.RegisterModel((*domain.PlayerGame)(nil))
// 	s.db.RegisterModel((*domain.GameInvite)(nil))
// 	s.db.RegisterModel((*domain.RecordChar)(nil))
// 	s.db.RegisterModel((*domain.RecordNPC)(nil))
// 	s.db.RegisterModel((*domain.RecordLocation)(nil))

// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.Game)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.GameSettings)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.Player)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.PlayerRegData)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.Token)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.Telegram)(nil)).Exec(ctx)

// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.Char)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.NPC)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.Location)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.Record)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.Session)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.Quest)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.QuestTask)(nil)).Exec(ctx)

// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.PlayerGame)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.GameInvite)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.RecordChar)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.RecordNPC)(nil)).Exec(ctx)
// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.RecordLocation)(nil)).Exec(ctx)

// 	_, _ = s.db.NewCreateTable().IfNotExists().Model((*domain.ApiLog)(nil)).Exec(ctx)

// 	_, err := s.db.NewCreateTable().IfNotExists().Model((*domain.ServiceFeedback)(nil)).Exec(ctx)
// 	if err != nil {
// 		fmt.Printf("error occured during index creation: %v", err.Error())
// 	}

// 	if _, err := s.db.NewCreateIndex().IfNotExists().Model((*domain.Token)(nil)).Index("token_idx").Column("token_hash").Exec(ctx); err != nil {
// 		fmt.Printf("error occured during index creation: %v", err.Error())
// 	}

// }
