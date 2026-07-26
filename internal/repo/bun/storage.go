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
	storage.Migrate(ctx)

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
		panic(fmt.Errorf("failed to create nanoid function: %w", err))
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
        $$ LANGUAGE plpgsql STABLE PARALLEL SAFE;
    `

	if _, err := db.ExecContext(context.Background(), `CREATE EXTENSION IF NOT EXISTS pgcrypto`); err != nil {
		return err
	}

	if _, err := db.ExecContext(context.Background(), nanoidFunction); err != nil {
		return err
	}

	return nil
}

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
