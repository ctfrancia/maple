// Package playerbus provides business access to the player domain.
package playerbus

import (
	"errors"
)

// Set of errors raised by the player domain.
var (
	ErrPlayerNotFound = errors.New("player not found")
)

type Storer interface {
	NewWithTx(tx any) (*Player, error)
	SavePlayer(player *Player) error
}

// Business manages the set of APIs for the player api access.
type Business struct {
	log       any
	playerbus any
}

// NewBusiness constructs a player business API for use.
func NewBusiness(log any, playerbus any) *Business {
}
