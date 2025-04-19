package reqData

type RecordInsert struct {
	Text     string `json:"text"`
	PlayerID int    `json:"-"`
	GameID   int    `json:"-"`
}
