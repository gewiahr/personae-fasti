package service

import (
	"context"
	"fmt"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/repo"
	gewiutils "personae-fasti/pkg/gewi-utils"
)

type EntitiesService struct {
	repo       repo.EntitiesRepository
	recordRepo repo.RecordRepository
}

func NewEntitiesService(repo repo.EntitiesRepository, recordRepo repo.RecordRepository) *EntitiesService {
	return &EntitiesService{
		repo:       repo,
		recordRepo: recordRepo,
	}
}

func (s *EntitiesService) GetPlayerCurrentGameChars(ctx context.Context, player *domain.Player) ([]domain.Char, error) {
	gameChars, err := s.repo.GetCurrentGameCharList(ctx, player.CurrentGameID, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных текущей игры", err)
	}
	return gameChars, nil
}

func (s *EntitiesService) GetPlayerCurrentGameNPCs(ctx context.Context, player *domain.Player) ([]domain.NPC, error) {
	npcList, err := s.repo.GetCurrentGameNPCList(ctx, player.CurrentGameID, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных текущей игры", err)
	}
	return npcList, nil
}

func (s *EntitiesService) GetPlayerCurrentGameLocations(ctx context.Context, player *domain.Player) ([]domain.Location, error) {
	locationList, err := s.repo.GetCurrentGameLocationList(ctx, player.CurrentGameID, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных текущей игры", err)
	}
	return locationList, nil
}

func (s *EntitiesService) GetPlayerCurrentGameCharByID(ctx context.Context, player *domain.Player, charID int) (*domain.Char, error) {
	char, err := s.repo.GetCurrentGameCharByID(ctx, player.CurrentGameID, charID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if char == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no character with id %d", charID))
	} else if char.HiddenBy != 0 && char.HiddenBy != player.ID {
		return nil, e.NewForbiddenError(fmt.Sprintf("char %d is not allowed to request for the player %d", char.ID, player.ID))
	}

	if len(char.Records) > 0 {
		char.Records, err = s.recordRepo.FilterAllowedRecords(ctx, char.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей персонажа", err)
		}
	}

	return char, nil
}

func (s *EntitiesService) GetPlayerCurrentGameNPCByID(ctx context.Context, player *domain.Player, npcID int) (*domain.NPC, error) {
	npc, err := s.repo.GetCurrentGameNPCByID(ctx, player.CurrentGameID, npcID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if npc == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no npc with id %d", npcID))
	} else if npc.HiddenBy != 0 && npc.HiddenBy != player.ID {
		return nil, e.NewForbiddenError(fmt.Sprintf("npc %d is not allowed to request for the player %d", npc.ID, player.ID))
	}

	if len(npc.Records) > 0 {
		npc.Records, err = s.recordRepo.FilterAllowedRecords(ctx, npc.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей персонажа", err)
		}
	}

	return npc, nil
}

func (s *EntitiesService) GetPlayerCurrentGameLocationByID(ctx context.Context, player *domain.Player, locationID int) (*domain.Location, error) {
	location, err := s.repo.GetCurrentGameLocationByID(ctx, player.CurrentGameID, locationID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных локации", err)
	} else if location == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no location with id %d", locationID))
	} else if location.HiddenBy != 0 && location.HiddenBy != player.ID {
		return nil, e.NewForbiddenError(fmt.Sprintf("location %d is not allowed to request for the player %d", location.ID, player.ID))
	}

	if len(location.Records) > 0 {
		location.Records, err = s.recordRepo.FilterAllowedRecords(ctx, location.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей локации", err)
		}
	}

	return location, nil
}

func (s *EntitiesService) GetPlayerCurrentGameLocationChildrenByID(ctx context.Context, player *domain.Player, locationID int) ([]domain.Location, error) {
	locations, err := s.repo.GetCurrentGameLocationChildrenByID(ctx, player.CurrentGameID, locationID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения записей дочерних локаций", err)
	}

	return locations, nil
}

func (s *EntitiesService) PostPlayerCurrentGameChar(ctx context.Context, player *domain.Player, charCreate *dto.CharCreate) (*domain.Char, error) {
	char := &domain.Char{
		Name:        charCreate.Name,
		Title:       charCreate.Title,
		Description: charCreate.Description,
		PlayerID:    player.ID,
		GameID:      player.CurrentGameID,
		HiddenBy:    gewiutils.TernaryInt(charCreate.Hidden, player.ID, 0),
	}

	char, err := s.repo.CreateChar(ctx, char)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if char == nil {
		return nil, e.NewNotFoundError("char not created")
	} else if char.HiddenBy != 0 && char.HiddenBy != player.ID {
		return nil, e.NewForbiddenError(fmt.Sprintf("char %d is not allowed to request for the player %d", char.ID, player.ID))
	}

	if len(char.Records) > 0 {
		char.Records, err = s.recordRepo.FilterAllowedRecords(ctx, char.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей персонажа", err)
		}
	}

	return char, nil
}

