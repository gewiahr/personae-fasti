package data

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

func (s *Storage) InitTables() {

	s.db.RegisterModel((*PlayerGame)(nil))
	s.db.RegisterModel((*RecordChar)(nil))
	s.db.RegisterModel((*RecordNPC)(nil))
	s.db.RegisterModel((*RecordLocation)(nil))

	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Game)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Player)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Telegram)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Char)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*NPC)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Location)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Record)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Session)(nil)).Exec(context.Background())
	//_, _ = s.db.NewCreateTable().IfNotExists().Model((*Mention)(nil)).Exec(context.Background())

	_, _ = s.db.NewCreateTable().IfNotExists().Model((*PlayerGame)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*RecordChar)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*RecordNPC)(nil)).Exec(context.Background())
	_, _ = s.db.NewCreateTable().IfNotExists().Model((*RecordLocation)(nil)).Exec(context.Background())

	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Log)(nil)).Exec(context.Background())

}

type Game struct {
	bun.BaseModel `bun:"table:game"`

	ID   int    `bun:"id,pk,autoincrement"`
	Name string `bun:"name,notnull"`

	GMID int64   `bun:"gm_id"`
	GM   *Player `bun:"rel:belongs-to,join:gm_id=id"`

	Sessions []Session `bun:"rel:has-many,join:id=game_id"`

	Players []Player `bun:"m2m:players_games,join:Game=Player"`
	Chars   []Char   `bun:"rel:has-many,join:id=game_id"`

	NPCs      []NPC      `bun:"rel:has-many,join:id=game_id"`
	Locations []Location `bun:"rel:has-many,join:id=game_id"`

	Records []Record `bun:"rel:has-many,join:id=game_id"`

	Created *time.Time `bun:"created,default:current_timestamp"`
	Deleted *time.Time `bun:"deleted,default:null"`
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

	Records []Record `bun:"rel:has-many,join:id=game_id"`

	CurrentGameID int   `bun:"current_game_id"`
	CurrentGame   *Game `bun:"rel:belongs-to,join:current_game_id=id"`

	Registered *time.Time `bun:"registeredTime,nullzero,notnull,default:current_timestamp"`
	LastAction *time.Time `bun:"lastActionTime,nullzero,notnull,default:current_timestamp"`
	Deleted    *time.Time `bun:"deleted,default:null"`
}

type Telegram struct {
	bun.BaseModel `bun:"table:telegram"`

	ID       int64  `bun:"id,pk"`
	Username string `bun:"username,notnull"`
	Lang     string `bun:"lang,default:'en'"`
}

type Char struct {
	bun.BaseModel `bun:"table:char"`

	ID    int    `bun:"id,pk,autoincrement"`
	Name  string `bun:"name,notnull"`
	Title string `bun:"title"`

	Description string `bun:"description"`

	PlayerID int `bun:"player_id"`
	GameID   int `bun:"game_id"`

	Records []Record `bun:"m2m:records_chars,join:Char=Record"`

	Created *time.Time `bun:"created,default:current_timestamp"`
	Deleted *time.Time `bun:"deleted,default:null"`
}

type PlayerGame struct {
	bun.BaseModel `bun:"players_games"`

	PlayerID int     `bun:"player_id,pk,autoincrement"`
	Player   *Player `bun:"rel:belongs-to,join:player_id=id"`
	GameID   int     `bun:"game_id,pk"`
	Game     *Game   `bun:"rel:belongs-to,join:game_id=id"`
}

type NPC struct {
	bun.BaseModel `bun:"table:npc"`

	ID    int    `bun:"id,pk,autoincrement"`
	Name  string `bun:"name,notnull"`
	Title string `bun:"title"`

	Description string `bun:"description"`

	GameID  int      `bun:"game_id"`
	Records []Record `bun:"m2m:records_npcs,join:NPC=Record"`

	CreatedByID int     `bun:"created_by_id"`
	CreatedBy   *Player `bun:"rel:belongs-to,join:created_by_id=id"`

	Created *time.Time `bun:"created,default:current_timestamp"`
	Deleted *time.Time `bun:"deleted,default:null"`
}

type Location struct {
	bun.BaseModel `bun:"table:location"`

	ID          int    `bun:"id,pk,autoincrement"`
	ParentID    int    `bun:"pid"` // TODO: Change to relations
	Name        string `bun:"name,notnull"`
	Title       string `bun:"title"`
	Description string `bun:"description"`

	GameID  int      `bun:"game_id"`
	Records []Record `bun:"m2m:records_locations,join:Location=Record"`

	CreatedByID int     `bun:"created_by_id"`
	CreatedBy   *Player `bun:"rel:belongs-to,join:created_by_id=id"`

	Created *time.Time `bun:"created,default:current_timestamp"`
	Deleted *time.Time `bun:"deleted,default:null"`
}

type Record struct {
	bun.BaseModel `bun:"table:record"`

	ID   int    `bun:"id,pk,autoincrement" json:"id"`
	Text string `bun:"text,notnull" json:"text"`

	Chars     []Char     `bun:"m2m:records_chars,join:Record=Char" json:"chars,omitempty"`
	NPCs      []NPC      `bun:"m2m:records_npcs,join:Record=NPC" json:"npcs,omitempty"`
	Locations []Location `bun:"m2m:records_locations,join:Record=Location" json:"locations,omitempty"`

	PlayerID int `bun:"player_id" json:"playerID"`
	GameID   int `bun:"game_id" json:"gameID"`

	Created *time.Time `bun:"created,nullzero,notnull,default:current_timestamp" json:"created"`
	Updated *time.Time `bun:"updated,nullzero,notnull,default:current_timestamp" json:"updated"`
	Deleted *time.Time `bun:"deleted,default:null"`
}

type RecordChar struct {
	bun.BaseModel `bun:"records_chars"`

	RecordID int     `bun:"record_id,pk,autoincrement"`
	Record   *Record `bun:"rel:belongs-to,join:record_id=id"`
	CharID   int     `bun:"char_id,pk"`
	Char     *Char   `bun:"rel:belongs-to,join:char_id=id"`
}

type RecordNPC struct {
	bun.BaseModel `bun:"records_npcs"`

	RecordID int     `bun:"record_id,pk,autoincrement"`
	Record   *Record `bun:"rel:belongs-to,join:record_id=id"`
	NPCID    int     `bun:"npc_id,pk"`
	NPC      *NPC    `bun:"rel:belongs-to,join:npc_id=id"`
}

type RecordLocation struct {
	bun.BaseModel `bun:"records_locations"`

	RecordID   int       `bun:"record_id,pk,autoincrement"`
	Record     *Record   `bun:"rel:belongs-to,join:record_id=id"`
	LocationID int       `bun:"location_id,pk"`
	Location   *Location `bun:"rel:belongs-to,join:location_id=id"`
}

type Session struct {
	bun.BaseModel `bun:"session"`

	ID int `bun:"id,pk,autoincrement" json:"id"`

	GameID int   `bun:"game_id,notnull"`
	Game   *Game `bun:"rel:belongs-to,join:game_id=id"`

	Number int    `bun:"number,notnull" json:"number"`
	Name   string `bun:",notnull,default:''" json:"name"`

	EndTime *time.Time `bun:"end_time,nullzero" json:"endTime"`
}
