package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"

	"personae-fasti/internal/domain"
	"personae-fasti/internal/dto"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/validation"
	"personae-fasti/internal/repo"
)

const externalImageURLMaxLength = 2048

type ImageService struct {
	images   repo.ImageRepository
	entities repo.EntitiesRepository
	storage  ImageObjectStorage
}

type ImageObjectStorage interface {
	Configured() bool
	Put(ctx context.Context, key, contentType string, data []byte) error
	Delete(ctx context.Context, key string) error
	PublicURL(key string) string
}

func NewImageService(images repo.ImageRepository, entities repo.EntitiesRepository, storage ImageObjectStorage) *ImageService {
	return &ImageService{images: images, entities: entities, storage: storage}
}

func (s *ImageService) List(ctx context.Context, player *domain.Player, entityType, entityExt string) ([]dto.ImageInfo, error) {
	entityID, hiddenBy, err := s.resolveEntityByExt(ctx, player.CurrentGameID, entityType, entityExt)
	if err != nil {
		return nil, err
	}
	if err := ensureHiddenContentEditable(hiddenBy, player.ID); err != nil {
		return nil, err
	}

	images, err := s.images.ListByEntity(ctx, player.CurrentGameID, entityType, entityID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения изображений", err)
	}
	result := make([]dto.ImageInfo, len(images))
	for i := range images {
		result[i] = s.imageToInfo(&images[i])
	}
	return result, nil
}

func (s *ImageService) CreateExternal(ctx context.Context, player *domain.Player, entityType, entityExt, rawURL string) (*dto.ImageInfo, error) {
	entityID, hiddenBy, err := s.resolveEntityByExt(ctx, player.CurrentGameID, entityType, entityExt)
	if err != nil {
		return nil, err
	}
	if err := ensureHiddenContentEditable(hiddenBy, player.ID); err != nil {
		return nil, err
	}

	externalURL, err := validateExternalImageURL(rawURL)
	if err != nil {
		return nil, err
	}
	image, err := s.images.CreateExternal(ctx, &domain.Image{
		EntityType:  entityType,
		EntityID:    entityID,
		GameID:      player.CurrentGameID,
		UploadedBy:  player.ID,
		SourceType:  domain.ImageSourceExternal,
		ExternalURL: externalURL,
		Status:      domain.ImageStatusComplete,
	})
	if err != nil {
		return nil, e.NewInternalError("Ошибка сохранения изображения", err)
	}
	info := s.imageToInfo(image)
	return &info, nil
}

func (s *ImageService) Upload(ctx context.Context, player *domain.Player, entityType, entityExt string, reader io.Reader, fileSize int64) (*dto.ImageInfo, error) {
	entityID, hiddenBy, err := s.resolveEntityByExt(ctx, player.CurrentGameID, entityType, entityExt)
	if err != nil {
		return nil, err
	}
	if err := ensureHiddenContentEditable(hiddenBy, player.ID); err != nil {
		return nil, err
	}
	quota, err := s.images.GetQuota(ctx, player.CurrentGameID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения квоты изображений", err)
	}
	if quota.MaxBytes <= 0 {
		return nil, e.NewForbiddenError("Загрузка файлов для этой игры недоступна")
	}
	if fileSize > quota.MaxFileBytes {
		return nil, e.NewFieldValidationError(map[string]string{"file": "Файл превышает допустимый размер"})
	}
	if s.storage == nil || !s.storage.Configured() {
		return nil, e.NewInternalError("Хранилище изображений не настроено", fmt.Errorf("image storage is not configured"))
	}

	processed, err := processUploadedImage(reader, quota.MaxFileBytes)
	if err != nil {
		return nil, e.NewFieldValidationError(map[string]string{"file": err.Error()})
	}
	imageExt, err := randomImageExt()
	if err != nil {
		return nil, e.NewInternalError("Ошибка создания идентификатора изображения", err)
	}
	if player.CurrentGame == nil || player.CurrentGame.ExtID == "" {
		return nil, e.NewInternalError("Текущая игра не загружена", fmt.Errorf("current game ext is missing"))
	}
	baseKey := fmt.Sprintf("img/%s/%s/%s/%s", player.CurrentGame.ExtID, entityType, entityExt, imageExt)
	image := &domain.Image{
		ExtID:       imageExt,
		EntityType:  entityType,
		EntityID:    entityID,
		GameID:      player.CurrentGameID,
		UploadedBy:  player.ID,
		SourceType:  domain.ImageSourceUploaded,
		StorageKey:  baseKey + ".webp",
		ThumbKey:    baseKey + "_thumb.webp",
		ContentType: "image/webp",
		ByteSize:    int64(len(processed.full) + len(processed.thumbnail)),
		Width:       processed.width,
		Height:      processed.height,
		Checksum:    processed.checksum,
		Status:      domain.ImageStatusPending,
	}
	if err := s.images.CreatePendingUpload(ctx, image); err != nil {
		return nil, s.uploadRepoError(err)
	}

	abort := func() {
		cleanupContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()
		_ = s.storage.Delete(cleanupContext, image.ThumbKey)
		_ = s.storage.Delete(cleanupContext, image.StorageKey)
		_ = s.images.AbortUpload(cleanupContext, image)
	}
	if err := s.storage.Put(ctx, image.StorageKey, image.ContentType, processed.full); err != nil {
		abort()
		return nil, e.NewInternalError("Ошибка загрузки изображения", err)
	}
	if err := s.storage.Put(ctx, image.ThumbKey, image.ContentType, processed.thumbnail); err != nil {
		abort()
		return nil, e.NewInternalError("Ошибка загрузки миниатюры", err)
	}
	image, err = s.images.CompleteUpload(ctx, image)
	if err != nil {
		abort()
		return nil, e.NewInternalError("Ошибка завершения загрузки изображения", err)
	}
	info := s.imageToInfo(image)
	return &info, nil
}

