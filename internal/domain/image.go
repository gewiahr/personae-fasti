package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type ImageSourceType string

const (
	ImageSourceUploaded ImageSourceType = "uploaded"
	ImageSourceExternal ImageSourceType = "external"
)

type ImageStatus string

const (
	ImageStatusPending  ImageStatus = "pending"
	ImageStatusComplete ImageStatus = "complete"
	ImageStatusDeleted  ImageStatus = "deleted"
)

type Image struct {
	bun.BaseModel `bun:"table:image"`

	ID          int64           `bun:"id,pk,autoincrement"`
	ExtID       string          `bun:"ext,unique,notnull,type:varchar(16),default:nanoid(16)"`
	EntityType  string          `bun:"entity_type,notnull"`
	EntityID    int             `bun:"entity_id,notnull"`
	GameID      int             `bun:"game_id,notnull"`
	UploadedBy  int             `bun:"uploaded_by,notnull"`
	SourceType  ImageSourceType `bun:"source_type,notnull"`
	StorageKey  string          `bun:"storage_key,nullzero"`
	ThumbKey    string          `bun:"thumb_key,nullzero"`
	ExternalURL string          `bun:"external_url,nullzero"`
	ContentType string          `bun:"content_type,nullzero"`
	ByteSize    int64           `bun:"byte_size,notnull,default:0"`
	Width       int             `bun:"width"`
	Height      int             `bun:"height"`
	Checksum    string          `bun:"checksum,nullzero"`
	IsMain      bool            `bun:"is_main,notnull,default:false"`
	Status      ImageStatus     `bun:"status,notnull,default:'pending'"`
	Created     *time.Time      `bun:"created,nullzero,notnull,default:current_timestamp"`
	Deleted     *time.Time      `bun:"deleted"`
}

type GameImageQuota struct {
	bun.BaseModel `bun:"table:game_image_quota"`

	GameID        int   `bun:"game_id,pk"`
	MaxBytes      int64 `bun:"max_bytes,notnull,default:0"`
	UsedBytes     int64 `bun:"used_bytes,notnull,default:0"`
	ReservedBytes int64 `bun:"reserved_bytes,notnull,default:0"`
	MaxFileBytes  int64 `bun:"max_file_bytes,notnull,default:5242880"`
	MaxImages     int   `bun:"max_images,notnull,default:0"`
}
