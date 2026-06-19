package mapper

import (
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
)

func PlayerToPlayersBriefArray(players []domain.Player) []dto.PlayerBrief {
	playerInfoArray := []dto.PlayerBrief{}
	for _, player := range players {
		playerInfoArray = append(playerInfoArray, PlayerToPlayersBrief(player))
	}

	return playerInfoArray
}

func PlayerToPlayersBrief(player domain.Player) dto.PlayerBrief {
	return dto.PlayerBrief{
		ExtID:    player.ExtID,
		Username: player.Username,
	}

}
