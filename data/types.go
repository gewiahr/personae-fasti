package data

type Suggestion struct {
	ID       int    `bun:"id" json:"id"`
	StringID string `bun:"sid" json:"sid"`
	Type     string `bun:"type" json:"type"`
	//TypeName string `bun:"typeName" json:"typeName"`
	Name   string `bun:"name" json:"name"`
	Hidden bool   `bun:"hidden" json:"hidden"`
}
