// Package ports represents the interface for how the core interfaces with the rest of the world
package ports

import (
	"net/http"

	"github.com/ctfrancia/maple/internal/core/domain"
)

type SystemHandler interface {
	HealthHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	NewConsumerHandler(w http.ResponseWriter, r *http.Request)
}

type SystemServicer interface {
	ProcessSystemHealthRequest() domain.System
}

type SystemAdapter interface {
	GetSystemInfo() domain.System
}

type SystemResponder interface {
	WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error
	ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any)
	FailedValidationResponse(w http.ResponseWriter, r *http.Request, errs map[string]string)
	BadRequestResponse(w http.ResponseWriter, r *http.Request, err error)
	ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error)
	InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request)
	ConflictResponse(w http.ResponseWriter, r *http.Request)
}
