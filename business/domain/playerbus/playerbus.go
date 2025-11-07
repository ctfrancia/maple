// Package playerbus provides business access to the player domain.
package playerbus

import (
	"errors"

	"github.com/google/uuid"
)

// Set of errors raised by the player domain.
var (
	ErrPlayerNotFound        = errors.New("player not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer declares the behaviour this package needs to persist and retrieve data.
type Storer interface {
	NewWithTx()
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

// NewBusiness creates a new instance of the player business.
func NewBusiness(log any, playerbus any, storer Storer) *Business {
	a := Business{}

	return &a
}

// Create creates a new player.
func (a *Business) Create() {
}

// Update updates an existing player.
func (a *Business) Update() {
}

// Delete deletes an existing player.
func (a *Business) Delete() {
}

// Query returns a list of players.
func (a *Business) Query() {
}
