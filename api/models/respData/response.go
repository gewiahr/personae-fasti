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

type PlayerSettings struct {
	CurrentGame GameInfo   `json:"currentGame"`
	PlayerGames []GameInfo `json:"playerGames"`
}

func FormPlayerSettings(playerGames []data.Game, currentGame data.Game) *PlayerSettings {
	var playerGameInfo []GameInfo
	for _, game := range playerGames {
		playerGameInfo = append(playerGameInfo, *GameToGameInfo(&game))
	}

	return &PlayerSettings{
		CurrentGame: *GameToGameInfo(&currentGame),
		PlayerGames: playerGameInfo,
	}
}

type GameInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	GMID  int    `json:"gmID"`
}

type GameRecords struct {
	Records     []data.Record  `json:"records"`
	Sessions    []data.Session `json:"sessions"`
	Players     []PlayerInfo   `json:"players"`
	CurrentGame GameInfo       `json:"currentGame"`
}

func FormGameRecords(p *data.Player, rs []data.Record, ps []data.Player, ss []data.Session) *GameRecords {
	gameRecords := GameRecords{
		Records:  rs,
		Sessions: ss,
		Players:  PlayersToPlayersInfoArray(ps),
		CurrentGame: GameInfo{
			ID:    p.CurrentGame.ID,
			Title: p.CurrentGame.Name,
			GMID:  p.CurrentGame.GMID,
		},
	}

	return &gameRecords
}

type CharInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`

	PlayerID int `json:"playerID"`
	GameID   int `json:"gameID"`
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

type NPCInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`

	GameID int `json:"gameID"`
}

type GameNPCs struct {
	NPCs        []NPCInfo `json:"npcs"`
	CurrentGame GameInfo  `json:"currentGame"`
}

type NPCPage struct {
	NPC     NPCFullInfo   `json:"npc"`
	Records []data.Record `json:"records"`
}

type NPCFullInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	GameID int `json:"gameID"`
}

type LocationInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`

	GameID int `json:"gameID"`
}

type GameLocations struct {
	Locations   []LocationInfo `json:"locations"`
	CurrentGame GameInfo       `json:"currentGame"`
}

type LocationPage struct {
	Location LocationFullInfo `json:"location"`
	Records  []data.Record    `json:"records"`
}

type LocationFullInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	GameID int `json:"gameID"`
}

type SuggestionData struct {
	Suggestions []data.Suggestion `json:"entities"`
}
