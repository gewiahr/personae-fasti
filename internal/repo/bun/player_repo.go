package bunrepo

import (
	"context"
	"database/sql"
	"personae-fasti/internal/domain"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/repo"

	"github.com/uptrace/bun"
)

var _ repo.PlayerRepository = (*PlayerRepo)(nil)

type PlayerRepo struct {
	db *bun.DB
}

func NewPlayerRepo(db *bun.DB) *PlayerRepo { return &PlayerRepo{db: db} }

func (r *PlayerRepo) GetByID(ctx context.Context, id int) (*domain.Player, error) {
	player := &domain.Player{ID: id}
	err := r.db.NewSelect().
		Model(player).
		WherePK().
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return player, nil
}

// TODO: add time and deletion verification
func (r *PlayerRepo) GetByToken(ctx context.Context, tokenHash string) (*domain.Player, error) {
	token := new(domain.Token)
	err := r.db.NewSelect().
		Model(token).
		Where("token_hash = ?", tokenHash).
		Relation("Player").
		Relation("Player.RegData").
		Relation("Player.CurrentGame").
		Relation("Player.CurrentGame.Settings"). // TODO: make proper current game load
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, e.NewNotFoundError("")
	}
	return token.Player, nil
}

func (r *PlayerRepo) GetByUsername(ctx context.Context, username string) (*domain.Player, error) {
	player := &domain.Player{}
	err := r.db.NewSelect().
		Model(player).
		Where("username = ?", username).
		Relation("RegData").
		Relation("CurrentGame").
		Relation("CurrentGame.Settings").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (r *PlayerRepo) CreatePlayer(ctx context.Context, player *domain.Player) (*domain.Player, error) {
	if err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(player).Returning("*").Exec(ctx); err != nil {
			return err
		}

		playerRegData := &domain.PlayerRegData{
			PlayerID:    player.ID,
			UsernameSet: true,
		}
		if _, err := tx.NewInsert().Model(playerRegData).Returning("*").Exec(ctx); err != nil {
			return err
		}

		player.RegData = playerRegData

		return nil
	}); err != nil {
		return nil, err
	}

	return player, nil
}

func (r *PlayerRepo) IsUsernameFree(ctx context.Context, username string) (bool, error) {
	count, err := r.db.NewSelect().Model(&domain.Player{}).Where("username = ?", username).Count(context.Background())
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *PlayerRepo) InsertToken(ctx context.Context, token *domain.Token) (*domain.Token, error) {
	if _, err := r.db.NewInsert().Model(token).Returning("*").Exec(context.Background(), token); err != nil {
		return nil, err
	}
	return token, nil
}

func (s *PlayerRepo) SetPlayerPassword(ctx context.Context, playerID int, passwordHash string) (*domain.Player, error) {
	player := &domain.Player{ID: playerID}
	if _, err := s.db.NewUpdate().Model(player).WherePK().Set("password_hash = ?", passwordHash).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return player, nil
}

func (r *PlayerRepo) GetPlayerWithGames(ctx context.Context, playerID int) (*domain.Player, error) {
	player := &domain.Player{ID: playerID}
	err := r.db.NewSelect().
		Model(player).
		WherePK().
		Relation("Games.GM").
		Relation("Invites.GM").
		Relation("CurrentGame.GM").
		Relation("CurrentGame.Settings").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (r *PlayerRepo) ChangeCurrentGame(ctx context.Context, playerID, gameID int) (*domain.Player, error) {
	player := &domain.Player{ID: playerID}
	_, err := r.db.NewUpdate().
		Model(player).
		WherePK().
		Set("current_game_id = ?", gameID).
		Returning("*").
		Exec(ctx, player)
	if err != nil {
		return nil, err
	}
	return player, nil
}

//GetInviteByExt

func (r *PlayerRepo) GetInvite(ctx context.Context, playerID int, inviteCode string) (*domain.GameInvite, error) {
	invite := &domain.GameInvite{
		PlayerID: playerID,
	}
	err := r.db.NewSelect().
		Model(invite).
		Where("player_id = ? AND code = ?", playerID, inviteCode).
		Scan(context.Background(), invite)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return invite, nil
}

func (r *PlayerRepo) AddPlayerToGame(ctx context.Context, playerID, gameID int) error {
	if err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		newPlayerGame := domain.PlayerGame{
			PlayerID: playerID,
			GameID:   gameID,
		}
		if _, err := tx.NewInsert().Model(&newPlayerGame).Returning("*").Exec(ctx, &newPlayerGame); err != nil {
			return err
		}

		gameInvite := domain.GameInvite{
			PlayerID: playerID,
			GameID:   gameID,
		}
		if _, err := tx.NewDelete().Model(&gameInvite).WherePK().Exec(context.Background()); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *PlayerRepo) DeleteInvite(ctx context.Context, invite *domain.GameInvite) error {
	if _, err := r.db.NewDelete().Model(invite).WherePK().Exec(context.Background()); err == sql.ErrNoRows {
		return err
	} else if err != nil {
		return err
	}

	return nil
}

// func (s *Storage) InvitePlayer(game *Game, player *Player) error {
// 	playerGame := &PlayerGame{
// 		PlayerID: player.ID,
// 		GameID:   game.ID,
// 	}
// 	if count, err := s.db.NewSelect().Model(playerGame).WherePK().Count(context.Background()); err != nil && err != sql.ErrNoRows {
// 		return err
// 	} else if count > 0 {
// 		return fmt.Errorf("player is a participant already")
// 	}

// 	invite := &GameInvite{
// 		PlayerID: player.ID,
// 		GameID:   game.ID,
// 	}
// 	if count, err := s.db.NewSelect().Model(invite).WherePK().Count(context.Background()); err != nil && err != sql.ErrNoRows {
// 		return err
// 	} else if count > 0 {
// 		return fmt.Errorf("player is invited already")
// 	}

// 	if _, err := s.db.NewInsert().Model(invite).Exec(context.Background(), invite); err != nil {
// 		return err
// 	}

// 	return nil
// }
