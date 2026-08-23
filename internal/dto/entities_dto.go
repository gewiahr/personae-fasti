package dto

/* Char */

type CharBrief struct {
	ExtID string `json:"ext"`
	Name  string `json:"name"`
	Title string `json:"title"`

	PlayerExtID string `json:"playerExt"`
	GameExtID   string `json:"gameExt"`
	Hidden      bool   `json:"hidden"`
}

type CharFull struct {
	ExtID       string `json:"ext"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	PlayerExtID string `json:"playerExt"`
	GameExtID   string `json:"gameExt"`
	Hidden      bool   `json:"hidden"`
}

// type GameChars struct {
// 	Chars       []CharInfo   `json:"chars"`
// 	Players     []PlayerInfo `json:"players"`
// 	CurrentGame GameInfo     `json:"currentGame"`
// }

type CharPage struct {
	Char    CharFull     `json:"char"`
	Records []RecordFull `json:"records"`
}

type CharCreate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
}

type CharUpdate struct {
	ExtID       string `json:"ext"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
}

/* NPC */

type NPCBrief struct {
	ExtID string `json:"ext"`
	Name  string `json:"name"`
	Title string `json:"title"`

	GameExtID string `json:"gameExt"`
	Hidden    bool   `json:"hidden"`
}

type NPCFull struct {
	ExtID       string `json:"ext"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	GameExtID string `json:"gameExt"`
	Hidden    bool   `json:"hidden"`
}

// type GameNPCs struct {
// 	NPCs        []NPCInfo `json:"npcs"`
// 	CurrentGame GameInfo  `json:"currentGame"`
// }

type NPCPage struct {
	NPC     NPCFull      `json:"npc"`
	Records []RecordFull `json:"records"`
}

type NPCCreate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
}

type NPCUpdate struct {
	ExtID       string `json:"ext"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
}

/* Location */

type LocationBrief struct {
	ExtID string `json:"ext"`
	Name  string `json:"name"`
	Title string `json:"title"`

	GameExtID string `json:"gameExt"`
	Hidden    bool   `json:"hidden"`
}

type LocationFull struct {
	ExtID       string `json:"ext"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	ParentExtID string `json:"parentExt"`

	GameExtID string `json:"gameExt"`
	Hidden    bool   `json:"hidden"`
}

// type GameLocations struct {
// 	Locations   []LocationInfo `json:"locations"`
// 	CurrentGame GameInfo       `json:"currentGame"`
// }

type LocationPage struct {
	Location LocationFull    `json:"location"`
	Records  []RecordFull    `json:"records"`
	Parent   *LocationBrief  `json:"parent"`
	Includes []LocationBrief `json:"includes"`
}

type LocationCreate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ParentExtID string `json:"parentExt"`
	Hidden      bool   `json:"hidden"`
}

type LocationUpdate struct {
	ExtID       string `json:"ext"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ParentExtID string `json:"parentExt"`
	ParentID    int    `json:"-"`
	Hidden      bool   `json:"hidden"`
}

/* Suggestion */

type Suggestion struct {
	ExtID    string `bun:"ext" json:"ext"`
	StringID string `bun:"sid" json:"sid"`
	Type     string `json:"type"`
	//TypeName string `json:"typeName"`
	Name string `json:"name"`
	// Hidden bool   `json:"hidden"`
	Secret bool `json:"secret"`
}
