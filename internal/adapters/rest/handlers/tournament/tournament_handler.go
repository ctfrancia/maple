// Package tournamenthandlers are the handlers for the tournament api
package tournamenthandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ctfrancia/maple/internal/adapters/rest/handlers/dto/tournament"
	"github.com/ctfrancia/maple/internal/adapters/rest/handlers/validator"
	"github.com/ctfrancia/maple/internal/adapters/rest/response"
	"github.com/ctfrancia/maple/internal/core/ports"
	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
)

type TournamentHandler struct {
	service   ports.TournamentServicer
	response  ports.SystemResponder
	logger    ports.Logger
	validator ports.ValidatorServicer
}

func NewTournamentHandler(log ports.Logger, ts ports.TournamentServicer) ports.TournamentHandler {
	handler := &TournamentHandler{
		service:   ts,
		response:  response.NewResponseWriter(log),
		logger:    log,
		validator: validator.NewValidator(),
	}

	return handler
}

// CreateTournamentHandler is the entrypoint for creating a tournament
func (h *TournamentHandler) CreateTournamentHandler(w http.ResponseWriter, r *http.Request) {
	var createTournamentRequest dto.CreateTournamentRequest
	if err := json.NewDecoder(r.Body).Decode(&createTournamentRequest); err != nil {
		h.response.ErrorResponse(w, r, http.StatusBadRequest, "invalid body")
		return
	}

	v := validator.NewValidator()
	v.Check(createTournamentRequest.Name != "", "name", "name is required")
	if !v.Valid() {
		h.response.FailedValidationResponse(w, r, v.ReturnErrors())
	}

	tournament := mapTournamentToDomain(createTournamentRequest)

	result, err := h.service.CreateTournament(r.Context(), tournament)
	if err != nil {
		h.response.ServerErrorResponse(w, r, err)
		return
	}

	resp := mapTournamentToDto(result)
	env := map[string]dto.Tournament{
		"tournament": resp,
	}

	// Successful response
	h.response.WriteJSON(w, http.StatusCreated, env, nil)
}

// FindTournamentHandler is the entrypoint for finding a tournament
func (h *TournamentHandler) FindTournamentHandler(w http.ResponseWriter, r *http.Request) {
	tournamentID := strings.TrimSpace(chi.URLParam(r, "id"))
	if tournamentID == "" {
		h.response.ErrorResponse(w, r, http.StatusBadRequest, "id is required")
		return
	}

	id, err := uuid.Parse(tournamentID)
	if err != nil {
		h.response.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.FindTournament(r.Context(), id)
	if err != nil {
		h.response.ServerErrorResponse(w, r, err)
		return
	}

	tournament := mapTournamentToDto(result)

	env := map[string]dto.Tournament{
		"tournament": tournament,
	}

	// Successful response
	h.response.WriteJSON(w, http.StatusOK, env, nil)
}

// ListTournamentsHandler is the entrypoint for listing tournaments
func (h *TournamentHandler) ListTournamentsHandler(w http.ResponseWriter, r *http.Request) {
}

// UpdateTournamentHandler is the entrypoint for updating a tournament
func (h *TournamentHandler) UpdateTournamentHandler(w http.ResponseWriter, r *http.Request) {
}

// DeleteTournamentHandler is the entrypoint for HARD deleting a tournament
func (h *TournamentHandler) DeleteTournamentHandler(w http.ResponseWriter, r *http.Request) {
}
