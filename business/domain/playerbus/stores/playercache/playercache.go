package playercache

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MemoryStore is a simple in-memory store for development and tests.
type MemoryStore struct {
	mu      sync.RWMutex
	players map[uuid.UUID][]models.Player
}

// NewMemoryStore returns a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		players: make(map[string][]models.Player),
	}
}

func (m *MemoryStore) UpsertPlayers(_ context.Context, tournamentID string, players []models.Player) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.players[tournamentID] = players
	if t, ok := m.tournaments[tournamentID]; ok {
		t.LastScraped = time.Now()
		t.Status = "scraped"
	}
	return nil
}
