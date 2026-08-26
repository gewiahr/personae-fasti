package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Player struct {
	bun.BaseModel `bun:"table:player"`

	ID           int     `bun:"id,pk,autoincrement"`
	ExtID        string  `bun:"ext,unique,notnull,type:varchar(16),default:nanoid(12)"`
	Username     string  `bun:"username,unique,notnull"`
	Email        string  `bun:"email,notnull,type:varchar(255),default:''"`
	PasswordHash string  `bun:"password_hash,notnull,type:varchar(255),default:''"`
	PersonalNote string  `bun:"personal_note,notnull,default:''"`
	Tokens       []Token `bun:"rel:has-many,join:id=player_id"`

	TelegramID int64     `bun:"telegram_id"`
	Telegram   *Telegram `bun:"rel:belongs-to,join:telegram_id=id"`

	Chars   []Char `bun:"rel:has-many,join:id=player_id"`
	Games   []Game `bun:"m2m:players_games,join:Player=Game"`
	Invites []Game `bun:"m2m:game_invites,join:Player=Game"`

	Feedback []ServiceFeedback `bun:"rel:has-many,join:id=player_id"`

	//Records []Record `bun:"rel:has-many,join:id=game_id"`

	CurrentGameID int   `bun:"current_game_id"`
	CurrentGame   *Game `bun:"rel:belongs-to,join:current_game_id=id"`

	Registered *time.Time `bun:"registeredTime,nullzero,notnull,default:current_timestamp"`
	LastAction *time.Time `bun:"lastActionTime,nullzero,notnull,default:current_timestamp"`
	Deleted    *time.Time `bun:"deleted,default:null"`

	RegData *PlayerRegData `bun:"rel:has-one,join:id=player_id"`
}

type PlayerRegData struct {
	bun.BaseModel `bun:"table:player_reg_data"`

	PlayerID int     `bun:"player_id,pk"`
	Player   *Player `bun:"rel:belongs-to,join:player_id=id"`

	UsernameSet bool `bun:"username_set,default:false"`
}

type Telegram struct {
	bun.BaseModel `bun:"table:telegram"`

	ID       int64  `bun:"id,pk"`
	Username string `bun:"username,notnull,default:''"`
	Lang     string `bun:"lang,default:'en'"`
	PicURL   string `bun:"pic_url,notnull,default:''"`
}

type Token struct {
	bun.BaseModel `bun:"table:tokens"`

	ID        int64     `bun:"id,pk,autoincrement"`
	PlayerID  int       `bun:"player_id,notnull"`
	TokenHash string    `bun:"token_hash,unique,notnull"`
	ExpiresAt time.Time `bun:"expires_at,notnull"`
	CreatedAt time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	Revoked   bool      `bun:"revoked,default:false"`

	Player *Player `bun:"rel:belongs-to,join:player_id=id"`
}
