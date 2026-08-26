package errorutils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ErrToApiError(err error) ApiError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		ae := ApiError{Message: appErr.Message, Fields: appErr.Fields}
		switch appErr.Type {
		case ErrNotFound:
			ae.Status = http.StatusNotFound
		case ErrValidation:
			ae.Status = http.StatusBadRequest
		case ErrUnauthorized:
			ae.Status = http.StatusUnauthorized
		case ErrForbidden:
			ae.Status = http.StatusForbidden
		default:
			ae.Status = http.StatusInternalServerError
			ae.inner = appErr.Inner
		}
		return ae
	}
	return ApiError{
		Status:  http.StatusInternalServerError,
		Message: "internal error",
		inner:   err,
	}
}

type ApiError struct {
	Status  int               `json:"-"`
	Message string            `json:"error"`
	Fields  map[string]string `json:"fields,omitempty"`
	inner   error             `json:"-"`
}

func (e ApiError) Respond(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	return json.NewEncoder(w).Encode(e)
}

func NewApiError(status int, msg string) ApiError {
	return ApiError{Status: status, Message: msg}
}