func (s *ImageService) uploadRepoError(err error) error {
	switch err {
	case repo.ErrUploadDisabled:
		return e.NewForbiddenError("Загрузка файлов для этой игры недоступна")
	case repo.ErrQuotaExceeded:
		return e.NewFieldValidationError(map[string]string{"file": "Недостаточно места в хранилище игры"})
	case repo.ErrImageLimit:
		return e.NewFieldValidationError(map[string]string{"file": "Достигнут лимит изображений игры"})
	default:
		return e.NewInternalError("Ошибка резервирования места для изображения", err)
	}
}

func (s *ImageService) SetMain(ctx context.Context, player *domain.Player, imageExt string) (*dto.ImageInfo, error) {
	image, err := s.getEditableImage(ctx, player, imageExt)
	if err != nil {
		return nil, err
	}
	image, err = s.images.SetMain(ctx, image)
	if err != nil {
		if err == repo.ErrNotFound {
			return nil, e.NewNotFoundError("Изображение не найдено")
		}
		return nil, e.NewInternalError("Ошибка изменения главного изображения", err)
	}
	info := s.imageToInfo(image)
	return &info, nil
}

func (s *ImageService) Delete(ctx context.Context, player *domain.Player, imageExt string) error {
	image, err := s.getEditableImage(ctx, player, imageExt)
	if err != nil {
		return err
	}
	if err := s.images.SoftDelete(ctx, image); err != nil {
		if err == repo.ErrNotFound {
			return e.NewNotFoundError("Изображение не найдено")
		}
		return e.NewInternalError("Ошибка удаления изображения", err)
	}
	return nil
}

func (s *ImageService) GetQuota(ctx context.Context, player *domain.Player) (*dto.GameImageQuota, error) {
	quota, err := s.images.GetQuota(ctx, player.CurrentGameID)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения квоты изображений", err)
	}
	return &dto.GameImageQuota{
		MaxBytes:      quota.MaxBytes,
		UsedBytes:     quota.UsedBytes,
		ReservedBytes: quota.ReservedBytes,
		MaxFileBytes:  quota.MaxFileBytes,
		MaxImages:     quota.MaxImages,
		UploadEnabled: quota.MaxBytes > 0,
	}, nil
}

func (s *ImageService) getEditableImage(ctx context.Context, player *domain.Player, imageExt string) (*domain.Image, error) {
	if imageExt == "" {
		return nil, e.NewNotFoundError("Изображение не найдено")
	}
	image, err := s.images.GetByExt(ctx, player.CurrentGameID, imageExt)
	if err != nil {
		return nil, e.NewInternalError("Ошибка получения изображения", err)
	}
	if image == nil {
		return nil, e.NewNotFoundError("Изображение не найдено")
	}
	_, hiddenBy, err := s.resolveEntityByID(ctx, player.CurrentGameID, image.EntityType, image.EntityID)
	if err != nil {
		return nil, err
	}
	if err := ensureHiddenContentEditable(hiddenBy, player.ID); err != nil {
		return nil, err
	}
	return image, nil
}

