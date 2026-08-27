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
	Invites  []PlayerBrief  `json:"invites"`
}

type GameSettings struct {
	AllowAllEditRecords bool `json:"allowAllEditRecords"`
}

type SessionBrief struct {
	Number  int        `json:"number"`
	Name    string     `json:"name"`
	EndTime *time.Time `json:"endTime"`
}

type SessionUpdate struct {
	Number    int       `json:"number"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"startTime"`
}

type GameInvite struct {
	PlayerExt  string `json:"playerExt"`
	GameExt    string `json:"gameExt"`
	GameTitle  string `json:"gameTitle"`
	InviteCode string `json:"inviteCode"`
}

type GameCreate struct {
	Title string `json:"title"`
}

type GameUpdate struct {
	Ext   string `json:"ext"`
	Title string `json:"title"`
	GMExt string `json:"gmExt"`
}

type GamePage struct {
	Game GameFull `json:"game"`
}

type PlayerSettingsResponse struct {
	CurrentGame   *GameFull    `json:"currentGame"`
	PlayerGames   []GameBrief  `json:"playerGames"`
	PlayerInvites []GameInvite `json:"playerInvites"`
}

type GameSettingsUpdate struct {
	GameExt             string `json:"gameExt"`
	AllowAllEditRecords bool   `json:"allowAllEditRecords"`
}
