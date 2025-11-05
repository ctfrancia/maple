// Package playerbus provides business access to the player domain.
package playerbus

import (
	"errors"

	"github.com/google/uuid"
)

// Set of errors raised by the player domain.
var (
	ErrPlayerNotFound = errors.New("player not found")
)

// Storer declares the behaviour this package needs to persist and retrieve data.
type Storer interface {
	NewWithTx()
	SavePlayer()
	Create()
	Update()
	Delete()
	Query()
}

// Business manages the set of APIs for the player api access.
type Business struct {
	log       any
	playerbus any
	storer    Storer
}

// NewBusiness constructs a player business API for use.
func NewBusiness(log any, playerbus any, storer Storer) *Business {
	b := Business{}

	return &b
}

// Create adds a new chess player to the system.
func (b *Business) Create() {
}

// Update modifies an existing chess player in the system.
func (b *Business) Update() {
}

// Delete removes an existing chess player from the system.
func (b *Business) Delete() {
}

// Query retrieves a list of existing chess players from the system.
func (b *Business) Query() {
}

// Count returns the total number of players in the system.
func (b *Business) Count() {
}

// QueryByID retrieves a single chess player by ID.
func (b *Business) QueryByID(ID uuid.UUID) {
}
