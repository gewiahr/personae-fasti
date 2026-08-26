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

func (s *EntitiesService) GetPlayerCurrentGameCharByExt(ctx context.Context, player *domain.Player, charExt string) (*domain.Char, error) {
	char, err := s.repo.GetCurrentGameCharByExt(ctx, player.CurrentGameID, charExt)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if char == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no character with ext %s", charExt))
	} else if err := ensureHiddenContentEditable(char.HiddenBy, player.ID); err != nil {
		return nil, err
	}
	if len(char.Records) > 0 {
		char.Records, err = s.recordRepo.FilterAllowedRecords(ctx, char.Records, player.ID)
	}
	return char, err
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

func (s *EntitiesService) GetPlayerCurrentGameNPCByExt(ctx context.Context, player *domain.Player, npcExt string) (*domain.NPC, error) {
	npc, err := s.repo.GetCurrentGameNPCByExt(ctx, player.CurrentGameID, npcExt)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if npc == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no npc with ext %s", npcExt))
	} else if err := ensureHiddenContentEditable(npc.HiddenBy, player.ID); err != nil {
		return nil, err
	}
	if len(npc.Records) > 0 {
		npc.Records, err = s.recordRepo.FilterAllowedRecords(ctx, npc.Records, player.ID)
	}
	return npc, err
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

func (s *EntitiesService) GetPlayerCurrentGameLocationByExt(ctx context.Context, player *domain.Player, locationExt string) (*domain.Location, error) {
	location, err := s.repo.GetCurrentGameLocationByExt(ctx, player.CurrentGameID, locationExt)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных локации", err)
	} else if location == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no location with ext %s", locationExt))
	} else if err := ensureHiddenContentEditable(location.HiddenBy, player.ID); err != nil {
		return nil, err
	}
	if len(location.Records) > 0 {
		location.Records, err = s.recordRepo.FilterAllowedRecords(ctx, location.Records, player.ID)
	}
	return location, err
}

func (s *EntitiesService) GetPlayerCurrentGameLocationChildrenByID(ctx context.Context, player *domain.Player, locationID int) ([]domain.Location, error) {
	locations, err := s.repo.GetCurrentGameLocationChildrenByID(ctx, player.CurrentGameID, locationID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения записей дочерних локаций", err)
	}

	return locations, nil
}

func (s *EntitiesService) PostPlayerCurrentGameChar(ctx context.Context, player *domain.Player, charCreate *dto.CharCreate) (*domain.Char, error) {
	if err := normalizeAndValidateEntity(&charCreate.Name, &charCreate.Title, charCreate.Description); err != nil {
		return nil, err
	}
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
	if err := normalizeAndValidateEntity(&npcCreate.Name, &npcCreate.Title, npcCreate.Description); err != nil {
		return nil, err
	}
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
	if err := normalizeAndValidateEntity(&locationCreate.Name, &locationCreate.Title, locationCreate.Description); err != nil {
		return nil, err
	}
	parentID, err := s.resolveLocationParent(ctx, player, locationCreate.ParentExtID, "")
	if err != nil {
		return nil, err
	}
	location := &domain.Location{
		Name:        locationCreate.Name,
		Title:       locationCreate.Title,
		Description: locationCreate.Description,
		GameID:      player.CurrentGameID,
		ParentID:    parentID,
		HiddenBy:    gewiutils.TernaryInt(locationCreate.Hidden, player.ID, 0),
	}

	location, err = s.repo.CreateLocation(ctx, location)
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
	existing, err := s.repo.GetCurrentGameCharByExt(ctx, player.CurrentGameID, charUpdate.ExtID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if existing == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no char with ext %s", charUpdate.ExtID))
	} else if err := ensureHiddenContentEditable(existing.HiddenBy, player.ID); err != nil {
		return nil, err
	}
	if err := normalizeAndValidateEntity(&charUpdate.Name, &charUpdate.Title, charUpdate.Description); err != nil {
		return nil, err
	}

	char, err := s.repo.EditChar(ctx, charUpdate, player.ID, player.CurrentGameID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if char == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no char with ext %s", charUpdate.ExtID))
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
	existing, err := s.repo.GetCurrentGameNPCByExt(ctx, player.CurrentGameID, NPCUpdate.ExtID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if existing == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no npc with ext %s", NPCUpdate.ExtID))
	} else if err := ensureHiddenContentEditable(existing.HiddenBy, player.ID); err != nil {
		return nil, err
	}
	if err := normalizeAndValidateEntity(&NPCUpdate.Name, &NPCUpdate.Title, NPCUpdate.Description); err != nil {
		return nil, err
	}

	npc, err := s.repo.EditNPC(ctx, NPCUpdate, player.ID, player.CurrentGameID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных персонажа", err)
	} else if npc == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no npc with ext %s", NPCUpdate.ExtID))
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
	existing, err := s.repo.GetCurrentGameLocationByExt(ctx, player.CurrentGameID, locationUpdate.ExtID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных локации", err)
	} else if existing == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no location with ext %s", locationUpdate.ExtID))
	} else if err := ensureHiddenContentEditable(existing.HiddenBy, player.ID); err != nil {
		return nil, err
	}
	if err := normalizeAndValidateEntity(&locationUpdate.Name, &locationUpdate.Title, locationUpdate.Description); err != nil {
		return nil, err
	}
	parentID, err := s.resolveLocationParent(ctx, player, locationUpdate.ParentExtID, locationUpdate.ExtID)
	if err != nil {
		return nil, err
	}
	locationUpdate.ParentID = parentID

	location, err := s.repo.EditLocation(ctx, locationUpdate, player.ID, player.CurrentGameID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения данных локации", err)
	} else if location == nil {
		return nil, e.NewNotFoundError(fmt.Sprintf("no location with ext %s", locationUpdate.ExtID))
	}

	if len(location.Records) > 0 {
		location.Records, err = s.recordRepo.FilterAllowedRecords(ctx, location.Records, player.ID)
		if err != nil {
			return nil, e.NewInternalError("Ошибка получения записей локации", err)
		}
	}

	return location, nil
}

func (s *EntitiesService) resolveLocationParent(ctx context.Context, player *domain.Player, parentExt, locationExt string) (int, error) {
	if parentExt == "" {
		return 0, nil
	}
	if parentExt == locationExt {
		return 0, e.NewFieldValidationError(map[string]string{"parentExt": "Место не может находиться внутри себя"})
	}
	parent, err := s.repo.GetCurrentGameLocationByExt(ctx, player.CurrentGameID, parentExt)
	if err != nil {
		return 0, e.NewInternalError("Ошибка получения родительской локации", err)
	}
	if parent == nil {
		return 0, e.NewFieldValidationError(map[string]string{"parentExt": "Родительское место не найдено в текущей игре"})
	}
	if err := ensureHiddenContentEditable(parent.HiddenBy, player.ID); err != nil {
		return 0, err
	}
	return parent.ID, nil
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
