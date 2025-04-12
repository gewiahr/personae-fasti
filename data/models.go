package data

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

func (s *Storage) InitTables() {

	s.db.RegisterModel((*PlayerGame)(nil))

	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Game)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Player)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Telegram)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Char)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*NPC)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Location)(nil)).Exec(context.Background())
	//_, _ = s.db.NewCreateTable().IfNotExists().Model((*Record)(nil)).Exec(context.Background())

	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Log)(nil)).Exec(context.Background())

}

type Game struct {
	bun.BaseModel `bun:"table:game"`

	ID   int    `bun:"id,pk,autoincrement"`
	Name string `bun:"name,notnull"`

	Players []Player `bun:"m2m:players_games,join:Game=Player"`
	Chars   []Char   `bun:"rel:has-many,join:id=game_id"`

	NPCs      []NPC      `bun:"rel:has-many,join:id=game_id"`
	Locations []Location `bun:"rel:has-many,join:id=game_id"`

	//Records []Record `bun:"rel:has-many,join:id=game_id"`

	//Comments []Comment `bun:"rel:has-many,join:id=game_id"`
}

type Player struct {
	bun.BaseModel `bun:"table:player"`

	ID         int       `bun:"id,pk,autoincrement"`
	Username   string    `bun:"username,unique,notnull"`
	AccessKey  string    `bun:"accesskey,notnull"`
	TelegramID int64     `bun:"telegram_id"`
	Telegram   *Telegram `bun:"rel:belongs-to,join:telegram_id=id"`

	Chars []Char `bun:"rel:has-many,join:id=player_id"`
	Games []Game `bun:"m2m:players_games,join:Player=Game"`

	//Comments []Comment `bun:"rel:has-many,join:id=player_id"`

	TimeRegistered time.Time `bun:"timeRegistered,nullzero,notnull,default:current_timestamp"`
	TimeLastAction time.Time `bun:"timeLastAction,nullzero,notnull,default:current_timestamp"`

	CurrentGameID int   `bun:"current_game_id"`
	CurrentGame   *Game `bun:"rel:belongs-to,join=current_game_id=id"`
}

type Telegram struct {
	bun.BaseModel `bun:"table:telegram"`

	ID       int64  `bun:"id,pk"`
	Username string `bun:"username,notnull"`
	Lang     string `bun:"lang,default:en"`
}

type Char struct {
	bun.BaseModel `bun:"table:char"`

	ID          int    `bun:"id,pk"`
	Name        string `bun:"name,notnull"`
	Description string `bun:"description"`

	PlayerID int `bun:"player_id"`
	GameID   int `bun:"game_id"`
}

type PlayerGame struct {
	bun.BaseModel `bun:"players_games"`

	PlayerID int     `bun:"player_id,pk"`
	Player   *Player `bun:"rel:belongs-to,join:player_id=id"`
	GameID   int     `bun:"game_id,pk"`
	Game     *Game   `bun:"rel:belongs-to,join:game_id=id"`
}

type NPC struct {
	bun.BaseModel `bun:"table:npc"`

	ID          int    `bun:"id,pk"`
	Name        string `bun:"name,notnull"`
	Description string `bun:"description"`

	GameID int `bun:"game_id"`
}

type Location struct {
	bun.BaseModel `bun:"table:location"`

	ID          int    `bun:"id,pk"`
	ParentID    int    `bun:"pid"` // TODO: Change to relations
	Name        string `bun:"name,notnull"`
	Description string `bun:"description"`

	GameID int `bun:"game_id"`
}