func (s *EntitiesService) PostPlayerCurrentGameNPC(ctx context.Context, player *domain.Player, npcCreate *dto.NPCCreate) (*domain.NPC, error) {
	npc := &domain.NPC{
		Name:        npcCreate.Name,
		Title:       npcCreate.Title,
		Description: npcCreate.Description,
		GameID:      player.CurrentGameID,
		HiddenBy:    gewiutils.TernaryInt(npcCreate.Hidden, player.ID, 0),
	}

	npc, err := s.repo.CreateNPC(ctx, npc)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if npc == nil {
		return nil, e.NewNotFoundError("npc not created")
	} else if npc.HiddenBy != 0 && npc.HiddenBy != player.ID {
		return nil, e.NewForbiddenError(fmt.Sprintf("npc %d is not allowed to request for the player %d", npc.ID, player.ID))
	}

	if len(npc.Records) > 0 {
		npc.Records, err = s.recordRepo.FilterAllowedRecords(ctx, npc.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей персонажа", err)
		}
	}

	return npc, nil
}

func (s *EntitiesService) PostPlayerCurrentGameLocation(ctx context.Context, player *domain.Player, locationCreate *dto.LocationCreate) (*domain.Location, error) {
	location := &domain.Location{
		Name:        locationCreate.Name,
		Title:       locationCreate.Title,
		Description: locationCreate.Description,
		GameID:      player.CurrentGameID,
		ParentID:    locationCreate.ParentID,
		HiddenBy:    gewiutils.TernaryInt(locationCreate.Hidden, player.ID, 0),
	}

	location, err := s.repo.CreateLocation(ctx, location)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных локации", err)
	} else if location == nil {
		return nil, e.NewNotFoundError("location not created")
	} else if location.HiddenBy != 0 && location.HiddenBy != player.ID {
		return nil, e.NewForbiddenError(fmt.Sprintf("location %d is not allowed to request for the player %d", location.ID, player.ID))
	}

	if len(location.Records) > 0 {
		location.Records, err = s.recordRepo.FilterAllowedRecords(ctx, location.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей локации", err)
		}
	}

	return location, nil
}

func (s *EntitiesService) EditPlayerCurrentGameChar(ctx context.Context, player *domain.Player, charUpdate *dto.CharUpdate) (*domain.Char, error) {
	existing, err := s.repo.GetCurrentGameCharByID(ctx, player.CurrentGameID, charUpdate.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if existing == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no char with id %d", charUpdate.ID))
	} else if err := ensureHiddenContentEditable(existing.HiddenBy, player.ID); err != nil {
		return nil, err
	}

	char, err := s.repo.EditChar(ctx, charUpdate, player.ID, player.CurrentGameID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if char == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no char with id %d", charUpdate.ID))
	}

	if len(char.Records) > 0 {
		char.Records, err = s.recordRepo.FilterAllowedRecords(ctx, char.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей персонажа", err)
		}
	}

	return char, nil
}

func (s *EntitiesService) EditPlayerCurrentGameNPC(ctx context.Context, player *domain.Player, NPCUpdate *dto.NPCUpdate) (*domain.NPC, error) {
	existing, err := s.repo.GetCurrentGameNPCByID(ctx, player.CurrentGameID, NPCUpdate.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if existing == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no npc with id %d", NPCUpdate.ID))
	} else if err := ensureHiddenContentEditable(existing.HiddenBy, player.ID); err != nil {
		return nil, err
	}

	npc, err := s.repo.EditNPC(ctx, NPCUpdate, player.ID, player.CurrentGameID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if npc == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no npc with id %d", NPCUpdate.ID))
	}

	if len(npc.Records) > 0 {
		npc.Records, err = s.recordRepo.FilterAllowedRecords(ctx, npc.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей персонажа", err)
		}
	}

	return npc, nil
}

func (s *EntitiesService) EditPlayerCurrentGameLocation(ctx context.Context, player *domain.Player, locationUpdate *dto.LocationUpdate) (*domain.Location, error) {
	existing, err := s.repo.GetCurrentGameLocationByID(ctx, player.CurrentGameID, locationUpdate.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных локации", err)
	} else if existing == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no location with id %d", locationUpdate.ID))
	} else if err := ensureHiddenContentEditable(existing.HiddenBy, player.ID); err != nil {
		return nil, err
	}

	location, err := s.repo.EditLocation(ctx, locationUpdate, player.ID, player.CurrentGameID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных локации", err)
	} else if location == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no location with id %d", locationUpdate.ID))
	}

	if len(location.Records) > 0 {
		location.Records, err = s.recordRepo.FilterAllowedRecords(ctx, location.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей локации", err)
		}
	}

	return location, nil
}

func ensureHiddenContentEditable(hiddenBy, playerID int) error {
	if hiddenBy != 0 && hiddenBy != playerID {
		return e.NewForbiddenError("content is hidden by another player")
	}
	return nil
}

func (s *EntitiesService) GetPlayerCurrentGameSuggestions(ctx context.Context, player *domain.Player) ([]dto.Suggestion, error) {
	suggestions, err := s.repo.GetCurrentGameSuggestionList(ctx, player.CurrentGameID, player.ID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных текущей игры", err)
	}
	return suggestions, nil
}
