package bunrepo

import (
	"context"
	"database/sql"
	"fmt"
	"personae-fasti/configs"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/repo"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type BunStorage struct {
	db *bun.DB
}

func NewBunStorage(c *configs.DBConfig) (*BunStorage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Name)
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithDialTimeout(dbTimeout(c.DialTimeout, 5*time.Second)),
		pgdriver.WithReadTimeout(dbTimeout(c.ReadTimeout, 120*time.Second)),
		pgdriver.WithWriteTimeout(dbTimeout(c.WriteTimeout, 30*time.Second)),
	))

	if err := sqldb.Ping(); err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	storage := &BunStorage{
		db: db,
	}
	storage.registerModels()

	return storage, nil
}

func dbTimeout(seconds int, fallback time.Duration) time.Duration {
	if seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}

func (s *BunStorage) registerModels() {
	models := []any{
		(*domain.PlayerGame)(nil),
		(*domain.GameInvite)(nil),
		(*domain.RecordChar)(nil),
		(*domain.RecordNPC)(nil),
		(*domain.RecordLocation)(nil),
	}
	for _, model := range models {
		s.db.RegisterModel(model)
	}
}

func (s *BunStorage) GameRepo() repo.GameRepository         { return NewGameRepo(s.db) }
func (s *BunStorage) PlayerRepo() repo.PlayerRepository     { return NewPlayerRepo(s.db) }
func (s *BunStorage) RecordRepo() repo.RecordRepository     { return NewRecordRepo(s.db) }
func (s *BunStorage) EntitiesRepo() repo.EntitiesRepository { return NewEntitiesRepo(s.db) }
func (s *BunStorage) QuestRepo() repo.QuestRepository       { return NewQuestRepo(s.db) }
func (s *BunStorage) LogRepo() repo.LogRepository           { return NewLogRepo(s.db) }
func (s *BunStorage) AppRepo() repo.AppRepository           { return NewAppRepo(s.db) }
func (s *BunStorage) ImageRepo() repo.ImageRepository       { return NewImageRepo(s.db) }

func (s *BunStorage) Migrate(ctx context.Context) error {
	err := addNanoID(s.db)
	if err != nil {
		return fmt.Errorf("failed to create nanoid function: %w", err)
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

	if err := s.applyLegacySchemaMigration(ctx); err != nil {
		return err
	}
	if err := s.applyPublicExtMigration(ctx); err != nil {
		return err
	}
	if err := s.applyPersonalNoteMigration(ctx); err != nil {
		return err
	}
	if err := s.applyImageMigration(ctx); err != nil {
		return err
	}
	if err := s.applyAPIRequestLoggingMigration(ctx); err != nil {
		return err
	}
	if err := s.applyStructuredAPILoggingMigration(ctx); err != nil {
		return err
	}
	if err := s.applyNullableAPILogFieldsMigration(ctx); err != nil {
		return err
	}
	if err := s.applyAPIClientMetadataMigration(ctx); err != nil {
		return err
	}
	if err := s.applyGameCreationRelationsMigration(ctx); err != nil {
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

func (s *BunStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
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

const legacySchemaMigration = "20260826_legacy_main_schema"

const imageMigration = "20260827_image_system"

const apiRequestLoggingMigration = "20260828_api_request_logging"

const structuredAPILoggingMigration = "20260829_structured_api_logging"

const nullableAPILogFieldsMigration = "20260830_nullable_api_log_fields"

const apiClientMetadataMigration = "20260831_api_client_metadata"

const gameCreationRelationsMigration = "20260901_game_creation_relations"

func (s *BunStorage) applyGameCreationRelationsMigration(ctx context.Context) error {
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
		ColumnExpr("EXISTS (SELECT 1 FROM schema_migration WHERE name = ?)", gameCreationRelationsMigration).
		Scan(ctx, &applied); err != nil {
		return fmt.Errorf("failed to check game creation relations migration: %w", err)
	}
	if applied {
		return nil
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		statements := []string{
			`INSERT INTO game_settings (game_id) SELECT id FROM game ON CONFLICT (game_id) DO NOTHING`,
			`INSERT INTO players_games (player_id, game_id) SELECT gm_id, id FROM game WHERE gm_id IS NOT NULL ON CONFLICT (player_id, game_id) DO NOTHING`,
			`INSERT INTO game_image_quota (game_id) SELECT id FROM game ON CONFLICT (game_id) DO NOTHING`,
		}
		for _, statement := range statements {
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("game creation relations migration failed: %w", err)
			}
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migration (name) VALUES (?)", gameCreationRelationsMigration); err != nil {
			return fmt.Errorf("failed to record game creation relations migration: %w", err)
		}
		return nil
	})
}

func (s *BunStorage) applyAPIClientMetadataMigration(ctx context.Context) error {
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
		ColumnExpr("EXISTS (SELECT 1 FROM schema_migration WHERE name = ?)", apiClientMetadataMigration).
		Scan(ctx, &applied); err != nil {
		return fmt.Errorf("failed to check API client metadata migration: %w", err)
	}
	if applied {
		return nil
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		statements := []string{
			`ALTER TABLE log_api ADD COLUMN IF NOT EXISTS ip VARCHAR(45)`,
			`ALTER TABLE log_api ADD COLUMN IF NOT EXISTS host VARCHAR(255)`,
			`CREATE INDEX IF NOT EXISTS log_api_ip_idx ON log_api (ip) WHERE ip IS NOT NULL`,
		}
		for _, statement := range statements {
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("API client metadata migration failed: %w", err)
			}
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migration (name) VALUES (?)", apiClientMetadataMigration); err != nil {
			return fmt.Errorf("failed to record API client metadata migration: %w", err)
		}
		return nil
	})
}

