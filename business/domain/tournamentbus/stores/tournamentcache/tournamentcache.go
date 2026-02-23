package tournamentcache

import (
	"sync"

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
