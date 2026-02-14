package service

import "net/http"

// ServiceError is a typed error with HTTP status code.
type ServiceError struct {
	Status  int
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

func ErrNotFound(msg string) *ServiceError {
	return &ServiceError{Status: http.StatusNotFound, Message: msg}
}

func ErrConflict(msg string) *ServiceError {
	return &ServiceError{Status: http.StatusConflict, Message: msg}
}

func ErrUnauthorized(msg string) *ServiceError {
	return &ServiceError{Status: http.StatusUnauthorized, Message: msg}
}

func ErrBadRequest(msg string) *ServiceError {
	return &ServiceError{Status: http.StatusBadRequest, Message: msg}
}

func ErrForbidden(msg string) *ServiceError {
	return &ServiceError{Status: http.StatusForbidden, Message: msg}
}
