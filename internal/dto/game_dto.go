package dto

import "time"

type GameBrief struct {
	ExtID   string `json:"ext"`
	Title   string `json:"title"`
	GMExtID string `json:"gmExt"`
}

type GameFull struct {
	ExtID   string `json:"ext"`
	Title   string `json:"title"`
	GMExtID string `json:"gmExt"`

	Settings *GameSettings  `json:"settings"`
	Sessions []SessionBrief `json:"sessions"`
	Players  []PlayerBrief  `json:"players"`
}

type GameSettings struct {
	AllowAllEditRecords bool `json:"allowAllEditRecords"`
}

type SessionBrief struct {
	Number  int        `json:"number"`
	Name    string     `json:"name"`
	EndTime *time.Time `json:"endTime"`
}

// type GamePage struct {
// 	Game    GameFullInfo `json:"game"`
// 	Players []PlayerInfo `json:"players"`
// 	Invites []PlayerInfo `json:"invites"`
// }

// type GameRecords struct {
// 	Records     []domain.Record  `json:"records"`
// 	Sessions    []domain.Session `json:"sessions"`
// 	Players     []PlayerInfo     `json:"players"`
// 	CurrentGame GameInfo         `json:"currentGame"`
// }

type PlayerSettingsResponse struct {
	CurrentGame   *GameFull    `json:"currentGame"`
	PlayerGames   []GameBrief  `json:"playerGames"`
	PlayerInvites []GameBrief  `json:"playerInvites"`
}