func (s *BunStorage) applyNullableAPILogFieldsMigration(ctx context.Context) error {
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
		ColumnExpr("EXISTS (SELECT 1 FROM schema_migration WHERE name = ?)", nullableAPILogFieldsMigration).
		Scan(ctx, &applied); err != nil {
		return fmt.Errorf("failed to check nullable API log fields migration: %w", err)
	}
	if applied {
		return nil
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		statements := []string{
			`ALTER TABLE log_api ALTER COLUMN player_id DROP NOT NULL, ALTER COLUMN player_id DROP DEFAULT`,
			`ALTER TABLE log_api ALTER COLUMN "request" DROP NOT NULL`,
			`ALTER TABLE log_api ALTER COLUMN "response" DROP NOT NULL`,
			`ALTER TABLE log_api ALTER COLUMN error_code DROP NOT NULL, ALTER COLUMN error_code DROP DEFAULT`,
			`ALTER TABLE log_api ALTER COLUMN internal_error DROP NOT NULL, ALTER COLUMN internal_error DROP DEFAULT`,
			`UPDATE log_api SET player_id = NULL WHERE player_id = 0`,
			`UPDATE log_api SET "request" = NULL WHERE "request" = ''`,
			`UPDATE log_api SET "response" = NULL WHERE "response" = ''`,
			`UPDATE log_api SET "error" = NULL WHERE "error" = ''`,
			`UPDATE log_api SET error_code = NULL WHERE error_code = ''`,
			`UPDATE log_api SET internal_error = NULL WHERE internal_error = ''`,
		}
		for _, statement := range statements {
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("nullable API log fields migration failed: %w", err)
			}
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migration (name) VALUES (?)", nullableAPILogFieldsMigration); err != nil {
			return fmt.Errorf("failed to record nullable API log fields migration: %w", err)
		}
		return nil
	})
}

func (s *BunStorage) applyStructuredAPILoggingMigration(ctx context.Context) error {
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
		ColumnExpr("EXISTS (SELECT 1 FROM schema_migration WHERE name = ?)", structuredAPILoggingMigration).
		Scan(ctx, &applied); err != nil {
		return fmt.Errorf("failed to check structured API logging migration: %w", err)
	}
	if applied {
		return nil
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		statements := []string{
			`ALTER TABLE log_api ADD COLUMN IF NOT EXISTS game_id BIGINT`,
			`ALTER TABLE log_api ADD COLUMN IF NOT EXISTS request_id VARCHAR(32) NOT NULL DEFAULT ''`,
			`ALTER TABLE log_api ADD COLUMN IF NOT EXISTS error_code VARCHAR(64) NOT NULL DEFAULT ''`,
			`ALTER TABLE log_api ADD COLUMN IF NOT EXISTS internal_error TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE log_api ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ`,
			`UPDATE log_api SET started_at = created - (time * INTERVAL '1 millisecond') WHERE started_at IS NULL`,
			`ALTER TABLE log_api ALTER COLUMN started_at SET DEFAULT current_timestamp, ALTER COLUMN started_at SET NOT NULL`,
			`CREATE INDEX IF NOT EXISTS log_api_created_idx ON log_api (created)`,
			`CREATE INDEX IF NOT EXISTS log_api_request_id_idx ON log_api (request_id) WHERE request_id <> ''`,
			`CREATE INDEX IF NOT EXISTS log_api_game_id_idx ON log_api (game_id) WHERE game_id IS NOT NULL`,
		}
		for _, statement := range statements {
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("structured API logging migration failed: %w", err)
			}
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migration (name) VALUES (?)", structuredAPILoggingMigration); err != nil {
			return fmt.Errorf("failed to record structured API logging migration: %w", err)
		}
		return nil
	})
}

