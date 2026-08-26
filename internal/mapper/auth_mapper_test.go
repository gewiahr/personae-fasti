package mapper

import (
	"testing"

	"personae-fasti/internal/domain"
)

func TestFormLoginInfoResponseHandlesMissingRelations(t *testing.T) {
	response := FormLoginInfoResponse("Web token", &domain.Player{
		ExtID:    "player-ext",
		Username: "player-name",
	})

	if response.Authorization != "Web token" {
		t.Fatalf("authorization = %q", response.Authorization)
	}
	if response.Player.ExtID != "player-ext" || response.Player.Username != "player-name" {
		t.Fatalf("unexpected player response: %+v", response.Player)
	}
	if response.Player.Settings != nil {
		t.Fatalf("settings = %+v, want nil", response.Player.Settings)
	}
	if response.Player.CurrentGameExtID != "" {
		t.Fatalf("current game ext = %q, want empty", response.Player.CurrentGameExtID)
	}
}
