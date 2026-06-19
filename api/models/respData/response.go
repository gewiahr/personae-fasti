package respData

// type LoginInfo struct {
// 	AccessKey   string       `json:"accesskey"`
// 	Player      PlayerInfo   `json:"player"`
// 	CurrentGame GameFullInfo `json:"currentGame"`
// }

// type PlayerSettings struct {
// 	CurrentGame   GameFullInfo `json:"currentGame"`
// 	PlayerGames   []GameInfo   `json:"playerGames"`
// 	PlayerInvites []GameInfo   `json:"playerInvites"`
// }

// func FormPlayerSettings(playerGames, playerInvites []domain.Game, currentGame *domain.Game) *PlayerSettings {
// 	var playerGameInfo []GameInfo
// 	for _, game := range playerGames {
// 		playerGameInfo = append(playerGameInfo, *GameToGameInfo(&game))
// 	}

// 	var playerInvitesInfo []GameInfo
// 	for _, invite := range playerInvites {
// 		playerInvitesInfo = append(playerInvitesInfo, *GameToGameInfo(&invite))
// 	}

// 	return &PlayerSettings{
// 		CurrentGame:   *GameToGameFullInfo(currentGame),
// 		PlayerGames:   playerGameInfo,
// 		PlayerInvites: playerInvitesInfo,
// 	}
// }

// type PlayerInfo struct {
// 	ID       int    `json:"id"`
// 	Username string `json:"username"`
// }

// type SuggestionData struct {
// 	Suggestions []dto.Suggestion `json:"entities"`
// }
