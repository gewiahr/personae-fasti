package repo

import (
	"context"
	"errors"
	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	"time"

	"github.com/uptrace/bun"
)

type Storage interface {
	GameRepo() GameRepository
	PlayerRepo() PlayerRepository
	RecordRepo() RecordRepository
	EntitiesRepo() EntitiesRepository
	LogRepo() LogRepository
	AppRepo() AppRepository
	ImageRepo() ImageRepository

	Migrate(ctx context.Context) error
	Close() error
}

type ImageRepository interface {
	ListByEntity(ctx context.Context, gameID int, entityType string, entityID int) ([]domain.Image, error)
	GetByExt(ctx context.Context, gameID int, imageExt string) (*domain.Image, error)
	CreateExternal(ctx context.Context, image *domain.Image) (*domain.Image, error)
	CreatePendingUpload(ctx context.Context, image *domain.Image) error
	CompleteUpload(ctx context.Context, image *domain.Image) (*domain.Image, error)
	AbortUpload(ctx context.Context, image *domain.Image) error
	SetMain(ctx context.Context, image *domain.Image) (*domain.Image, error)
	SoftDelete(ctx context.Context, image *domain.Image) error
	GetQuota(ctx context.Context, gameID int) (*domain.GameImageQuota, error)
}

var (
	ErrNotFound       = errors.New("not found")
	ErrDBInternal     = errors.New("db internal")
	ErrUploadDisabled = errors.New("image upload disabled")
	ErrQuotaExceeded  = errors.New("image quota exceeded")
	ErrImageLimit     = errors.New("image count limit exceeded")
)

type GameRepository interface {
	// Get Game object with Players, Settings and Sessions
	// Used for loading current games
	GetCurrentGame(ctx context.Context, playerCurrentGameID int) (*domain.Game, error)

	GetByExt(ctx context.Context, gameExt string) (*domain.Game, error)
	Create(ctx context.Context, game *domain.Game) (*domain.Game, error)
	UpdateByExt(ctx context.Context, game *domain.Game) (*domain.Game, error)

	GetCurrentGameSession(ctx context.Context, gameID int) (*domain.Session, error)
	GetGameSessionByNumber(ctx context.Context, gameID int, sessionNumber int) (*domain.Session, error)

	CreateNewSession(ctx context.Context, game *domain.Game) (*domain.Session, error)
	EditSession(ctx context.Context, game *domain.Game, sessionUpdate *dto.SessionUpdate) (*domain.Session, error)
	RemoveLastSession(ctx context.Context, game *domain.Game) error

	GetPlayerInvites(ctx context.Context, playerID int) ([]domain.GameInvite, error)
	GetGameInvites(ctx context.Context, gameID int) ([]domain.GameInvite, error)

	InvitePlayer(ctx context.Context, invite *domain.GameInvite) error
	DeleteInvite(ctx context.Context, invite *domain.GameInvite) error

	UpdateGameSettings(ctx context.Context, gameID int, settingsUpdate *dto.GameSettingsUpdate) (*domain.GameSettings, error)

	// Create(ctx context.Context, game *domain.Game) error
	// GetByID(ctx context.Context, id int) (*domain.Game, error)
}

type PlayerRepository interface {
	GetByID(ctx context.Context, id int) (*domain.Player, error)
	GetByToken(ctx context.Context, tokenHash string) (*domain.Player, error)
	GetByUsername(ctx context.Context, username string) (*domain.Player, error)
	GetPlayerWithGames(ctx context.Context, playerID int) (*domain.Player, error)
	CreatePlayer(ctx context.Context, player *domain.Player) (*domain.Player, error)
	IsUsernameFree(ctx context.Context, username string) (bool, error)
	InsertToken(ctx context.Context, token *domain.Token) (*domain.Token, error)
	SetPlayerPassword(ctx context.Context, playerID int, passwordHash string) (*domain.Player, error)
	ChangeCurrentGame(ctx context.Context, playerID, gameID int) (*domain.Player, error)
	GetPersonalNote(ctx context.Context, playerID int) (string, error)
	UpdatePersonalNote(ctx context.Context, playerID int, note string) (string, error)

	GetInvite(ctx context.Context, playerID int, inviteCode string) (*domain.GameInvite, error)
	//DeleteInvite(ctx context.Context, invite *domain.GameInvite) error
	AddPlayerToGame(ctx context.Context, playerID, gameID int) error
}

