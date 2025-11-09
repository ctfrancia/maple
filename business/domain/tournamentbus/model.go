package tournamentbus

import (
	"github.com/ctfrancia/maple/business/types/name"
	"github.com/ctfrancia/maple/business/types/pgn"
	"github.com/google/uuid"
	"time"
)

type Schedule struct {
	Date    time.Time
	Winner  uuid.UUID
	Players [2]uuid.UUID // 0: white, 1: black
	PGN     pgn.PGN
}

// NewTournament represents the information we need to create a new tournament.
type NewTournament struct {
	Name *string
}

// Tournament represents information about an individual tournament.
type Tournament struct {
	ID          uuid.UUID
	Name        name.Name
	Description description.Description // types/description/description.go
	Poster      poster.Poster           // types/poster/poster.go
	MatchMaking matchmaking.MatchMaking // types/matchmaking/matchmaking.go
	Location    location.Location       // types/location/location.go
	Schedule    []Schedule              // types/schedule/schedule.go ??
}

// UpdateTournament contains information needed to update a user.
type UpdateTournament struct {
}
