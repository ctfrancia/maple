// Package dto is the data transfer object for the REST API
package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTournamentRequest struct {
	Name string `json:"name"`
}

type TournamentStatus string

const (
	TournamentStatusActive    TournamentStatus = "active"
	TournamentStatusDraft     TournamentStatus = "draft"
	TournamentStatusInactive  TournamentStatus = "inactive"
	TournamentStatusSuspended TournamentStatus = "suspended"
	TournamentStatusPending   TournamentStatus = "pending"
	TournamentStatusCompleted TournamentStatus = "completed"
)

type Tournament struct {
	ID                 string           `json:"id"` // public uuid
	Name               string           `json:"name"`
	Location           Location         `json:"location"`
	Description        string           `json:"description"`
	OpenToPublic       bool             `json:"open_to_public"`
	OpenToSpectators   bool             `json:"open_to_spectators"`
	OpenToRegistration bool             `json:"open_to_registration"`
	Registration       Registration     `json:"registration"`
	Arbitrator         string           `json:"arbitrator"` // name of the person
	PairingMethod      string           `json:"pairing_method"`
	Matches            []Match          `json:"matches,omitempty"`
	Players            []string         `json:"players,omitempty"` // this will be there public IDS
	NumberOfPlayers    int              `json:"number_of_players"` // how many are participating
	Schedule           []Schedule       `json:"schedule,omitempty"`
	Results            []Result         `json:"results"`
	Status             TournamentStatus `json:"status"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	SoftDeletedAt      *time.Time       `json:"soft_deleted_at,omitempty"` // omit if not soft deleted
}

type Location struct {
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

type RegistrationStatus string

type Registration struct {
	Status    RegistrationStatus `json:"status"`
	StartTime time.Time          `json:"start_time"`
	EndTime   time.Time          `json:"end_time"`
	Fee       int64              `json:"fee"`
	PrizePool int64              `json:"prize_pool"`
	Payment   []Payment          `json:"payout"`
}

type Match struct {
	UUID          string     `json:"id"`            // public uuid
	TournamentID  string     `json:"tournament_id"` // public uuid
	Winner        string     `json:"winner"`        // public uuid
	Location      string     `json:"location"`
	City          string     `json:"city"`
	State         string     `json:"state"`
	Country       string     `json:"country"`
	Rated         bool       `json:"rated"`
	WhitePlayer   string     `json:"white_player"` // public uuid
	BlackPlayer   string     `json:"black_player"` // public uuid
	PGN           string     `json:"pgn"`          // Portable Game Notation
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	SoftDeletedAt *time.Time `json:"soft_deleted_at,omitempty"` // omit if not soft deleted
}

type Schedule struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type Result struct {
	Player string `json:"name"`
	Prize  int64  `json:"prize"`
}

type Payment struct {
	Place  int8   `json:"place"` // first,second, etc.
	Amount int64  `json:"amount"`
	Other  string `json:"other"` // maybe they get a book or a subscription
}
