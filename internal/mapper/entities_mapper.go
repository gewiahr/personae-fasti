package mapper

import (
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
)

func CharToCharBriefArray(chars []domain.Char, gameExt, playerExt string) []dto.CharBrief {
	charBriefArray := []dto.CharBrief{}
	for _, char := range chars {
		charBriefArray = append(charBriefArray, CharToCharBrief(char, gameExt, playerExt))
	}

	return charBriefArray
}

func CharToCharBrief(char domain.Char, gameExt, playerExt string) dto.CharBrief {
	return dto.CharBrief{
		ExtID:       char.ExtID,
		Name:        char.Name,
		Title:       char.Title,
		PlayerExtID: playerExt,
		GameExtID:   gameExt,
		Hidden:      char.HiddenBy != 0,
	}
}

func CharToCharFullInfo(char *domain.Char, gameExt, playerExt string) *dto.CharFull {
	return &dto.CharFull{
		ExtID:       char.ExtID,
		Name:        char.Name,
		Title:       char.Title,
		Description: char.Description,
		PlayerExtID: playerExt,
		GameExtID:   gameExt,
		Hidden:      char.HiddenBy != 0,
	}
}

func NPCToNPCBriefArray(npcs []domain.NPC, gameExt string) []dto.NPCBrief {
	npcBriefArray := []dto.NPCBrief{}
	for _, npc := range npcs {
		npcBriefArray = append(npcBriefArray, NPCToNPCBrief(npc, gameExt))
	}

	return npcBriefArray
}

func NPCToNPCBrief(npc domain.NPC, gameExt string) dto.NPCBrief {
	return dto.NPCBrief{
		ExtID:     npc.ExtID,
		Name:      npc.Name,
		Title:     npc.Title,
		GameExtID: gameExt,
		Hidden:    npc.HiddenBy != 0,
	}
}

func NPCToNPCFullInfo(npc *domain.NPC, gameExt string) *dto.NPCFull {
	return &dto.NPCFull{
		ExtID:       npc.ExtID,
		Name:        npc.Name,
		Title:       npc.Title,
		Description: npc.Description,
		GameExtID:   gameExt,
		Hidden:      npc.HiddenBy != 0,
	}
}

func LocationToLocationBriefArray(locations []domain.Location, gameExt string) []dto.LocationBrief {
	locationBriefArray := []dto.LocationBrief{}
	for _, location := range locations {
		locationBriefArray = append(locationBriefArray, *LocationToLocationBrief(location, gameExt))
	}

	return locationBriefArray
}

func LocationToLocationBrief(location domain.Location, gameExt string) *dto.LocationBrief {
	return &dto.LocationBrief{
		ExtID:     location.ExtID,
		Name:      location.Name,
		Title:     location.Title,
		GameExtID: gameExt,
		Hidden:    location.HiddenBy != 0,
	}
}

func LocationToLocationFullInfo(location *domain.Location, gameExt string) *dto.LocationFull {
	return &dto.LocationFull{
		ExtID:       location.ExtID,
		ParentExtID: locationParentExt(location),
		Name:        location.Name,
		Title:       location.Title,
		Description: location.Description,
		GameExtID:   gameExt,
		Hidden:      location.HiddenBy != 0,
	}
}

func locationParentExt(location *domain.Location) string {
	if location.Parent == nil {
		return ""
	}
	return location.Parent.ExtID
}
