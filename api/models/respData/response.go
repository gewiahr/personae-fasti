package respData

import "personae-fasti/data"

type LoginInfo struct {
	AccessKey   string     `json:"accesskey"`
	Player      PlayerInfo `json:"player"`
	CurrentGame GameInfo   `json:"currentGame"`
}

type PlayerInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type GameInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type CharInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`

	PlayerID int `json:"playerID"`
	GameID   int `json:"gameID"`
}

type GameRecords struct {
	Records     []data.Record `json:"records"`
	Players     []PlayerInfo  `json:"players"`
	CurrentGame GameInfo      `json:"currentGame"`
}

type GameChars struct {
	Chars       []CharInfo   `json:"chars"`
	Players     []PlayerInfo `json:"players"`
	CurrentGame GameInfo     `json:"currentGame"`
}

type CharPage struct {
	Char    CharFullInfo  `json:"char"`
	Records []data.Record `json:"records"`
}

type CharFullInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	PlayerID int `json:"playerID"`
	GameID   int `json:"gameID"`
}

type SuggestionData struct {
	Suggestions []data.Suggestion `json:"entities"`
}
