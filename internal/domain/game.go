package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Game struct {
	bun.BaseModel `bun:"table:game"`

	ID    int    `bun:"id,pk,autoincrement"`
	ExtID string `bun:"ext,unique,notnull,type:varchar(16)"`
	Name  string `bun:"name,notnull"`

	GMID int     `bun:"gm_id"`
	GM   *Player `bun:"rel:belongs-to,join:gm_id=id"`

	Sessions []Session `bun:"rel:has-many,join:id=game_id"`

	Players []Player `bun:"m2m:players_games,join:Game=Player"`
	Invites []Player `bun:"m2m:game_invites,join:Game=Player"`

	Chars     []Char     `bun:"rel:has-many,join:id=game_id"`
	NPCs      []NPC      `bun:"rel:has-many,join:id=game_id"`
	Locations []Location `bun:"rel:has-many,join:id=game_id"`

	Records []Record `bun:"rel:has-many,join:id=game_id"`
	Quests  []Quest  `bun:"rel:has-many,join:id=game_id"`

	Settings *GameSettings `bun:"rel:has-one,join:id=game_id"`

	Created *time.Time `bun:"created,default:current_timestamp"`
	Deleted *time.Time `bun:"deleted,default:null"`
}

type GameSettings struct {
	bun.BaseModel `bun:"table:game_settings"`

	GameID int   `bun:"game_id,pk"`
	Game   *Game `bun:"rel:belongs-to,join:game_id=id"`

	AllowAllEditRecords bool `bun:"allow_all_edit_records,default:false"`
}

type PlayerGame struct {
	bun.BaseModel `bun:"players_games"`

	PlayerID int     `bun:"player_id,pk"`
	Player   *Player `bun:"rel:belongs-to,join:player_id=id"`
	GameID   int     `bun:"game_id,pk"`
	Game     *Game   `bun:"rel:belongs-to,join:game_id=id"`
	// Status   PlayerGameStatus `bun:"status,default:1"`
}

// type PlayerGameStatus int8

// const (
// 	PlayerGameStatusBanned      PlayerGameStatus = -2
// 	PlayerGameStatusLeft        PlayerGameStatus = -1
// 	PlayerGameStatusInvited     PlayerGameStatus = 0
// 	PlayerGameStatusParticipant PlayerGameStatus = 1
// 	PlayerGameStatusSpectator   PlayerGameStatus = 2
// )

type GameInvite struct {
	bun.BaseModel `bun:"game_invites"`

	PlayerID int     `bun:"player_id,pk"`
	Player   *Player `bun:"rel:belongs-to,join:player_id=id"`
	GameID   int     `bun:"game_id,pk"`
	Game     *Game   `bun:"rel:belongs-to,join:game_id=id"`

	Code string `bun:"code,unique,notnull,type:varchar(16),default:nanoid(16)"`

	Created *time.Time `bun:"created,default:current_timestamp"`
}

type Session struct {
	bun.BaseModel `bun:"session"`

	ID int `bun:"id,pk,autoincrement"`

	GameID int   `bun:"game_id,notnull"`
	Game   *Game `bun:"rel:belongs-to,join:game_id=id"`

	Number int    `bun:"number,notnull"`
	Name   string `bun:",notnull,default:''"`

	EndTime *time.Time `bun:"end_time,nullzero"`
}
