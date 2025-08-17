// Package inmemory provides a inmemory implementation of the tournament repository
// this is used for testing purposes
package inmemory

import (
	"sync"
	"time"

	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/ctfrancia/maple/internal/core/ports"
	"github.com/google/uuid"
)

// InMemoryTournamentRepository is an inmemory tournament repository
type InMemoryTournamentRepository struct {
	mu          sync.RWMutex
	tournaments map[uuid.UUID]domain.Tournament
}

// NewInMemoryTournamentRepository creates a new inmemory tournament repository
func NewInMemoryTournamentRepository() ports.TournamentRepository {
	return &InMemoryTournamentRepository{
		tournaments: make(map[uuid.UUID]domain.Tournament),
	}
}

func (ir *InMemoryTournamentRepository) CreateTournament(tournament domain.Tournament) (domain.Tournament, error) {
	ir.mu.Lock()
	defer ir.mu.Unlock()

	tournament.ID = len(ir.tournaments) + 1
	tournament.PublicID = uuid.New()
	tournament.CreatedAt = time.Now()
	tournament.UpdatedAt = time.Now()

	ir.tournaments[tournament.PublicID] = tournament

	return tournament, nil
}

func (ir *InMemoryTournamentRepository) FindTournament(id uuid.UUID) (domain.Tournament, error) {
	ir.mu.RLock()
	defer ir.mu.RUnlock()

	found, ok := ir.tournaments[id]
	if !ok {
		return domain.Tournament{}, domain.ErrTournamentNotFound
	}

	return found, nil
}

func (ir *InMemoryTournamentRepository) ListTournaments(params any) ([]domain.Tournament, error) {
	ir.mu.RLock()
	defer ir.mu.RUnlock()

	tournaments := make([]domain.Tournament, 0, len(ir.tournaments))
	for _, tournament := range ir.tournaments {
		tournaments = append(tournaments, tournament)
	}

	return tournaments, nil
}
