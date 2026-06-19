package service

import (
	"context"
	"encoding/json"
	"fmt"
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

func (s *LogService) InsertLog(playerID int, r *http.Request, reqBody []byte, response any, httpCode int, errMsg string, reqStartTime time.Time) error {
	completedAt := time.Now().UTC()
	completedIn := completedAt.Sub(reqStartTime).Milliseconds()

	responseString := ""
	if response != nil {
		respBytes, err := json.Marshal(response)
		if err != nil {
			responseString = fmt.Sprintf("error marshal response for log: %w", err)
		}
		responseString = string(respBytes)
	}

	log := &domain.ApiLog{
		PlayerID: playerID,

		URI:     r.URL.Path,
		Method:  r.Method,
		Request: string(reqBody),

		Response: responseString,
		Code:     httpCode,
		Error:    errMsg,
		Time:     completedIn,

		Created: completedAt,
	}

	// TODO: repo error handling
	return s.repo.Insert(context.Background(), log)
}
