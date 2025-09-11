package respData

import "personae-fasti/data"

func GameToGameInfo(game *data.Game) *GameInfo {
	return &GameInfo{
		ID:    game.ID,
		Title: game.Name,
		GMID:  game.GMID,
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
			HiddenBy: char.HiddenBy,
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
		HiddenBy:    char.HiddenBy,
	}
}

func NPCToNPCInfoArray(npcs []data.NPC) []NPCInfo {
	npcInfoArray := []NPCInfo{}
	for _, npc := range npcs {
		npcInfoArray = append(npcInfoArray, NPCInfo{
			ID:       npc.ID,
			Name:     npc.Name,
			Title:    npc.Title,
			GameID:   npc.GameID,
			HiddenBy: npc.HiddenBy,
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
		HiddenBy:    npc.HiddenBy,
	}
}

func LocationToLocationInfo(location *data.Location) *LocationInfo {
	return &LocationInfo{
		ID:       location.ID,
		Name:     location.Name,
		Title:    location.Title,
		GameID:   location.GameID,
		HiddenBy: location.HiddenBy,
	}
}

func LocationToLocationInfoArray(locations []data.Location) []LocationInfo {
	locationInfoArray := []LocationInfo{}
	for _, location := range locations {
		locationInfoArray = append(locationInfoArray, LocationInfo{
			ID:       location.ID,
			Name:     location.Name,
			Title:    location.Title,
			GameID:   location.GameID,
			HiddenBy: location.HiddenBy,
		})
	}

	return locationInfoArray
}

func LocationToLocationFullInfo(location *data.Location) *LocationFullInfo {
	return &LocationFullInfo{
		ID:          location.ID,
		Name:        location.Name,
		Title:       location.Title,
		Description: location.Description,
		ParentID:    location.ParentID,
		GameID:      location.GameID,
		HiddenBy:    location.HiddenBy,
	}
}

func QuestToQuestInfoArray(quests []data.Quest) []QuestInfo {
	questInfoArray := []QuestInfo{}
	for _, quest := range quests {
		finishedQuest := false
		if quest.Finished != nil {
			finishedQuest = true
		}
		questInfoArray = append(questInfoArray, QuestInfo{
			ID:         quest.ID,
			Name:       quest.Name,
			Title:      quest.Title,
			GameID:     quest.GameID,
			Successful: quest.Successful,
			HiddenBy:   quest.HiddenBy,
			Finished:   finishedQuest,
		})
	}

	return questInfoArray
}

func QuestToQuestFullInfo(quest *data.Quest) *QuestFullInfo {
	finishedQuest := false
	if quest.Finished != nil {
		finishedQuest = true
	}
	return &QuestFullInfo{
		ID:          quest.ID,
		Name:        quest.Name,
		Title:       quest.Title,
		Description: quest.Description,
		ParentID:    quest.ParentID,
		ChildID:     quest.ChildID,
		HeadID:      quest.HeadID,
		GameID:      quest.GameID,
		Successful:  quest.Successful,
		HiddenBy:    quest.HiddenBy,
		Finished:    finishedQuest,
	}
}

func TaskToTaskFullInfoArray(tasks []data.QuestTask) []QuestTaskFullInfo {
	taskInfoArray := []QuestTaskFullInfo{}
	for _, task := range tasks {
		finishedTask := false
		if task.Finished != nil {
			finishedTask = true
		}
		taskInfoArray = append(taskInfoArray, QuestTaskFullInfo{
			ID:          task.ID,
			QuestID:     task.QuestID,
			Name:        task.Name,
			Description: task.Description,
			GameID:      task.GameID,
			Type:        int(task.Type),
			Capacity:    task.Capacity,
			Current:     task.Current,
			HiddenBy:    task.HiddenBy,
			Finished:    finishedTask,
		})
	}

	return taskInfoArray
}
