package reqData

type RecordInsert struct {
	Text     string `json:"text"`
	PlayerID int    `json:"-"`
	GameID   int    `json:"-"`
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
