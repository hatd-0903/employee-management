package models

import "net/http"

// AppError is the standard error type returned by services and repositories.
// Handlers translate it into the JSON error envelope {"error": ..., "code": ...}.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func ErrValidation(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message, nil)
}

func ErrNotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message, nil)
}

func ErrConflict(message string) *AppError {
	return NewAppError(http.StatusConflict, message, nil)
}

func ErrUnauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message, nil)
}

func ErrInternal(err error) *AppError {
	return NewAppError(http.StatusInternalServerError, "internal server error", err)
}

func ErrTimeout(err error) *AppError {
	return NewAppError(http.StatusGatewayTimeout, "request timed out", err)
}
