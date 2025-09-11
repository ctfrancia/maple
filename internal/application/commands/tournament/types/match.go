package types

import "github.com/google/uuid"

// Match represents a match between two players
type Match struct {
	ID           uuid.UUID // this is the public ID
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
