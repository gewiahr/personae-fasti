package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Record struct {
	bun.BaseModel `bun:"table:record"`

	ID   int    `bun:"id,pk,autoincrement"`
	Text string `bun:"text,notnull"`

	Chars     []Char     `bun:"m2m:records_chars,join:Record=Char"`
	NPCs      []NPC      `bun:"m2m:records_npcs,join:Record=NPC"`
	Locations []Location `bun:"m2m:records_locations,join:Record=Location"`

	PlayerID int     `bun:"player_id"`
	Player   *Player `bun:"rel:belongs-to,join:player_id=id"`
	GameID   int     `bun:"game_id"`
	Game     *Game   `bun:"rel:belongs-to,join:game_id=id"`
	HiddenBy int     `bun:"hidden_by,default:0"`

	QuestID int    `bun:"quest_id"`
	Quest   *Quest `bun:"rel:belongs-to,join:quest_id=id"`

	Created *time.Time `bun:"created,nullzero,notnull,default:current_timestamp"`
	Updated *time.Time `bun:"updated,nullzero,notnull,default:current_timestamp"`
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
