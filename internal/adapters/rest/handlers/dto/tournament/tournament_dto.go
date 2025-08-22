// Package dto is the data transfer object for the REST API
package dto

import "time"

type CreateTournamentRequest struct {
	Name string `json:"name"`
}

type Tournament struct {
	ID                 string        `json:"id"` // public uuid
	Name               string        `json:"name"`
	Location           string        `json:"location"` // address or location
	Description        string        `json:"description"`
	OpenToPublic       bool          `json:"open_to_public"`
	OpenToSpectators   bool          `json:"open_to_spectators"`
	OpenToRegistration bool          `json:"open_to_registration"`
	Registration       *Registration `json:"registration,omitempty"`
	Arbitrator         string        `json:"arbitrator"` // name of the person
	PairingMethod      string        `json:"pairing_method,omitempty"`
	Matches            *[]Match      `json:"matches,omitempty"`
	Players            []string      `json:"players,omitempty"`           // this will be there public IDS
	NumberOfPlayers    int           `json:"number_of_players,omitempty"` // how many are participating
	Schedule           *[]Schedule   `json:"schedule,omitempty"`
	Results            *[]Result     `json:"results,omitempty"`
	Status             string        `json:"status"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	SoftDeletedAt      time.Time     `json:"soft_deleted_at"`
}

type Registration struct {
	Status    string
	StartTime time.Time
	EndTime   time.Time
	Fee       int64
	PrizePool int64
	Payout    Payout
}

type Match struct {
	UUID          string    `json:"id"`            // public uuid
	TournamentID  string    `json:"tournament_id"` // public uuid
	Winner        string    `json:"winner"`        // public uuid
	Location      string    `json:"location"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	Rated         bool      `json:"rated"`
	WhitePlayer   string    `json:"white_player"` // public uuid
	BlackPlayer   string    `json:"black_player"` // public uuid
	PGN           string    `json:"pgn"`          // Portable Game Notation
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	SoftDeletedAt time.Time `json:"soft_deleted_at"`
}

type Schedule struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type Result struct {
	Player string `json:"name"`
	Prize  int64  `json:"prize"`
}

type Payout struct {
	Place  int8   `json:"place"` // first,second, etc.
	Amount int64  `json:"amount"`
	Other  string `json:"other"` // maybe they get a book or a subscription
}
