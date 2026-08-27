package mapper

import (
	"testing"

	"personae-fasti/internal/domain"
)

func TestGameToGameFullHandlesMissingRelations(t *testing.T) {
	response := GameToGameFull(&domain.Game{ExtID: "game-ext", Name: "Game"})
	if response.GMExtID != "" {
		t.Fatalf("GM ext = %q, want empty", response.GMExtID)
	}
	if response.Settings == nil || response.Settings.AllowAllEditRecords {
		t.Fatalf("unexpected settings: %+v", response.Settings)
	}
	if response.Players == nil || response.Sessions == nil || response.Invites == nil {
		t.Fatalf("collections must be empty arrays: %+v", response)
	}
}
