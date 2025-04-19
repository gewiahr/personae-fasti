package reqData

type RecordInsert struct {
	Text     string `json:"text"`
	PlayerID int    `json:"-"`
	GameID   int    `json:"-"`
}

type CharUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
