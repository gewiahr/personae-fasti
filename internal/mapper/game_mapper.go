package mapper

import (
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
)

func GameToGameBrief(game *domain.Game) *dto.GameBrief {
	return &dto.GameBrief{
		ExtID:   game.ExtID,
		Title:   game.Name,
		GMExtID: game.GM.ExtID,
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
	return &dto.GameFull{
		ExtID:   game.ExtID,
		Title:   game.Name,
		GMExtID: game.GM.ExtID,

		Settings: &dto.GameSettings{
			AllowAllEditRecords: game.Settings.AllowAllEditRecords,
		},
		Sessions: SessionToSessionBriefArray(game.Sessions),
		Players:  PlayerToPlayersBriefArray(game.Players), // TODO: add observers and left
	}
}

func SessionToSessionBrief(session domain.Session) dto.SessionBrief {
	return dto.SessionBrief{
		Number:  session.Number,
		Name:    session.Name,
		EndTime: session.EndTime,
	}
}

func SessionToSessionBriefArray(sessions []domain.Session) []dto.SessionBrief {
	sessionInfoArray := []dto.SessionBrief{}
	for _, session := range sessions {
		sessionInfoArray = append(sessionInfoArray, SessionToSessionBrief(session))
	}

	return sessionInfoArray
}
