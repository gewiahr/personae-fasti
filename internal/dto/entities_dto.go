package dto

/* Char */

type CharBrief struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`

	PlayerExtID string `json:"playerExt"`
	GameExtID   string `json:"gameExt"`
	Hidden      bool   `json:"hidden"`
}

type CharFull struct {
	ID          int    `json:"id"`
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
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
}

/* NPC */

type NPCBrief struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`

	GameExtID string `json:"gameExt"`
	Hidden    bool   `json:"hidden"`
}

type NPCFull struct {
	ID          int    `json:"id"`
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
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
}

/* Location */

type LocationBrief struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`

	ParentID int `json:"pid"`

	GameExtID string `json:"gameExt"`
	Hidden    bool   `json:"hidden"`
}

type LocationFull struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	ParentID int `json:"pid"`

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
	ParentID    int    `json:"pid"`
	Hidden      bool   `json:"hidden"`
}

type LocationUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ParentID    int    `json:"pid"`
	Hidden      bool   `json:"hidden"`
}

/* Suggestion */

type Suggestion struct {
	ID       int    `json:"id"`
	StringID string `bun:"sid" json:"sid"`
	Type     string `json:"type"`
	//TypeName string `json:"typeName"`
	Name string `json:"name"`
	// Hidden bool   `json:"hidden"`
	Secret bool `json:"secret"`
}