func (s *BunStorage) applyAPIRequestLoggingMigration(ctx context.Context) error {
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
		ColumnExpr("EXISTS (SELECT 1 FROM schema_migration WHERE name = ?)", apiRequestLoggingMigration).
		Scan(ctx, &applied); err != nil {
		return fmt.Errorf("failed to check API request logging migration: %w", err)
	}
	if applied {
		return nil
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Old deployments silently ignored log insert errors. Repair a sequence that
		// may have fallen behind after a database restore before logging resumes.
		if _, err := tx.ExecContext(ctx, `
			SELECT setval(
				pg_get_serial_sequence('log_api', 'id'),
				COALESCE((SELECT MAX(id) + 1 FROM log_api), 1),
				false
			)
		`); err != nil {
			return fmt.Errorf("failed to repair API log sequence: %w", err)
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migration (name) VALUES (?)", apiRequestLoggingMigration); err != nil {
			return fmt.Errorf("failed to record API request logging migration: %w", err)
		}
		return nil
	})
}

func (s *BunStorage) applyLegacySchemaMigration(ctx context.Context) error {
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
		ColumnExpr("EXISTS (SELECT 1 FROM schema_migration WHERE name = ?)", legacySchemaMigration).
		Scan(ctx, &applied); err != nil {
		return fmt.Errorf("failed to check legacy schema migration: %w", err)
	}
	if applied {
		return nil
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		statements := []string{
			`ALTER TABLE game ADD COLUMN IF NOT EXISTS ext VARCHAR(16)`,
			`UPDATE game SET ext = nanoid(12) WHERE ext IS NULL OR ext = ''`,
			`CREATE UNIQUE INDEX IF NOT EXISTS game_ext_idx ON game (ext)`,
			`ALTER TABLE game ALTER COLUMN ext SET DEFAULT nanoid(12), ALTER COLUMN ext SET NOT NULL`,

			`ALTER TABLE player ADD COLUMN IF NOT EXISTS ext VARCHAR(16)`,
			`ALTER TABLE player ADD COLUMN IF NOT EXISTS email VARCHAR(255)`,
			`ALTER TABLE player ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255)`,
			`UPDATE player SET ext = nanoid(12) WHERE ext IS NULL OR ext = ''`,
			`UPDATE player SET email = '' WHERE email IS NULL`,
			`UPDATE player SET password_hash = '' WHERE password_hash IS NULL`,
			`CREATE UNIQUE INDEX IF NOT EXISTS player_ext_idx ON player (ext)`,
			`ALTER TABLE player ALTER COLUMN ext SET DEFAULT nanoid(12), ALTER COLUMN ext SET NOT NULL`,
			`ALTER TABLE player ALTER COLUMN email SET DEFAULT '', ALTER COLUMN email SET NOT NULL`,
			`ALTER TABLE player ALTER COLUMN password_hash SET DEFAULT '', ALTER COLUMN password_hash SET NOT NULL`,
			`ALTER TABLE player DROP COLUMN IF EXISTS accesskey`,

			`ALTER TABLE telegram ADD COLUMN IF NOT EXISTS pic_url VARCHAR NOT NULL DEFAULT ''`,

			`DO $$
			BEGIN
				IF EXISTS (
					SELECT 1 FROM information_schema.columns
					WHERE table_schema = 'public' AND table_name = 'log_api' AND column_name = 'user'
				) AND NOT EXISTS (
					SELECT 1 FROM information_schema.columns
					WHERE table_schema = 'public' AND table_name = 'log_api' AND column_name = 'player_id'
				) THEN
					ALTER TABLE log_api RENAME COLUMN "user" TO player_id;
				END IF;

				IF EXISTS (
					SELECT 1 FROM information_schema.columns
					WHERE table_schema = 'public' AND table_name = 'log_api' AND column_name = 'http_code'
				) AND NOT EXISTS (
					SELECT 1 FROM information_schema.columns
					WHERE table_schema = 'public' AND table_name = 'log_api' AND column_name = 'code'
				) THEN
					ALTER TABLE log_api RENAME COLUMN http_code TO code;
				END IF;

				IF EXISTS (
					SELECT 1 FROM information_schema.columns
					WHERE table_schema = 'public' AND table_name = 'log_api' AND column_name = 'time'
						AND data_type = 'timestamp with time zone'
				) AND NOT EXISTS (
					SELECT 1 FROM information_schema.columns
					WHERE table_schema = 'public' AND table_name = 'log_api' AND column_name = 'created'
				) THEN
					ALTER TABLE log_api RENAME COLUMN "time" TO created;
				END IF;
			END $$`,
			`ALTER TABLE log_api ADD COLUMN IF NOT EXISTS player_id BIGINT NOT NULL DEFAULT 0`,
			`ALTER TABLE log_api ADD COLUMN IF NOT EXISTS code BIGINT NOT NULL DEFAULT 0`,
			`ALTER TABLE log_api ADD COLUMN IF NOT EXISTS time BIGINT NOT NULL DEFAULT 0`,
			`ALTER TABLE log_api ADD COLUMN IF NOT EXISTS created TIMESTAMPTZ NOT NULL DEFAULT current_timestamp`,
		}

		for _, statement := range statements {
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("legacy schema migration failed: %w", err)
			}
		}

		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migration (name) VALUES (?)", legacySchemaMigration); err != nil {
			return fmt.Errorf("failed to record legacy schema migration: %w", err)
		}
		return nil
	})
}

