package errorutils

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrValidation   = errors.New("validation")
	ErrInternal     = errors.New("internal")
)

type AppError struct {
	Type    error
	Message string
	Inner   error
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Inner }

func NewNotFoundError(msg string) *AppError {
	return &AppError{Type: ErrNotFound, Message: msg}
}
func NewValidationError(msg string) *AppError {
	return &AppError{Type: ErrValidation, Message: msg}
}
func NewUnauthorizedError(msg string) *AppError {
	return &AppError{Type: ErrUnauthorized, Message: msg}
}
func NewForbiddenError(msg string) *AppError {
	return &AppError{Type: ErrForbidden, Message: msg}
}
func NewInternalError(msg string, inner error) *AppError {
	return &AppError{Type: ErrInternal, Message: "internal error", Inner: inner}
}
