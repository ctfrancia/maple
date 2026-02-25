package tournamentcache

import (
	"context"
	"sync"
	"time"

	"github.com/ctfrancia/maple/business/domain/tournamentbus/stores/tournamentdb"
)

// MemoryStore is a simple in-memory store for development and tests.
type MemoryStore struct {
	mu          sync.RWMutex
	tournaments map[string]tournamentdb.Tournament
	//players     map[string][]models.Player
}

// NewMemoryStore returns a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tournaments: make(map[string]tournamentdb.Tournament),
		//players:     make(map[string][]models.Player),
	}
}

func (m *MemoryStore) KnownTournamentIDs(_ context.Context, federation string) (map[string]bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ids := make(map[string]bool)
	for id, t := range m.tournaments {
		if federation == "" || t.Federation == federation {
			ids[id] = true
		}
	}
	return ids, nil
}

func (m *MemoryStore) UpsertTournament(_ context.Context, t tournamentdb.Tournament) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.tournaments[t.PublicID.String()]
	now := time.Now()
	if !exists {
		t.CreatedAt = now
		t.Status = "discovered"
	}
	t.UpdatedAt = now
	m.tournaments[t.PublicID.String()] = t
	return !exists, nil
}

func (m *MemoryStore) UpdateStatus(_ context.Context, id, status string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t, ok := m.tournaments[id]; ok {
		t.Status = status
		t.UpdatedAt = time.Now()
	}

	return nil
}

func (m *MemoryStore) TournamentsToScrape(_ context.Context, maxAge time.Duration) ([]tournamentdb.Tournament, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []tournamentdb.Tournament
	cutoff := time.Now().Add(-maxAge)
	for _, t := range m.tournaments {
		if t.Status == "discovered" || t.LastScrapedAt.Before(cutoff) {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *MemoryStore) Close() error { return nil }
