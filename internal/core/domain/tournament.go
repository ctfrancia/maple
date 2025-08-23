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
	TournamentStatusDraft     TournamentStatus = "draft"
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
	Location           Location // address or location
	Creator            Player
	Contact            Contact
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
	Payment   []Payment
}

type PaymentType string

const (
	PaymentTypeMonetary PaymentType = "monetary" //money
	PaymentTypePhysical PaymentType = "physical" // e.g. book/lesson/etc.
	PaymentTypeOther    PaymentType = "other"
)

type Payment struct {
	Place  int   // 1st, 2nd, etc.
	Amount int64 // if type is monetary
	Type   PaymentType
}

type Result struct {
	PlayerID Player
	Prize    int64
}

type Contact struct {
	Name  string
	Email string
	Phone string
}
