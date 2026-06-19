package dto

import "time"

type RecordFull struct {
	ID   int    `json:"id"`
	Text string `json:"text"`

	PlayerExtID string `json:"playerExt"`
	GameExtID   string `json:"gameExt"`
	Hidden      bool   `json:"hidden"`

	QuestID int `json:"questID"`

	Created *time.Time `json:"created"`
	Updated *time.Time `json:"updated"`
	Deleted *time.Time `json:"deleted,omitempty"`
}

type RecordInsert struct {
	Text    string `json:"text"`
	Hidden  bool   `json:"hidden"`
	QuestID int    `json:"questID"`
}

type RecordUpdate struct {
	ID      int    `json:"id"`
	Text    string `json:"text"`
	Hidden  bool   `json:"hidden"`
	QuestID int    `json:"questID"`
}
