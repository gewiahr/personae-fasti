package mapper

import (
	"testing"

	"personae-fasti/internal/domain"
)

func TestEntityMappersExposePublicExt(t *testing.T) {
	const ext = "QRTWXY346789"

	tests := []struct {
		name string
		got  string
	}{
		{name: "char brief", got: CharToCharBrief(domain.Char{ExtID: ext}, "game", "player").ExtID},
		{name: "char full", got: CharToCharFullInfo(&domain.Char{ExtID: ext}, "game", "player").ExtID},
		{name: "npc brief", got: NPCToNPCBrief(domain.NPC{ExtID: ext}, "game").ExtID},
		{name: "npc full", got: NPCToNPCFullInfo(&domain.NPC{ExtID: ext}, "game").ExtID},
		{name: "location brief", got: LocationToLocationBrief(domain.Location{ExtID: ext}, "game").ExtID},
		{name: "location full", got: LocationToLocationFullInfo(&domain.Location{ExtID: ext}, "game").ExtID},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != ext {
				t.Fatalf("got ext %q, want %q", test.got, ext)
			}
		})
	}
}

func TestQuestMappersExposePublicExt(t *testing.T) {
	const ext = "QRTWXY346789"
	quest := domain.Quest{ExtID: ext}

	if got := QuestToQuestBrief(quest, "game").ExtID; got != ext {
		t.Fatalf("brief ext = %q, want %q", got, ext)
	}
	if got := QuestToQuestFullInfo(&quest, "game").ExtID; got != ext {
		t.Fatalf("full ext = %q, want %q", got, ext)
	}
}