func (s *ImageService) resolveEntityByExt(ctx context.Context, gameID int, entityType, entityExt string) (int, int, error) {
	switch entityType {
	case "char":
		entity, err := s.entities.GetCurrentGameCharByExt(ctx, gameID, entityExt)
		if err != nil {
			return 0, 0, e.NewInternalError("Ошибка получения персонажа", err)
		}
		if entity == nil {
			return 0, 0, e.NewNotFoundError("Персонаж не найден")
		}
		return entity.ID, entity.HiddenBy, nil
	case "npc":
		entity, err := s.entities.GetCurrentGameNPCByExt(ctx, gameID, entityExt)
		if err != nil {
			return 0, 0, e.NewInternalError("Ошибка получения персонажа", err)
		}
		if entity == nil {
			return 0, 0, e.NewNotFoundError("Персонаж не найден")
		}
		return entity.ID, entity.HiddenBy, nil
	case "location":
		entity, err := s.entities.GetCurrentGameLocationByExt(ctx, gameID, entityExt)
		if err != nil {
			return 0, 0, e.NewInternalError("Ошибка получения места", err)
		}
		if entity == nil {
			return 0, 0, e.NewNotFoundError("Место не найдено")
		}
		return entity.ID, entity.HiddenBy, nil
	default:
		return 0, 0, e.NewFieldValidationError(map[string]string{"entityType": "Неизвестный тип сущности"})
	}
}

func (s *ImageService) resolveEntityByID(ctx context.Context, gameID int, entityType string, entityID int) (int, int, error) {
	switch entityType {
	case "char":
		entity, err := s.entities.GetCurrentGameCharByID(ctx, gameID, entityID)
		if err != nil {
			return 0, 0, e.NewInternalError("Ошибка получения персонажа", err)
		}
		if entity == nil {
			return 0, 0, e.NewNotFoundError("Персонаж не найден")
		}
		return entity.ID, entity.HiddenBy, nil
	case "npc":
		entity, err := s.entities.GetCurrentGameNPCByID(ctx, gameID, entityID)
		if err != nil {
			return 0, 0, e.NewInternalError("Ошибка получения персонажа", err)
		}
		if entity == nil {
			return 0, 0, e.NewNotFoundError("Персонаж не найден")
		}
		return entity.ID, entity.HiddenBy, nil
	case "location":
		entity, err := s.entities.GetCurrentGameLocationByID(ctx, gameID, entityID)
		if err != nil {
			return 0, 0, e.NewInternalError("Ошибка получения места", err)
		}
		if entity == nil {
			return 0, 0, e.NewNotFoundError("Место не найдено")
		}
		return entity.ID, entity.HiddenBy, nil
	default:
		return 0, 0, e.NewNotFoundError("Сущность изображения не найдена")
	}
}

func validateExternalImageURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", e.NewFieldValidationError(map[string]string{"url": "Введите ссылку на изображение"})
	}
	if validation.CharacterCount(rawURL) > externalImageURLMaxLength {
		return "", e.NewFieldValidationError(map[string]string{"url": "Ссылка не может быть длиннее 2048 символов"})
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" || parsed.User != nil {
		return "", e.NewFieldValidationError(map[string]string{"url": "Введите абсолютную HTTPS-ссылку"})
	}
	hostname := strings.ToLower(strings.TrimSuffix(parsed.Hostname(), "."))
	if hostname == "localhost" || strings.HasSuffix(hostname, ".localhost") {
		return "", e.NewFieldValidationError(map[string]string{"url": "Локальные адреса недоступны"})
	}
	if ip := net.ParseIP(hostname); ip != nil && (ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified()) {
		return "", e.NewFieldValidationError(map[string]string{"url": "Локальные адреса недоступны"})
	}
	parsed.Fragment = ""
	return parsed.String(), nil
}

func (s *ImageService) imageToInfo(image *domain.Image) dto.ImageInfo {
	imageURL, thumbnailURL := image.ExternalURL, image.ExternalURL
	if image.SourceType == domain.ImageSourceUploaded && s.storage != nil {
		imageURL = s.storage.PublicURL(image.StorageKey)
		thumbnailURL = s.storage.PublicURL(image.ThumbKey)
	}
	return dto.ImageInfo{
		ExtID:        image.ExtID,
		URL:          imageURL,
		ThumbnailURL: thumbnailURL,
		Width:        image.Width,
		Height:       image.Height,
		IsMain:       image.IsMain,
		SourceType:   string(image.SourceType),
		Created:      image.Created,
	}
}

func randomImageExt() (string, error) {
	data := make([]byte, 12)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}
