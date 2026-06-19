package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Char struct {
	bun.BaseModel `bun:"table:char"`

	ID    int    `bun:"id,pk,autoincrement"`
	Name  string `bun:"name,notnull"`
	Title string `bun:"title"`

	Description string `bun:"description"`

	PlayerID int     `bun:"player_id"`
	Player   *Player `bun:"rel:belongs-to,join:player_id=id"`
	GameID   int     `bun:"game_id"`
	Game     *Game   `bun:"rel:belongs-to,join:game_id=id"`
	HiddenBy int     `bun:"hidden_by,default:0"`

	Records []Record `bun:"m2m:records_chars,join:Char=Record"`

	Created *time.Time `bun:"created,default:current_timestamp"`
	Deleted *time.Time `bun:"deleted,default:null"`
}

type NPC struct {
	bun.BaseModel `bun:"table:npc"`

	ID    int    `bun:"id,pk,autoincrement"`
	Name  string `bun:"name,notnull"`
	Title string `bun:"title"`

	Description string `bun:"description"`

	GameID  int      `bun:"game_id"`
	Game    *Game    `bun:"rel:belongs-to,join:game_id=id"`
	Records []Record `bun:"m2m:records_npcs,join:NPC=Record"`

	CreatedByID int     `bun:"created_by_id"`
	CreatedBy   *Player `bun:"rel:belongs-to,join:created_by_id=id"`
	HiddenBy    int     `bun:"hidden_by,default:0"`

	Created *time.Time `bun:"created,default:current_timestamp"`
	Deleted *time.Time `bun:"deleted,default:null"`
}

type Location struct {
	bun.BaseModel `bun:"table:location"`

	ID          int    `bun:"id,pk,autoincrement"`
	Name        string `bun:"name,notnull"`
	Title       string `bun:"title"`
	Description string `bun:"description"`

	GameID   int       `bun:"game_id"`
	Game     *Game     `bun:"rel:belongs-to,join:game_id=id"`
	ParentID int       `bun:"pid,nullzero"`
	Parent   *Location `bun:"rel:belongs-to,join:pid=id"`
	Records  []Record  `bun:"m2m:records_locations,join:Location=Record"`

	CreatedByID int     `bun:"created_by_id"`
	CreatedBy   *Player `bun:"rel:belongs-to,join:created_by_id=id"`
	HiddenBy    int     `bun:"hidden_by,default:0"`

	Created *time.Time `bun:"created,default:current_timestamp"`
	Deleted *time.Time `bun:"deleted,default:null"`
}
