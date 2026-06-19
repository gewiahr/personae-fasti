package bunrepo

import (
	"context"
	"database/sql"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	"personae-fasti/internal/repo"
	gewiutils "personae-fasti/pkg/gewi-utils"

	"github.com/uptrace/bun"
)

var _ repo.EntitiesRepository = (*EntitiesRepo)(nil)

type EntitiesRepo struct {
	db *bun.DB
}

func NewEntitiesRepo(db *bun.DB) *EntitiesRepo {
	return &EntitiesRepo{db: db}
}

func (r *EntitiesRepo) GetCurrentGameCharList(ctx context.Context, gameID, playerID int) ([]domain.Char, error) {
	charList := []domain.Char{}
	err := r.db.NewSelect().
		Model(&charList).
		Where("game_id = ? AND (hidden_by = 0 OR hidden_by = ?)", gameID, playerID).
		Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return charList, nil
}

func (r *EntitiesRepo) GetCurrentGameNPCList(ctx context.Context, gameID, playerID int) ([]domain.NPC, error) {
	npcList := []domain.NPC{}
	err := r.db.NewSelect().
		Model(&npcList).
		Where("game_id = ? AND (hidden_by = 0 OR hidden_by = ?)", gameID, playerID).
		Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return npcList, nil
}

func (r *EntitiesRepo) GetCurrentGameLocationList(ctx context.Context, gameID, playerID int) ([]domain.Location, error) {
	locationList := []domain.Location{}
	err := r.db.NewSelect().
		Model(&locationList).
		Where("game_id = ? AND (hidden_by = 0 OR hidden_by = ?)", gameID, playerID).
		Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return locationList, nil
}

func (r *EntitiesRepo) GetCurrentGameCharByID(ctx context.Context, gameID, charID int) (*domain.Char, error) {
	char := domain.Char{ID: charID, GameID: gameID}
	err := r.db.NewSelect().
		Model(&char).
		WherePK().
		Relation("Records").
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &char, nil
}

// func (s *Storage) CreateChar(charCreate *reqData.CharCreate, player *Player) (*Char, error) {
// 	if charCreate.Name == "" {
// 		return nil, fmt.Errorf("name cannot be empty")
// 	}

// 	var hiddenBy = 0
// 	if charCreate.Hidden {
// 		hiddenBy = player.ID
// 	}

// 	char := Char{
// 		Name:        charCreate.Name,
// 		Title:       charCreate.Title,
// 		Description: charCreate.Description,
// 		HiddenBy:    hiddenBy,
// 		PlayerID:    player.ID,
// 		GameID:      player.CurrentGameID,
// 	}

// 	_, err := s.db.NewInsert().Model(&char).
// 		Column("name", "title", "description", "hidden_by", "player_id", "game_id").
// 		Returning("*").Exec(context.Background(), &char)
// 	//Exec(context.Background())

// 	return &char, err
// }

// func (s *Storage) UpdateChar(charUpdate *reqData.CharUpdate, char *Char, player *Player) (*Char, error) {
// 	if charUpdate.Name == "" {
// 		return nil, fmt.Errorf("name cannot be empty")
// 	}

// 	var hiddenBy = 0
// 	if charUpdate.Hidden {
// 		hiddenBy = player.ID
// 	}

// 	_, err := s.db.NewUpdate().Model(char).WherePK().
// 		Set("name = ?", charUpdate.Name).
// 		Set("title = ?", charUpdate.Title).
// 		Set("description = ?", charUpdate.Description).
// 		Set("hidden_by = ?", hiddenBy).
// 		Returning("*").Exec(context.Background())
// 	return char, err
// }

func (r *EntitiesRepo) GetCurrentGameNPCByID(ctx context.Context, gameID, npcID int) (*domain.NPC, error) {
	npc := domain.NPC{ID: npcID, GameID: gameID}
	err := r.db.NewSelect().
		Model(&npc).
		WherePK().
		Relation("Records").
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &npc, nil
}

func (r *EntitiesRepo) GetCurrentGameLocationByID(ctx context.Context, gameID, locationID int) (*domain.Location, error) {
	location := domain.Location{ID: locationID, GameID: gameID}
	err := r.db.NewSelect().
		Model(&location).
		WherePK().
		Relation("Records").
		Relation("Parent"). // ** PARENT is not guarded if hidden ** //
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &location, nil
}

