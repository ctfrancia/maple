package response

import (
	"encoding/json"
	"go.uber.org/zap"
	"maps"
	"net/http"
)

type Helper struct {
	logger *zap.Logger
}

func NewHelper(logger *zap.Logger) *Helper {
	return &Helper{
		logger: logger,
	}
}

type envelope map[string]any

func (h *Helper) WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// js, err := json.MarshalIndent(data, "", "\t")
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	maps.Copy(w.Header(), headers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (h *Helper) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}
	err := h.WriteJSON(w, status, env, nil)
	if err != nil {
		h.logError(r, err)
		w.WriteHeader(500)
	}
}

func (h *Helper) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errs map[string]string) {
	env := envelope{"errors": errs}
	h.ErrorResponse(w, r, http.StatusUnprocessableEntity, env)
}

func (h *Helper) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (h *Helper) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	h.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (h *Helper) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid credentials"
	h.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (h *Helper) ConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "a record already exists with this email address"
	h.ErrorResponse(w, r, http.StatusConflict, message)
}

func (h *Helper) logError(r *http.Request, err error) {
	fields := []zap.Field{
		zap.String("method", r.Method),
		zap.String("uri", r.URL.RequestURI()),
	}
	h.logger.Error(err.Error(), fields...)
}
