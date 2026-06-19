package httputils

import (
	"context"
	"net/http"
	"net/url"
	"personae-fasti/internal/domain"
	"strconv"
)

type RequestData[Body any] struct {
	Context context.Context
	Body    Body
	Player  *domain.Player
	Query   url.Values
	Request *http.Request
}

// TODO: add error
func GetPathValueInt(r *http.Request, param string) int {
	wrongValue := -1

	value := r.PathValue(param)
	if len(value) != 0 {
		valueInt, err := strconv.Atoi(value)
		if err == nil {
			return valueInt
		}
	}

	return wrongValue
}
