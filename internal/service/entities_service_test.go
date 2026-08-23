package service

import (
	"context"
	"errors"
	"testing"

	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/repo"
)

type entitiesRepoStub struct {
	repo.EntitiesRepository
	char       *domain.Char
	npc        *domain.NPC
	location   *domain.Location
	editCalled bool
	editGameID int
}

func (r *entitiesRepoStub) GetCurrentGameCharByID(_ context.Context, gameID, charID int) (*domain.Char, error) {
	if r.char == nil || r.char.ID != charID || r.char.GameID != gameID {
		return nil, nil
	}
	copy := *r.char
	return &copy, nil
}

func (r *entitiesRepoStub) GetCurrentGameCharByExt(_ context.Context, gameID int, ext string) (*domain.Char, error) {
	if r.char == nil || r.char.ExtID != ext || r.char.GameID != gameID {
		return nil, nil
	}
	copy := *r.char
	return &copy, nil
}

func (r *entitiesRepoStub) GetCurrentGameNPCByID(_ context.Context, gameID, npcID int) (*domain.NPC, error) {
	if r.npc == nil || r.npc.ID != npcID || r.npc.GameID != gameID {
		return nil, nil
	}
	copy := *r.npc
	return &copy, nil
}

func (r *entitiesRepoStub) GetCurrentGameNPCByExt(_ context.Context, gameID int, ext string) (*domain.NPC, error) {
	if r.npc == nil || r.npc.ExtID != ext || r.npc.GameID != gameID {
		return nil, nil
	}
	copy := *r.npc
	return &copy, nil
}

func (r *entitiesRepoStub) GetCurrentGameLocationByID(_ context.Context, gameID, locationID int) (*domain.Location, error) {
	if r.location == nil || r.location.ID != locationID || r.location.GameID != gameID {
		return nil, nil
	}
	copy := *r.location
	return &copy, nil
}

func (r *entitiesRepoStub) GetCurrentGameLocationByExt(_ context.Context, gameID int, ext string) (*domain.Location, error) {
	if r.location == nil || r.location.ExtID != ext || r.location.GameID != gameID {
		return nil, nil
	}
	copy := *r.location
	return &copy, nil
}

func (r *entitiesRepoStub) EditChar(_ context.Context, update *dto.CharUpdate, playerID, gameID int) (*domain.Char, error) {
	r.editCalled = true
	r.editGameID = gameID
	return &domain.Char{ExtID: update.ExtID, GameID: gameID, HiddenBy: hiddenBy(update.Hidden, playerID)}, nil
}

func (r *entitiesRepoStub) EditNPC(_ context.Context, update *dto.NPCUpdate, playerID, gameID int) (*domain.NPC, error) {
	r.editCalled = true
	r.editGameID = gameID
	return &domain.NPC{ExtID: update.ExtID, GameID: gameID, HiddenBy: hiddenBy(update.Hidden, playerID)}, nil
}

func (r *entitiesRepoStub) EditLocation(_ context.Context, update *dto.LocationUpdate, playerID, gameID int) (*domain.Location, error) {
	r.editCalled = true
	r.editGameID = gameID
	return &domain.Location{ExtID: update.ExtID, GameID: gameID, HiddenBy: hiddenBy(update.Hidden, playerID)}, nil
}

func hiddenBy(hidden bool, playerID int) int {
	if hidden {
		return playerID
	}
	return 0
}

func TestEditEntityAllowsVisibleAndOwnHiddenContent(t *testing.T) {
	for _, hiddenByID := range []int{0, 12} {
		t.Run(map[bool]string{true: "own hidden", false: "visible"}[hiddenByID != 0], func(t *testing.T) {
			repository := &entitiesRepoStub{char: &domain.Char{ID: 4, ExtID: "char-ext", GameID: 8, HiddenBy: hiddenByID}}
			service := NewEntitiesService(repository, nil)
			player := &domain.Player{ID: 12, CurrentGameID: 8}

			_, err := service.EditPlayerCurrentGameChar(context.Background(), player, &dto.CharUpdate{ExtID: "char-ext"})
			if err != nil {
				t.Fatalf("edit returned an unexpected error: %v", err)
			}
			if !repository.editCalled {
				t.Fatal("expected repository edit to be called")
			}
			if repository.editGameID != player.CurrentGameID {
				t.Fatalf("edit used game %d, want %d", repository.editGameID, player.CurrentGameID)
			}
		})
	}
}

func TestEditEntityRejectsContentHiddenByAnotherPlayerBeforeWrite(t *testing.T) {
	tests := []struct {
		name string
		repo *entitiesRepoStub
		edit func(*EntitiesService, *domain.Player) error
	}{
		{
			name: "char",
			repo: &entitiesRepoStub{char: &domain.Char{ID: 1, ExtID: "char-ext", GameID: 8, HiddenBy: 99}},
			edit: func(s *EntitiesService, p *domain.Player) error {
				_, err := s.EditPlayerCurrentGameChar(context.Background(), p, &dto.CharUpdate{ExtID: "char-ext"})
				return err
			},
		},
		{
			name: "npc",
			repo: &entitiesRepoStub{npc: &domain.NPC{ID: 2, ExtID: "npc-ext", GameID: 8, HiddenBy: 99}},
			edit: func(s *EntitiesService, p *domain.Player) error {
				_, err := s.EditPlayerCurrentGameNPC(context.Background(), p, &dto.NPCUpdate{ExtID: "npc-ext"})
				return err
			},
		},
		{
			name: "location",
			repo: &entitiesRepoStub{location: &domain.Location{ID: 3, ExtID: "location-ext", GameID: 8, HiddenBy: 99}},
			edit: func(s *EntitiesService, p *domain.Player) error {
				_, err := s.EditPlayerCurrentGameLocation(context.Background(), p, &dto.LocationUpdate{ExtID: "location-ext"})
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			player := &domain.Player{ID: 12, CurrentGameID: 8}
			err := test.edit(NewEntitiesService(test.repo, nil), player)
			assertAppErrorType(t, err, e.ErrForbidden)
			if test.repo.editCalled {
				t.Fatal("repository edit was called before authorization")
			}
		})
	}
}

func TestEditEntityRejectsCrossGameIDBeforeWrite(t *testing.T) {
	repository := &entitiesRepoStub{char: &domain.Char{ID: 4, ExtID: "char-ext", GameID: 9}}
	service := NewEntitiesService(repository, nil)
	player := &domain.Player{ID: 12, CurrentGameID: 8}

	_, err := service.EditPlayerCurrentGameChar(context.Background(), player, &dto.CharUpdate{ExtID: "char-ext"})
	assertAppErrorType(t, err, e.ErrNotFound)
	if repository.editCalled {
		t.Fatal("repository edit was called for a cross-game entity")
	}
}

func assertAppErrorType(t *testing.T, err error, want error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected %v error, got nil", want)
	}
	var appErr *e.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Type != want {
		t.Fatalf("got error type %v, want %v", appErr.Type, want)
	}
}
