package dto

import "time"

type ExternalImageCreate struct {
	URL string `json:"url"`
}

type ImageInfo struct {
	ExtID        string     `json:"ext"`
	URL          string     `json:"url"`
	ThumbnailURL string     `json:"thumbnailUrl"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	IsMain       bool       `json:"isMain"`
	SourceType   string     `json:"sourceType"`
	Created      *time.Time `json:"created"`
}

type ImageList struct {
	Images []ImageInfo `json:"images"`
}

type GameImageQuota struct {
	MaxBytes      int64 `json:"maxBytes"`
	UsedBytes     int64 `json:"usedBytes"`
	ReservedBytes int64 `json:"reservedBytes"`
	MaxFileBytes  int64 `json:"maxFileBytes"`
	MaxImages     int   `json:"maxImages"`
	UploadEnabled bool  `json:"uploadEnabled"`
}
