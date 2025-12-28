package apierror

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"error"`
	Message    string `json:"message"`
	Cause      error  `json:"-"`
}

func (e *APIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *APIError) Unwrap() error {
	return e.Cause
}

func New(statusCode int, code string, message string, cause ...error) *APIError {
	var c error
	if len(cause) > 0 {
		c = cause[0]
	}

	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Cause:      c,
	}
}

func NewBadRequest(message string, cause ...error) *APIError {
	return New(http.StatusBadRequest, CodeBadRequest, message, cause...)
}

func NewNotFound(message string, cause ...error) *APIError {
	return New(http.StatusNotFound, CodeNotFound, message, cause...)
}

func NewUnauthorized(message string, cause ...error) *APIError {
	return New(http.StatusUnauthorized, CodeUnauthorized, message, cause...)
}

func NewForbidden(message string, cause ...error) *APIError {
	return New(http.StatusForbidden, CodeForbidden, message, cause...)
}

func NewInternalServer(message string, cause ...error) *APIError {
	msg := message
	if msg == "" {
		msg = "A internal error has been ocurred."
	}
	return New(http.StatusInternalServerError, CodeInternalServer, msg, cause...)
}
