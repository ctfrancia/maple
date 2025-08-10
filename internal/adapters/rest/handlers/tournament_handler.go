package handlers

// TODO: Implement this
// endpoints:
// - GET /tournaments_basic
// description:
// this gets the tournaments within a given time range
// for
// - GET /tournaments_detailed
// description:
// this gets the tournaments within a given time range with expanded information

import (
	"net/http"

	"github.com/ctfrancia/maple/internal/core/ports"
)

type TournamentHandler struct {
	useCase ports.TournamentUseCase
	logger  ports.Logger
}

// NewTournamentHandler creates a new TournamentHandler
func NewTournamentHandler(lg ports.Logger, uc ports.TournamentUseCase) ports.TournamentHandler {
	return &TournamentHandler{
		logger:  lg,
		useCase: uc,
	}
}

// GetTournamentsBasic handles the GET /tournaments_basic request
func (h *TournamentHandler) GetTournamentsBasic(w http.ResponseWriter, r *http.Request) {}
