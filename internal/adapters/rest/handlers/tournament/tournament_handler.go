// Package tournamenthandlers are the handlers for the tournament api
package tournamenthandlers

import (
	"net/http"

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
}

func (h *TournamentHandler) FindTournamentHandler(w http.ResponseWriter, r *http.Request) {
}

func (h *TournamentHandler) ListTournamentsHandler(w http.ResponseWriter, r *http.Request) {
}
