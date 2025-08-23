package ports

import (
	"context"
	"net/http"

	dto "github.com/ctfrancia/maple/internal/adapters/http/handlers/dto/tournament"
	commands "github.com/ctfrancia/maple/internal/application/commands/tournament"
	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/google/uuid"
)

// TournamentHandler is for our incomming http requests
type TournamentHandler interface {
	CreateTournamentHandler(w http.ResponseWriter, r *http.Request)
	FindTournamentHandler(w http.ResponseWriter, r *http.Request)
	ListTournamentsHandler(w http.ResponseWriter, r *http.Request)
	UpdateTournamentHandler(w http.ResponseWriter, r *http.Request)
	DeleteTournamentHandler(w http.ResponseWriter, r *http.Request)
}

// TournamentServicer is for our application layer
type TournamentServicer interface {
	CreateTournament(ctx context.Context, tournament commands.CreateTournamentCommand) (domain.Tournament, error)
	ListTournaments(ctx context.Context) ([]domain.Tournament, error)
	FindTournament(ctx context.Context, cmd commands.FindTournamentCommand) (domain.Tournament, error)
}

// TournamentRepository  is for our persistence layer
type TournamentRepository interface {
	CreateTournament(tournament domain.Tournament) (domain.Tournament, error)
	FindTournament(id uuid.UUID) (domain.Tournament, error)
	ListTournaments(params any) ([]domain.Tournament, error)
}

type TournamentMapper interface {
	MapToCommand(dto dto.CreateTournamentRequest) commands.CreateTournamentCommand
	MapToFindCommand(ID uuid.UUID) commands.FindTournamentCommand
}
