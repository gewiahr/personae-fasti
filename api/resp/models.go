package resp

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

type GameRecords struct {
	Records     []data.Record `json:"records"`
	Players     []PlayerInfo  `json:"players"`
	CurrentGame GameInfo      `json:"currentGame"`
}
