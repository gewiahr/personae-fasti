package respData

import "personae-fasti/data"

func GameToGameInfo(game *data.Game) *GameInfo {
	return &GameInfo{
		ID:    game.ID,
		Title: game.Name,
	}
}

func PlayersToPlayersInfoArray(players []data.Player) []PlayerInfo {
	playerInfoArray := []PlayerInfo{}
	for _, player := range players {
		playerInfoArray = append(playerInfoArray, PlayerInfo{
			ID:       player.ID,
			Username: player.Username,
		})
	}

	return playerInfoArray
}

func CharToCharInfoArray(chars []data.Char) []CharInfo {
	charInfoArray := []CharInfo{}
	for _, char := range chars {
		charInfoArray = append(charInfoArray, CharInfo{
			ID:       char.ID,
			Name:     char.Name,
			Title:    char.Title,
			PlayerID: char.PlayerID,
			GameID:   char.GameID,
		})
	}

	return charInfoArray
}

func CharToCharFullInfo(char *data.Char) *CharFullInfo {
	return &CharFullInfo{
		ID:          char.ID,
		Name:        char.Name,
		Title:       char.Title,
		Description: char.Description,
		PlayerID:    char.PlayerID,
		GameID:      char.GameID,
	}
}

func NPCToNPCInfoArray(npcs []data.NPC) []NPCInfo {
	npcInfoArray := []NPCInfo{}
	for _, npc := range npcs {
		npcInfoArray = append(npcInfoArray, NPCInfo{
			ID:     npc.ID,
			Name:   npc.Name,
			Title:  npc.Title,
			GameID: npc.GameID,
		})
	}

	return npcInfoArray
}

func NPCToNPCFullInfo(npc *data.NPC) *NPCFullInfo {
	return &NPCFullInfo{
		ID:          npc.ID,
		Name:        npc.Name,
		Title:       npc.Title,
		Description: npc.Description,
		GameID:      npc.GameID,
	}
}
