package inmemory

import (
	"sync"

	"github.com/ctfrancia/maple/internal/core/ports"
)

type InMemoryTournamentProvider struct {
	repository ports.TournamentRepository
	mu         *sync.RWMutex
}

func NewTournamentRepositoryProvider(repo ports.TournamentRepository) ports.TournamentRepositoryProvider {
	return &InMemoryTournamentProvider{
		repository: repo,
		mu:         &sync.RWMutex{},
	}
}

func (itp *InMemoryTournamentProvider) WriteTx(do func(ports.TournamentRepository) error) error {
	itp.mu.Lock()
	defer itp.mu.Unlock()

	return do(itp.repository)
}

func (itp *InMemoryTournamentProvider) ReadTx(do func(ports.TournamentRepository) error) error {
	itp.mu.RLock()
	defer itp.mu.RUnlock()

	return do(itp.repository)
}
