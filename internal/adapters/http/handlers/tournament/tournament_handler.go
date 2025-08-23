// Package tournamenthandlers are the handlers for the tournament api
package tournamenthandlers

import (
	"encoding/json"
	"net/http"
	"strings"

	dto "github.com/ctfrancia/maple/internal/adapters/http/handlers/dto/tournament"
	"github.com/ctfrancia/maple/internal/adapters/http/handlers/validator"
	"github.com/ctfrancia/maple/internal/adapters/http/response"
	"github.com/ctfrancia/maple/internal/core/ports"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TournamentHandler struct {
	service   ports.TournamentServicer
	response  ports.SystemResponder
	logger    ports.Logger
	validator ports.ValidatorServicer
	mapper    ports.TournamentMapper
}

func NewTournamentHandler(log ports.Logger, ts ports.TournamentServicer) ports.TournamentHandler {
	handler := &TournamentHandler{
		service:   ts,
		response:  response.NewResponseWriter(log),
		logger:    log,
		validator: validator.NewValidator(),
		mapper:    NewTournamentMapper(),
	}

	return handler
}

// CreateTournamentHandler is the entrypoint for creating a tournament
func (h *TournamentHandler) CreateTournamentHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Receive DTO from JSON
	var ctr dto.CreateTournamentRequest
	if err := json.NewDecoder(r.Body).Decode(&ctr); err != nil {
		h.response.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Map DTO to Command
	cmd := h.mapper.MapToCommand(ctr)

	// 3. Validate command
	if err := cmd.Validate(); err != nil {
		h.response.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Execute command via service
	result, err := h.service.CreateTournament(r.Context(), cmd)
	if err != nil {
		h.response.ServerErrorResponse(w, r, err)
		return
	}

	resp := mapTournamentToDto(result)

	env := map[string]dto.TournamentResponse{
		"tournament": resp,
	}

	// Successful response
	h.response.WriteJSON(w, http.StatusCreated, env, nil)
}

// FindTournamentHandler is the entrypoint for finding a tournament
func (h *TournamentHandler) FindTournamentHandler(w http.ResponseWriter, r *http.Request) {
	tournamentIDStr := strings.TrimSpace(chi.URLParam(r, "id"))
	if tournamentIDStr == "" {
		h.response.ErrorResponse(w, r, http.StatusBadRequest, "id is required")
		return
	}

	// Parse and validate the UUID here - fail fast
	ID, err := uuid.Parse(tournamentIDStr)
	if err != nil {
		h.response.ErrorResponse(w, r, http.StatusBadRequest, "invalid tournament ID format")
		return
	}

	cmd := h.mapper.MapToFindCommand(ID)

	result, err := h.service.FindTournament(r.Context(), cmd)
	if err != nil {
		h.response.ServerErrorResponse(w, r, err)
		return
	}

	tournament := mapTournamentToDto(result)

	env := map[string]dto.TournamentResponse{
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
