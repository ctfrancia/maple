package ports

import (
	"context"
	"net/http"

	"github.com/ctfrancia/maple/internal/core/domain"
)

// TournamentHandler defines the interface for handling the incoming requests
type TournamentHandler interface {
	GetTournamentsBasic(w http.ResponseWriter, r *http.Request)
}

// TournamentServicer defines the interface for the business logic
type TournamentUseCase interface {
<<<<<<< HEAD
	ProcessTournamentRequest() ([]domain.Tournament, error)
=======
	//ProcessTournamentRequest() ([]domain.Tournament, error)
	CreateNewTournament() (domain.Tournament, error)
>>>>>>> 3e96be3 (wip: end of day commit)
}

// TournamentRepository defines the interface for interacting with the database
type TournamentRepository interface {
	//CreateTournament(ctx context.Context, tournament domain.Tournament) error
	GetTournaments(ctx context.Context, page int, size int) ([]domain.Tournament, error)
	//GetTournament(ctx context.Context, id int) (domain.Tournament, error)
	//UpdateTournament(ctx context.Context, tournament domain.Tournament) error
	//DeleteTournament(ctx context.Context, id int) error
}
