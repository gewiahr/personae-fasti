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
			ae.Code = "not_found"
		case ErrValidation:
			ae.Status = http.StatusBadRequest
			ae.Code = "validation_error"
		case ErrUnauthorized:
			ae.Status = http.StatusUnauthorized
			ae.Code = "unauthorized"
		case ErrForbidden:
			ae.Status = http.StatusForbidden
			ae.Code = "forbidden"
		default:
			ae.Status = http.StatusInternalServerError
			ae.Code = "internal_error"
			ae.inner = appErr.Inner
		}
		return ae
	}
	return ApiError{
		Status:  http.StatusInternalServerError,
		Message: "Внутренняя ошибка сервера",
		Code:    "internal_error",
		inner:   err,
	}
}

type ApiError struct {
	Status    int               `json:"-"`
	Message   string            `json:"error"`
	Code      string            `json:"code,omitempty"`
	RequestID string            `json:"requestId,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
	inner     error             `json:"-"`
}

func (e ApiError) Respond(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	return json.NewEncoder(w).Encode(e)
}

func NewApiError(status int, msg string) ApiError {
	return ApiError{Status: status, Message: msg, Code: "http_error"}
}

func (e ApiError) Internal() error { return e.inner }