type RecordRepository interface {
	GetCurrentGameRecordList(ctx context.Context, gameID, playerID int) ([]domain.Record, error)
	//GetListByGame(ctx context.Context, gameID int) ([]domain.Record, error)
	GetRecord(ctx context.Context, playerID int, recordID int) (*domain.Record, error)
	PostRecord(ctx context.Context, record *domain.Record) (*domain.Record, error)
	EditRecord(ctx context.Context, recordUpdate *dto.RecordUpdate, playerID int) (*domain.Record, error)
	SoftDeleteRecord(ctx context.Context, playerID int, recordID int) error

	FilterAllowedRecords(ctx context.Context, records []domain.Record, playerID int) ([]domain.Record, error)

	// TODO: make tx param general sql one
	InsertMentionsForRecord(ctx context.Context, tx bun.Tx, record *domain.Record) error
	DeleteMentionsForRecord(ctx context.Context, tx bun.Tx, record *domain.Record) error
}

type EntitiesRepository interface {
	GetCurrentGameCharList(ctx context.Context, gameID, playerID int) ([]domain.Char, error)
	GetCurrentGameNPCList(ctx context.Context, gameID, playerID int) ([]domain.NPC, error)
	GetCurrentGameLocationList(ctx context.Context, gameID, playerID int) ([]domain.Location, error)

	GetCurrentGameCharByID(ctx context.Context, gameID, charID int) (*domain.Char, error)
	GetCurrentGameNPCByID(ctx context.Context, gameID, npcID int) (*domain.NPC, error)
	GetCurrentGameLocationByID(ctx context.Context, gameID, locationID int) (*domain.Location, error)
	GetCurrentGameCharByExt(ctx context.Context, gameID int, charExt string) (*domain.Char, error)
	GetCurrentGameNPCByExt(ctx context.Context, gameID int, npcExt string) (*domain.NPC, error)
	GetCurrentGameLocationByExt(ctx context.Context, gameID int, locationExt string) (*domain.Location, error)
	GetCurrentGameLocationChildrenByID(ctx context.Context, gameID, locationID int) ([]domain.Location, error)

	CreateChar(ctx context.Context, char *domain.Char) (*domain.Char, error)
	CreateNPC(ctx context.Context, npc *domain.NPC) (*domain.NPC, error)
	CreateLocation(ctx context.Context, location *domain.Location) (*domain.Location, error)

	EditChar(ctx context.Context, charUpdate *dto.CharUpdate, playerID, gameID int) (*domain.Char, error)
	EditNPC(ctx context.Context, npcUpdate *dto.NPCUpdate, playerID, gameID int) (*domain.NPC, error)
	EditLocation(ctx context.Context, locationUpdate *dto.LocationUpdate, playerID, gameID int) (*domain.Location, error)

	GetCurrentGameSuggestionList(ctx context.Context, gameID, playerID int) ([]dto.Suggestion, error) // TODO: make db-level suggestion objest
}

type QuestRepository interface {
	GetCurrentGameQuestList(ctx context.Context, gameID, playerID int) ([]domain.Quest, error)
	GetPlayerCurrentGameQuestByID(ctx context.Context, gameID, questID int) (*domain.Quest, error)
	GetPlayerCurrentGameQuestByExt(ctx context.Context, gameID int, questExt string) (*domain.Quest, error)
	CreatePlayerCurrentGameQuest(ctx context.Context, questCreateData *dto.QuestCreateData, playerID, gameID int) (*domain.Quest, error)
	EditPlayerCurrentGameQuest(ctx context.Context, questUpdateData *dto.QuestUpdateData, playerID, gameID int) (*domain.Quest, error)
	//DeleteQuest(ctx context.Context, questID int) error

	FinishPlayerCurrentGameQuest(ctx context.Context, questID, gameID, playerID int, successful bool) (*domain.Quest, error)
	ResetPlayerCurrentGameQuest(ctx context.Context, questID, gameID, playerID int) (*domain.Quest, error)

	GetTasksByQuest(ctx context.Context, quest *domain.Quest) ([]domain.QuestTask, error)
	// CreateTasks(ctx context.Context, tasks []domain.QuestTask) ([]domain.QuestTask, error)
	// EditTasks(ctx context.Context, tasksUpdate []dto.TaskUpdate, playerID int) ([]domain.QuestTask, error)
	UpdateTasks(ctx context.Context, tasksPatch []dto.TaskPatch, quest *domain.Quest) ([]domain.QuestTask, error)

	FilterAllowedTasks(ctx context.Context, tasks []domain.QuestTask, playerID int) ([]domain.QuestTask, error)
}

type LogRepository interface {
	Insert(ctx context.Context, log *domain.ApiLog) error
	Prune(ctx context.Context, successBefore, errorBefore *time.Time) error
}

type AppRepository interface {
	InsertFeedback(ctx context.Context, feedback *domain.ServiceFeedback) (*domain.ServiceFeedback, error)
	// Insert(ctx context.Context, log *domain.ApiLog) error
}
