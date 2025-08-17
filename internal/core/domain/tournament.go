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

// PairingMethod - the method used to pair players
type PairingMethod string

const (
	PairingMethodNone       PairingMethod = "none"
	PairingMethodDraw       PairingMethod = "draw"
	PairingMethodRoundRobin PairingMethod = "round_robin"
)

type Tournament struct {
	ID                 int       // private
	PublicID           uuid.UUID // this is the public ID
	Name               string
	Location           string // address or location
	Description        string
	OpenToPublic       bool
	OpenToSpectators   bool
	OpenToRegistration bool
	Registration       Registration
	PairingMethod      PairingMethod
	Matches            []Match
	Players            []Player
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

type Registration struct {
	StartTime time.Time
	EndTime   time.Time
	Fee       float64
	PrizePool float64
	Payment   Payment
}

type Payment struct {
	Currency      string
	AmountFirst   float64
	AmountSecond  float64
	AmountThird   float64
	AmountFourth  float64
	AmountFifth   float64
	AmountSixth   float64
	AmountSeventh float64
	AmountEighth  float64
	AmountNinth   float64
	AmountTenth   float64
}
