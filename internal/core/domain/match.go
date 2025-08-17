package domain

import (
	"time"

	"github.com/google/uuid"
)

// Match represents a match between two players
type Match struct {
	ID           int
	TournamentID uuid.UUID
	Rated        bool
	WhitePlayer  WhitePlayer
	BlackPlayer  BlackPlayer
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// WhitePlayer represents a player playing white
type WhitePlayer struct {
	ID        int
	PublicID  uuid.UUID
	FirstName string
	LastName  string
}

// BlackPlayer represents a player playing black
type BlackPlayer struct {
	ID        int
	PublicID  uuid.UUID
	FirstName string
	LastName  string
}
