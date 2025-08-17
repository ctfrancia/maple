package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrTournamentNotFound = errors.New("tournament not found")

// TournamentStatus - tournament states
type TournamentStatus string

const (
	TournamentStatusActive    TournamentStatus = "active"
	TournamentStatusInactive  TournamentStatus = "inactive"
	TournamentStatusSuspended TournamentStatus = "suspended"
	TournamentStatusPending   TournamentStatus = "pending"
	TournamentStatusCompleted TournamentStatus = "completed"
)

type Tournament struct {
	ID                 int       // private
	PublicID           uuid.UUID // this is the public ID
	Name               string
	Location           string
	Description        string
	OpenToPublic       bool
	OpenToSpectators   bool
	OpenToRegistration bool
	RegistrationOpen   bool
	NumberOfPlayers    int
	Schedule           []Schedule
	Status             TournamentStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Schedule struct {
	StartTime time.Time
	EndTime   time.Time
}
