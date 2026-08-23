package mapper

import (
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
)

func RecordToRecordFullArray(records []domain.Record, playerExt, gameExt string) []dto.RecordFull {
	recordsFull := make([]dto.RecordFull, 0, len(records))
	for _, r := range records {
		recordsFull = append(recordsFull, RecordToRecordFull(r, playerExt, gameExt))
	}
	return recordsFull
}

func RecordToRecordFull(r domain.Record, playerExt, gameExt string) dto.RecordFull {
	return dto.RecordFull{
		ID:          r.ID,
		Text:        r.Text,
		PlayerExtID: playerExt,
		GameExtID:   gameExt,
		QuestExtID:  recordQuestExt(r),
		Hidden:      r.HiddenBy != 0,
		Created:     r.Created,
		Updated:     r.Updated,
		Deleted:     r.Deleted,
	}
}

func recordQuestExt(record domain.Record) string {
	if record.Quest == nil {
		return ""
	}
	return record.Quest.ExtID
}
