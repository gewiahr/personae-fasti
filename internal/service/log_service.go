package service

import (
	"context"
	"time"

	"personae-fasti/internal/domain"
	"personae-fasti/internal/repo"
)

type LogService struct {
	repo repo.LogRepository
}

func NewLogService(repo repo.LogRepository) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) InsertLog(log *domain.ApiLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.repo.Insert(ctx, log)
}

func (s *LogService) Prune(successDays, errorDays int) error {
	var successBefore, errorBefore *time.Time
	now := time.Now().UTC()
	if successDays > 0 {
		value := now.AddDate(0, 0, -successDays)
		successBefore = &value
	}
	if errorDays > 0 {
		value := now.AddDate(0, 0, -errorDays)
		errorBefore = &value
	}
	if successBefore == nil && errorBefore == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.repo.Prune(ctx, successBefore, errorBefore)
}
