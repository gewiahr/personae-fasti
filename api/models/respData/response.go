package respData

import (
	"personae-fasti/data"
	"time"
)

type LoginInfo struct {
	AccessKey   string       `json:"accesskey"`
	Player      PlayerInfo   `json:"player"`
	CurrentGame GameFullInfo `json:"currentGame"`
}

type PlayerInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type PlayerSettings struct {
	CurrentGame GameFullInfo `json:"currentGame"`
	PlayerGames []GameInfo   `json:"playerGames"`
}

func FormPlayerSettings(playerGames []data.Game, currentGame *data.Game) *PlayerSettings {
	var playerGameInfo []GameInfo
	for _, game := range playerGames {
		playerGameInfo = append(playerGameInfo, *GameToGameInfo(&game))
	}

	return &PlayerSettings{
		CurrentGame: *GameToGameFullInfo(currentGame),
		PlayerGames: playerGameInfo,
	}
}

type GameInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	GMID  int    `json:"gmID"`
}

type GameFullInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	GMID  int    `json:"gmID"`

	Settings *GameSettings `json:"settings"`
	Sessions []SessionInfo `json:"sessions"`
}

type GameRecords struct {
	Records     []data.Record  `json:"records"`
	Sessions    []data.Session `json:"sessions"`
	Players     []PlayerInfo   `json:"players"`
	CurrentGame GameInfo       `json:"currentGame"`
}

type GameSettings struct {
	AllowAllEditRecords bool `json:"allowAllEditRecords"`
}

type SessionInfo struct {
	Number  int        `json:"number"`
	Name    string     `json:"name"`
	EndTime *time.Time `json:"endTime"`
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
	HiddenBy int `json:"hiddenBy"`
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
	HiddenBy int `json:"hiddenBy"`
}

type NPCInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`

	GameID   int `json:"gameID"`
	HiddenBy int `json:"hiddenBy"`
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

	GameID   int `json:"gameID"`
	HiddenBy int `json:"hiddenBy"`
}

type LocationInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`

	GameID   int `json:"gameID"`
	HiddenBy int `json:"hiddenBy"`
}

type GameLocations struct {
	Locations   []LocationInfo `json:"locations"`
	CurrentGame GameInfo       `json:"currentGame"`
}

type LocationPage struct {
	Location LocationFullInfo `json:"location"`
	Records  []data.Record    `json:"records"`
	Parent   *LocationInfo    `json:"parent"`
	Includes []LocationInfo   `json:"includes"`
}

type LocationFullInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	ParentID int `json:"pid"`

	GameID   int `json:"gameID"`
	HiddenBy int `json:"hiddenBy"`
}

type QuestInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	GameID     int  `json:"gameID"`
	HiddenBy   int  `json:"hiddenBy"`
	Successful bool `json:"successful"`
	Finished   bool `json:"finished"`
}

type GameQuests struct {
	Quests      []QuestInfo `json:"quests"`
	CurrentGame GameInfo    `json:"currentGame"`
}

type QuestPage struct {
	Quest   QuestFullInfo       `json:"quest"`
	Tasks   []QuestTaskFullInfo `json:"tasks"`
	Records []data.Record       `json:"records"`
}

type QuestFullInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`

	ParentID int `json:"parentID"`
	ChildID  int `json:"childID"`
	HeadID   int `json:"headID"`

	GameID     int  `json:"gameID"`
	HiddenBy   int  `json:"hiddenBy"`
	Successful bool `json:"successful"`
	Finished   bool `json:"finished"`
}

type QuestTaskFullInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	QuestID int `json:"questID"`

	Type     int  `json:"type"`
	Capacity int  `json:"capacity"`
	Current  int  `json:"current"`
	Finished bool `json:"finished"`

	GameID   int `json:"gameID"`
	HiddenBy int `json:"hiddenBy"`
}

type SuggestionData struct {
	Suggestions []data.Suggestion `json:"entities"`
}