func (r *EntitiesRepo) GetCurrentGameLocationChildrenByID(ctx context.Context, gameID, locationID int) ([]domain.Location, error) {
	var locations []domain.Location
	err := r.db.NewSelect().
		Model(&locations).
		Where("game_id = ? AND pid = ?", gameID, locationID).
		Scan(ctx, &locations)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *EntitiesRepo) CreateChar(ctx context.Context, char *domain.Char) (*domain.Char, error) {
	_, err := r.db.NewInsert().
		Model(char).
		Column("name", "title", "description", "hidden_by", "player_id", "game_id").
		Returning("*").
		Exec(ctx, char)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return char, nil
}

func (r *EntitiesRepo) CreateNPC(ctx context.Context, npc *domain.NPC) (*domain.NPC, error) {
	_, err := r.db.NewInsert().
		Model(npc).
		Column("name", "title", "description", "hidden_by", "created_by_id", "game_id").
		Returning("*").
		Exec(ctx, npc)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return npc, nil
}

func (r *EntitiesRepo) CreateLocation(ctx context.Context, location *domain.Location) (*domain.Location, error) {
	_, err := r.db.NewInsert().
		Model(location).
		Column("name", "title", "description", "hidden_by", "created_by_id", "game_id", "pid").
		Returning("*").
		Exec(ctx, location)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return location, nil
}

func (r *EntitiesRepo) EditChar(ctx context.Context, charUpdate *dto.CharUpdate, playerID int) (*domain.Char, error) {
	var char domain.Char

	_, err := r.db.NewUpdate().
		Model(&char).
		Where("id = ?", charUpdate.ID).
		Set("name = ?", charUpdate.Name).
		Set("title = ?", charUpdate.Title).
		Set("description = ?", charUpdate.Description).
		Set("hidden_by = ?", gewiutils.TernaryInt(charUpdate.Hidden, playerID, 0)).
		Returning("*").
		Exec(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &char, nil
}

func (r *EntitiesRepo) EditNPC(ctx context.Context, npcUpdate *dto.NPCUpdate, playerID int) (*domain.NPC, error) {
	var npc domain.NPC

	_, err := r.db.NewUpdate().
		Model(&npc).
		Where("id = ?", npcUpdate.ID).
		Set("name = ?", npcUpdate.Name).
		Set("title = ?", npcUpdate.Title).
		Set("description = ?", npcUpdate.Description).
		Set("hidden_by = ?", gewiutils.TernaryInt(npcUpdate.Hidden, playerID, 0)).
		Returning("*").
		Exec(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &npc, nil
}

func (r *EntitiesRepo) EditLocation(ctx context.Context, locationUpdate *dto.LocationUpdate, playerID int) (*domain.Location, error) {
	var location domain.Location

	_, err := r.db.NewUpdate().
		Model(&location).
		Where("id = ?", locationUpdate.ID).
		Set("name = ?", locationUpdate.Name).
		Set("title = ?", locationUpdate.Title).
		Set("description = ?", locationUpdate.Description).
		Set("hidden_by = ?", gewiutils.TernaryInt(locationUpdate.Hidden, playerID, 0)).
		Set("pid = ?", locationUpdate.ParentID).
		Returning("*").
		Exec(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &location, nil
}

func (r *EntitiesRepo) GetCurrentGameSuggestionList(ctx context.Context, gameID, playerID int) ([]dto.Suggestion, error) {
	suggestions := []dto.Suggestion{}

	err := r.db.NewRaw(
		`SELECT
			id,
			CONCAT('char:', id) as sid,
			'char' as type,
			name,
			CASE
				WHEN hidden_by != 0 AND hidden_by != ?0 THEN false
				ELSE true
			END as secret
		FROM char
		WHERE game_id = ?1

		UNION ALL

		SELECT
			id,
			CONCAT('npc:', id) as sid,
			'npc' as type,
			name,
			CASE
				WHEN hidden_by = 0 OR hidden_by = ?0 THEN false
				ELSE true
			END as secret
		FROM npc
		WHERE game_id = ?1

		UNION ALL

		SELECT
			id,
			CONCAT('location:', id) as sid,
			'location' as type,
			name,
			CASE
				WHEN hidden_by = 0 OR hidden_by = ?0 THEN false
				ELSE true
			END as secret
		FROM location
		WHERE game_id = ?1`,
		playerID, gameID,
	).Scan(ctx, &suggestions)

	return suggestions, err
}
