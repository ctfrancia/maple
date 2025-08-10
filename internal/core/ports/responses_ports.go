package ports

import (
	"net/http"
)

type Envelope map[string]any

type ResponseHelper interface {
	WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error
	ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any)
	FailedValidationResponse(w http.ResponseWriter, r *http.Request, errs map[string]string)
	BadRequestResponse(w http.ResponseWriter, r *http.Request, err error)
	ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error)
	InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request)
	ConflictResponse(w http.ResponseWriter, r *http.Request)
}
