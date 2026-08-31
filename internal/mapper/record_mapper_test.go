package mapper

import (
	"testing"

	"personae-fasti/internal/domain"
)

func TestRecordMapperResolvesOwnerAndQuestByInternalID(t *testing.T) {
	players := []domain.Player{
		{ID: 1, ExtID: "requesting-player"},
		{ID: 2, ExtID: "record-owner"},
	}
	quests := []domain.Quest{
		{ID: 10, ExtID: "other-quest"},
		{ID: 20, ExtID: "record-quest"},
	}
	record := domain.Record{
		ID:       30,
		PlayerID: 2,
		QuestID:  20,
		Player:   &players[0],
		Quest:    &quests[0],
	}

	got := RecordToRecordFull(record, players, quests, "game-ext")

	if got.PlayerExtID != "record-owner" {
		t.Fatalf("player ext = %q, want record owner ext", got.PlayerExtID)
	}
	if got.QuestExtID != "record-quest" {
		t.Fatalf("quest ext = %q, want record quest ext", got.QuestExtID)
	}
}

func TestRecordMapperLeavesUnknownRelationsEmpty(t *testing.T) {
	got := RecordToRecordFull(domain.Record{PlayerID: 99, QuestID: 99}, nil, nil, "game-ext")

	if got.PlayerExtID != "" {
		t.Fatalf("player ext = %q, want empty", got.PlayerExtID)
	}
	if got.QuestExtID != "" {
		t.Fatalf("quest ext = %q, want empty", got.QuestExtID)
	}
}
