package data

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"time"

// 	"personae-fasti/configs"

// 	"github.com/uptrace/bun"
// 	"github.com/uptrace/bun/dialect/pgdialect"
// 	"github.com/uptrace/bun/driver/pgdriver"
// 	"github.com/uptrace/bun/extra/bundebug"
// )

// type Storage struct {
// 	db *bun.DB
// }

// type Log struct {
// 	bun.BaseModel `bun:"table:log_api"`

// 	ID       int       `json:"id" bun:",pk,autoincrement"`
// 	Time     time.Time `json:"time" bun:",notnull,default:now()"`
// 	User     int       `json:"user" bun:",notnull"`
// 	URI      string    `json:"uri" bun:",notnull"`
// 	Method   string    `json:"method" bun:",notnull"`
// 	Request  string    `json:"request" bun:",notnull"`
// 	Response string    `json:"response" bun:",notnull"`
// 	HTTPCode int       `json:"httpCode" bun:",notnull"`
// 	Error    string    `json:"error" bun:""`
// }

// func NewStorage(c *configs.Main) *Storage {

// }

// func (s *Storage) Log(log *Log, ctx context.Context) {
// 	s.db.NewInsert().Model(log).Exec(ctx)
// }

// func (s *Storage) LogAuthorized(log *Log, ctx context.Context) {
// 	s.db.NewInsert().Model(log).Exec(ctx)
// }
