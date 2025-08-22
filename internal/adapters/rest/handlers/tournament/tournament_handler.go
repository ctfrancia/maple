// Package tournamenthandlers are the handlers for the tournament api
package tournamenthandlers

import (
	"encoding/json"
	"net/http"

	"github.com/ctfrancia/maple/internal/adapters/rest/handlers/dto/tournament"
	"github.com/ctfrancia/maple/internal/adapters/rest/handlers/validator"
	"github.com/ctfrancia/maple/internal/adapters/rest/response"
	"github.com/ctfrancia/maple/internal/core/ports"
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

func (h *TournamentHandler) CreateTournamentHandler(w http.ResponseWriter, r *http.Request) {
	var createTournamentRequest dto.CreateTournamentRequest
	if err := json.NewDecoder(r.Body).Decode(&createTournamentRequest); err != nil {
		h.response.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}
	// TODO: validate request
	// h.logger.Info(r.Context(), "create tournament request received", createTournamentRequest)
	isValidRequest := isValidRequest(createTournamentRequest)
	if !isValidRequest {
		h.response.FailedValidationResponse(w, r, h.validator.ReturnErrors())
		return
	}
	tournament := mapTournamentToDomain(createTournamentRequest)

	tournament, err := h.service.CreateTournament(r.Context(), tournament)
	if err != nil {
		h.response.ServerErrorResponse(w, r, err)
		return
	}

	// Successful response
	h.response.WriteJSON(w, http.StatusCreated, tournament, nil)
}

func (h *TournamentHandler) FindTournamentHandler(w http.ResponseWriter, r *http.Request) {
}

func (h *TournamentHandler) ListTournamentsHandler(w http.ResponseWriter, r *http.Request) {
}
