package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"personae-fasti/opt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Storage struct {
	db *bun.DB
}

type Log struct {
	bun.BaseModel `bun:"table:log_api"`

	ID       int       `json:"id" bun:",pk,autoincrement"`
	Time     time.Time `json:"time" bun:",notnull,default:now()"`
	User     int       `json:"user" bun:",notnull"`
	URI      string    `json:"uri" bun:",notnull"`
	Method   string    `json:"method" bun:",notnull"`
	Request  string    `json:"request" bun:",notnull"`
	Response string    `json:"response" bun:",notnull"`
	HTTPCode int       `json:"httpCode" bun:",notnull"`
	Error    string    `json:"error" bun:""`
}

func NewStorage(c *opt.Conf) *Storage {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name)
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(pgdb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	storage := &Storage{
		db: db,
	}
	storage.InitTables()

	return storage

}

func (s *Storage) Log(log *Log, ctx context.Context) {
	s.db.NewInsert().Model(log).Exec(ctx)
}