func (s *BunStorage) applyImageMigration(ctx context.Context) error {
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
		ColumnExpr("EXISTS (SELECT 1 FROM schema_migration WHERE name = ?)", imageMigration).
		Scan(ctx, &applied); err != nil {
		return fmt.Errorf("failed to check image migration: %w", err)
	}
	if applied {
		return nil
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		statements := []string{
			`CREATE TABLE IF NOT EXISTS image (
				id BIGSERIAL PRIMARY KEY,
				ext VARCHAR(16) UNIQUE NOT NULL DEFAULT nanoid(16),
				entity_type TEXT NOT NULL CHECK (entity_type IN ('char', 'npc', 'location')),
				entity_id INT NOT NULL,
				game_id INT NOT NULL REFERENCES game(id),
				uploaded_by INT NOT NULL REFERENCES player(id),
				source_type TEXT NOT NULL CHECK (source_type IN ('uploaded', 'external')),
				storage_key TEXT,
				thumb_key TEXT,
				external_url TEXT,
				content_type TEXT,
				byte_size BIGINT NOT NULL DEFAULT 0,
				width INT,
				height INT,
				checksum TEXT,
				is_main BOOLEAN NOT NULL DEFAULT false,
				status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'complete', 'deleted')),
				created TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
				deleted TIMESTAMPTZ,
				CHECK (
					(source_type = 'external' AND external_url IS NOT NULL AND storage_key IS NULL AND thumb_key IS NULL)
					OR (source_type = 'uploaded' AND external_url IS NULL)
				)
			)`,
			`CREATE UNIQUE INDEX IF NOT EXISTS image_single_main_idx ON image(entity_type, entity_id) WHERE is_main AND deleted IS NULL`,
			`CREATE INDEX IF NOT EXISTS image_entity_idx ON image(entity_type, entity_id, created) WHERE deleted IS NULL`,
			`CREATE INDEX IF NOT EXISTS image_game_idx ON image(game_id) WHERE deleted IS NULL`,
			`CREATE TABLE IF NOT EXISTS game_image_quota (
				game_id INT PRIMARY KEY REFERENCES game(id) ON DELETE CASCADE,
				max_bytes BIGINT NOT NULL DEFAULT 0 CHECK (max_bytes >= 0),
				used_bytes BIGINT NOT NULL DEFAULT 0 CHECK (used_bytes >= 0),
				reserved_bytes BIGINT NOT NULL DEFAULT 0 CHECK (reserved_bytes >= 0),
				max_file_bytes BIGINT NOT NULL DEFAULT 5242880 CHECK (max_file_bytes > 0),
				max_images INT NOT NULL DEFAULT 0 CHECK (max_images >= 0)
			)`,
			`INSERT INTO game_image_quota (game_id) SELECT id FROM game ON CONFLICT (game_id) DO NOTHING`,
		}
		for _, statement := range statements {
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("image migration failed: %w", err)
			}
		}

		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migration (name) VALUES (?)", imageMigration); err != nil {
			return fmt.Errorf("failed to record image migration: %w", err)
		}
		return nil
	})
}

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
