package mapper

import (
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
)

func FormLoginInfoResponse(auth string, player *domain.Player) *dto.LoginInfo {
	loginInfo := dto.LoginInfo{
		Authorization: auth,
		Player: dto.LoginPlayerInfo{
			ExtID:    player.ExtID,
			Username: player.Username,
			Settings: nil,
		},
	}

	if player.RegData.UsernameSet {
		loginInfo.Player.Settings = &dto.LoginPlayerInfoSettings{
			CouldChangeUsername: false,
		}
	}

	if player.CurrentGame != nil {
		loginInfo.Player.CurrentGameExtID = player.CurrentGame.ExtID
	}

	return &loginInfo
}
