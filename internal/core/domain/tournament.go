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

type RegistrationStatus string

const (
	RegistrationStatusOpen   RegistrationStatus = "open"
	RegistrationStatusClosed RegistrationStatus = "closed"
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
	Arbitrator         string
	PairingMethod      PairingMethod
	Matches            []Match
	Players            []Player // no more no less than 2 white/black
	NumberOfPlayers    int      // how many are participating
	Schedule           []Schedule
	Results            []Result
	Status             TournamentStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
	SoftDeletedAt      time.Time
	DeletedAt          time.Time
}

type Schedule struct {
	StartTime time.Time
	EndTime   time.Time
}

type Registration struct {
	Status    RegistrationStatus
	StartTime time.Time
	EndTime   time.Time
	Fee       int64
	PrizePool int64
	Payment   Payment
}

type Payment struct {
	Currency      string
	AmountFirst   int64
	AmountSecond  int64
	AmountThird   int64
	AmountFourth  int64
	AmountFifth   int64
	AmountSixth   int64
	AmountSeventh int64
	AmountEighth  int64
	AmountNinth   int64
	AmountTenth   int64
}

type Result struct {
	PlayerID Player
	Prize    int64
}
