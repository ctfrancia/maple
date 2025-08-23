// Package inmemory provides a inmemory implementation of the tournament repository
// this is used for testing purposes
package inmemory

import (
	"time"

	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/ctfrancia/maple/internal/core/ports"
	"github.com/google/uuid"
)

type InMemoryTournamentRepository struct {
	tournaments map[uuid.UUID]domain.Tournament
}

func NewInMemoryTournamentRepository() ports.TournamentRepository {
	return &InMemoryTournamentRepository{
		tournaments: make(map[uuid.UUID]domain.Tournament),
	}
}

func (ir *InMemoryTournamentRepository) CreateTournament(tournament domain.Tournament) (domain.Tournament, error) {
	tournament.ID = len(ir.tournaments) + 1
	tournament.PublicID = uuid.New()
	tournament.CreatedAt = time.Now()
	tournament.UpdatedAt = time.Now()

	ir.tournaments[tournament.PublicID] = tournament

	return tournament, nil
}

func (ir *InMemoryTournamentRepository) FindTournament(id uuid.UUID) (domain.Tournament, error) {
	found, ok := ir.tournaments[id]
	if !ok {
		return domain.Tournament{}, domain.ErrTournamentNotFound
	}

	return found, nil
}

func (ir *InMemoryTournamentRepository) ListTournaments(params any) ([]domain.Tournament, error) {
	tournaments := make([]domain.Tournament, 0, len(ir.tournaments))
	for _, tournament := range ir.tournaments {
		tournaments = append(tournaments, tournament)
	}

	return tournaments, nil
}
