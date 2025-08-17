package domain

import "github.com/google/uuid"

// Player represents a player in a tournament/match (human or computer)
// if they have a public ID, they are a human
// and also represents a "User" When they login
type Player struct {
	IsHuman         bool
	ID              int
	PublicID        uuid.UUID
	Username        string
	Email           string
	Password        string
	FirstName       string
	LastName        string
	Website         string
	ClubAffiliation string
	FIDE            Fide
	Regional        Regional
}

type Fide struct {
	Rating string
	URL    string
	Title  string
}

type Regional struct {
	Country string
	City    string
	Rating  string
	Title   string
}
