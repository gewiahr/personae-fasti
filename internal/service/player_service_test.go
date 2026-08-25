package service

import (
	"context"
	"errors"
	"testing"

	"personae-fasti/internal/domain"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/repo"
)

type playerNoteRepoStub struct {
	repo.PlayerRepository
	note           string
	err            error
	readPlayerID   int
	updatePlayerID int
}

func (r *playerNoteRepoStub) GetPersonalNote(_ context.Context, playerID int) (string, error) {
	r.readPlayerID = playerID
	return r.note, r.err
}

func (r *playerNoteRepoStub) UpdatePersonalNote(_ context.Context, playerID int, note string) (string, error) {
	r.updatePlayerID = playerID
	if r.err != nil {
		return "", r.err
	}
	r.note = note
	return note, nil
}

func TestPersonalNoteUsesAuthenticatedPlayer(t *testing.T) {
	repository := &playerNoteRepoStub{note: "old note"}
	service := NewPlayerService(repository, nil)
	player := &domain.Player{ID: 42}

	note, err := service.GetPersonalNote(context.Background(), player)
	if err != nil {
		t.Fatalf("get personal note returned an unexpected error: %v", err)
	}
	if note != "old note" {
		t.Fatalf("get personal note returned %q, want %q", note, "old note")
	}
	if repository.readPlayerID != player.ID {
		t.Fatalf("get used player %d, want %d", repository.readPlayerID, player.ID)
	}

	note, err = service.UpdatePersonalNote(context.Background(), player, "new note")
	if err != nil {
		t.Fatalf("update personal note returned an unexpected error: %v", err)
	}
	if note != "new note" || repository.note != "new note" {
		t.Fatalf("update returned %q and stored %q, want %q", note, repository.note, "new note")
	}
	if repository.updatePlayerID != player.ID {
		t.Fatalf("update used player %d, want %d", repository.updatePlayerID, player.ID)
	}
}

func TestPersonalNoteWrapsRepositoryErrors(t *testing.T) {
	repository := &playerNoteRepoStub{err: errors.New("database unavailable")}
	service := NewPlayerService(repository, nil)
	player := &domain.Player{ID: 42}

	if _, err := service.GetPersonalNote(context.Background(), player); err == nil {
		t.Fatal("get personal note returned nil error")
	} else {
		assertAppErrorType(t, err, e.ErrInternal)
	}

	if _, err := service.UpdatePersonalNote(context.Background(), player, "note"); err == nil {
		t.Fatal("update personal note returned nil error")
	} else {
		assertAppErrorType(t, err, e.ErrInternal)
	}
}
