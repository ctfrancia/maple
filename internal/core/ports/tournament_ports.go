package ports

import (
	"context"
	"net/http"

	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/google/uuid"
)

type TournamentHandler interface {
	CreateTournamentHandler(w http.ResponseWriter, r *http.Request)
	FindTournamentHandler(w http.ResponseWriter, r *http.Request)
	ListTournamentsHandler(w http.ResponseWriter, r *http.Request)
}

type TournamentServicer interface {
	CreateTournament(ctx context.Context, tournament domain.Tournament) (domain.Tournament, error)
}

type TournamentRepository interface {
	CreateTournament(tournament domain.Tournament) (domain.Tournament, error)
	FindTournament(id uuid.UUID) (domain.Tournament, error)
}
