package middleware

import "net/http"

type AppError struct {
	Code    int               `json:"-"`
	Message string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func BadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}
func Unauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}
func NotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}
func Conflict(msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg}
}
func UnprocessableEntity(msg string) *AppError {
	return &AppError{Code: http.StatusUnprocessableEntity, Message: msg}
}
func (e *AppError) WithDetails(details map[string]string) *AppError {
	e.Details = details
	return e
}
