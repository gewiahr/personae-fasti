package mapper

import (
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
)

func GameToGameBrief(game *domain.Game) *dto.GameBrief {
	gmExtID := ""
	if game.GM != nil {
		gmExtID = game.GM.ExtID
	}
	return &dto.GameBrief{
		ExtID:   game.ExtID,
		Title:   game.Name,
		GMExtID: gmExtID,
	}
}

func GameToGameBriefArray(games []domain.Game) []dto.GameBrief {
	gameInfoArray := []dto.GameBrief{}
	for _, game := range games {
		gameInfoArray = append(gameInfoArray, *GameToGameBrief(&game))
	}

	return gameInfoArray
}

func GameToGameFull(game *domain.Game) *dto.GameFull {
	gmExtID := ""
	if game.GM != nil {
		gmExtID = game.GM.ExtID
	}
	allowAllEditRecords := false
	if game.Settings != nil {
		allowAllEditRecords = game.Settings.AllowAllEditRecords
	}
	return &dto.GameFull{
		ExtID:   game.ExtID,
		Title:   game.Name,
		GMExtID: gmExtID,

		Settings: &dto.GameSettings{
			AllowAllEditRecords: allowAllEditRecords,
		},
		Sessions: SessionToSessionBriefArray(game.Sessions),
		Players:  PlayerToPlayersBriefArray(game.Players), // TODO: add observers and left
		Invites:  PlayerToPlayersBriefArray(game.Invites), // TODO: add observers and left
	}
}

func SessionToSessionBrief(session *domain.Session) *dto.SessionBrief {
	return &dto.SessionBrief{
		Number:  session.Number,
		Name:    session.Name,
		EndTime: session.EndTime,
	}
}

func SessionToSessionBriefArray(sessions []domain.Session) []dto.SessionBrief {
	sessionInfoArray := []dto.SessionBrief{}
	for _, session := range sessions {
		sessionInfoArray = append(sessionInfoArray, *SessionToSessionBrief(&session))
	}

	return sessionInfoArray
}

func InviteToGameInvite(invite *domain.GameInvite) *dto.GameInvite {
	return &dto.GameInvite{
		PlayerExt:  invite.Player.ExtID,
		GameExt:    invite.Game.ExtID,
		GameTitle:  invite.Game.Name,
		InviteCode: invite.Code,
	}
}

func InviteToGameInviteArray(invites []domain.GameInvite) []dto.GameInvite {
	gameInviteArray := []dto.GameInvite{}
	for _, invite := range invites {
		gameInviteArray = append(gameInviteArray, *InviteToGameInvite(&invite))
	}

	return gameInviteArray
}
