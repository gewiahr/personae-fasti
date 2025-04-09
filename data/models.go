package data

import (
	"context"
)

func (s *Storage) InitTables() {

	// s.db.RegisterModel((*MODEL)(nil))

	_, _ = s.db.NewCreateTable().IfNotExists().Model((*Log)(nil)).Exec(context.Background())

}

// type Game struct {
// 	bun.BaseModel `bun:"table:game"`

// 	ID       int      `bun:"id,pk,autoincrement"`
// 	Name     string   `bun:"name"`
// 	Entities []Entity `bun:"rel:has-many,join:id=game_id"`
// }
