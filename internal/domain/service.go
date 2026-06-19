package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type ApiLog struct {
	bun.BaseModel `bun:"table:log_api"`

	ID       int64 `bun:",pk,autoincrement"`
	PlayerID int   `bun:"player_id,notnull"`

	URI     string `bun:",notnull"`
	Method  string `bun:",notnull"`
	Request string `bun:",notnull"`

	Response string `bun:",notnull"`
	Code     int    `bun:",notnull"`
	Error    string `bun:""`
	Time     int64  `bun:""`

	Created time.Time `bun:",notnull,default:now()"`
}

type ServiceFeedback struct {
	bun.BaseModel `bun:"table:service_feedback"`

	ID       int    `bun:"id,pk,autoincrement"`
	Type     string `bun:"type,notnull,default:''"`
	Text     string `bun:"text,notnull"`
	Response string `bun:"response"`

	PlayerID int     `bun:"player_id"`
	Player   *Player `bun:"rel:belongs-to,join:player_id=id"`
	GameID   int     `bun:"game_id,notnull"`
	Game     *Game   `bun:"rel:belongs-to,join:game_id=id"`

	Created *time.Time `bun:"created,nullzero,notnull,default:current_timestamp"`
	Updated *time.Time `bun:"updated,nullzero,notnull,default:current_timestamp"`
}
