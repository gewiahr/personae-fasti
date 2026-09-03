package mapper

import (
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
)

func RecordToRecordFullArray(records []domain.Record, players []domain.Player, quests []domain.Quest, gameExt string) []dto.RecordFull {
	recordsFull := make([]dto.RecordFull, 0, len(records))
	for _, r := range records {
		recordsFull = append(recordsFull, RecordToRecordFull(r, players, quests, gameExt))
	}
	return recordsFull
}

func RecordToRecordFull(r domain.Record, players []domain.Player, quests []domain.Quest, gameExt string) dto.RecordFull {
	return dto.RecordFull{
		ExtID:       r.ExtID,
		Text:        r.Text,
		PlayerExtID: recordPlayerExt(r, players),
		GameExtID:   gameExt,
		QuestExtID:  recordQuestExt(r, quests),
		Hidden:      r.HiddenBy != 0,
		Created:     r.Created,
		Updated:     r.Updated,
		Deleted:     r.Deleted,
	}
}

func recordPlayerExt(record domain.Record, players []domain.Player) string {
	for _, player := range players {
		if player.ID == record.PlayerID {
			return player.ExtID
		}
	}
	return ""
}

func recordQuestExt(record domain.Record, quests []domain.Quest) string {
	for _, quest := range quests {
		if quest.ID == record.QuestID {
			return quest.ExtID
		}
	}
	return ""
}
