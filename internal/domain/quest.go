package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Quest struct {
	bun.BaseModel `bun:"quest"`

	ID    int    `bun:"id,pk,autoincrement"`
	ExtID string `bun:"ext,unique,notnull,type:varchar(16),default:nanoid(12)"`

	GameID int   `bun:"game_id,notnull"`
	Game   *Game `bun:"rel:belongs-to,join:game_id=id"`

	Name        string `bun:",notnull,default:''"`
	Title       string `bun:",notnull,default:''"`
	Description string `bun:",notnull,default:''"`

	Records []Record `bun:"rel:has-many,join:id=quest_id"`

	ParentID int    `bun:"parent_id"`
	Parent   *Quest `bun:"rel:belongs-to,join:parent_id=id"`
	ChildID  int    `bun:"child_id"`
	Child    *Quest `bun:"rel:belongs-to,join:child_id=id"`
	HeadID   int    `bun:"head_id"`
	Head     *Quest `bun:"rel:belongs-to,join:head_id=id"`

	Tasks []QuestTask `bun:"rel:has-many,join:id=quest_id"`

	Successful bool `bun:"successful,default:false"`

	HiddenBy int `bun:"hidden_by,default:0"`

	Created  *time.Time `bun:"created,default:current_timestamp"`
	Deleted  *time.Time `bun:"deleted,default:null"`
	Finished *time.Time `bun:"finished,default:null"`
}

type QuestTaskType int

const (
	Binary QuestTaskType = iota
	Decimal
)

type QuestTask struct {
	bun.BaseModel `bun:"quest_task"`

	ID int `bun:"id,pk,autoincrement"`

	GameID  int    `bun:"game_id,notnull"`
	Game    *Game  `bun:"rel:belongs-to,join:game_id=id"`
	QuestID int    `bun:"quest_id,notnull"`
	Quest   *Quest `bun:"rel:belongs-to,join:quest_id=id"`

	Name        string        `bun:",notnull,default:''"`
	Description string        `bun:",notnull,default:''"`
	Type        QuestTaskType `bun:",default:0"`
	Capacity    int           `bun:"capacity,default:0"`
	Current     int           `bun:"current,default:0"`

	HiddenBy int `bun:"hidden_by,default:0"`

	Finished *time.Time `bun:"finished,default:null"`
}
