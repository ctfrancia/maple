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
	"fmt"
	"net/http"

	"github.com/ctfrancia/maple/internal/core/ports"
)

type TournamentHandler struct {
	useCase   ports.TournamentUseCase
	logger    ports.Logger
	responses ports.ResponseHelper
}

// NewTournamentHandler creates a new TournamentHandler
func NewTournamentHandler(lg ports.Logger, uc ports.TournamentUseCase) ports.TournamentHandler {
	return &TournamentHandler{
		logger:  lg,
		useCase: uc,
	}
}

// GetTournamentsBasic handles the GET /tournaments_basic request
func (h *TournamentHandler) GetTournamentsBasic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetTournamentsBasic")
	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")
	fmt.Println("start ", startDate, " end ", endDate)

	tournaments, err := h.useCase.ProcessTournamentRequest()
	if err != nil {
		h.logger.Error(r.Context(), "Error processing tournament request", ports.Error("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{"tournaments": tournaments}
	h.responses.WriteJSON(w, http.StatusOK, data, nil)
}
