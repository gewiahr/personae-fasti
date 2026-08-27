package service

import (
	"context"
	"net/http"
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

func (s *LogService) InsertLog(playerID int, r *http.Request, requestBody, responseBody string, httpCode int, reqStartTime time.Time) error {
	completedAt := time.Now().UTC()
	completedIn := completedAt.Sub(reqStartTime).Milliseconds()

	errMsg := ""
	if httpCode >= http.StatusBadRequest {
		errMsg = responseBody
	}

	log := &domain.ApiLog{
		PlayerID: playerID,

		URI:     r.URL.Path,
		Method:  r.Method,
		Request: requestBody,

		Response: responseBody,
		Code:     httpCode,
		Error:    errMsg,
		Time:     completedIn,

		Created: completedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.repo.Insert(ctx, log)
}
