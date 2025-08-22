package domain

import (
	"time"

	"github.com/google/uuid"
)

// Match represents a match between two players
type Match struct {
	ID           int // private
	UUID         uuid.UUID
	TournamentID Tournament
	Winner       Player
	Location     Location
	City         string
	State        string
	Country      string
	Rated        bool
	WhitePlayer  Player
	BlackPlayer  Player
	PGN          string // Portable Game Notation
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
