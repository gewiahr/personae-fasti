package reqData

type RecordInsert struct {
	Text     string `json:"text"`
	Hidden   bool   `json:"hidden"`
	PlayerID int    `json:"-"`
	GameID   int    `json:"-"`
}

type RecordUpdate struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Hidden bool   `json:"hidden"`
}

type CharCreate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CharUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type NPCCreate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type NPCUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type LocationCreate struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type LocationUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GameChange struct {
	GameID int `json:"gameID"`
}
