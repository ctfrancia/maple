// Package playerbus provides business access to the player domain.
package playerbus

import (
	"context"
	"errors"

	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/google/uuid"
)

// Set of errors raised by the player domain.
var (
	ErrPlayerNotFound        = errors.New("player not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer declares the behaviour this package needs to persist
// and retrieve data.
type Storer interface {
	Create(ctx context.Context, np NewPlayer) error
	Query(ctx context.Context, filter any, orderBy any, page any) ([]Player, error)
	Update(ctx context.Context, id uuid.UUID, np NewPlayer) error
	Delete(ctx context.Context, id uuid.UUID) error
	// Used for scraping
	UpsertPlayers(ctx context.Context, tournamentID uuid.UUID, players []Player) error
}

// Business manages the set of APIs for the player api access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness creates a new instance of the player business.
func NewBusiness(log any, playerbus any, storer Storer) *Business {
	a := Business{}

	return &a
}

// Create creates a new player.
func (a *Business) Create(ctx context.Context, np NewPlayer) error {
	return nil
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
