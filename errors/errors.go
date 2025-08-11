package errors

import "net/http"

type AppError struct {
	StatusCode int
	Message    string
	Debug      string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(statusCode int, message string, debug string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Debug:      debug,
	}
}

var (
	ErrInvalidID       = New(http.StatusBadRequest, "invalid id", "ID must be a positive integer")
	ErrInvalidInput    = New(http.StatusBadRequest, "invalid input", "Malformed or missing request fields")
	ErrNotFound        = New(http.StatusNotFound, "resource not found", "Entity not found in database")
	ErrInternal        = New(http.StatusInternalServerError, "internal server error", "Unexpected system error")
	ErrConflict        = New(http.StatusConflict, "service name conflict", "Service with same name already exists")
	ErrServiceNotFound = New(http.StatusNotFound, "service not found", "Service with give id doesnt exists")
)
