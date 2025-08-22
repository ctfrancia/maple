package domain

import "github.com/google/uuid"

type Timezone string

// TODO: add eu, us, etc
const (
	TimezoneUTC Timezone = "UTC"
	TimezonePST Timezone = "PST"
	TimezoneEST Timezone = "EST"
)

// Location represents a location in the world
type Location struct {
	ID         int
	PublicID   uuid.UUID
	ClubAffil  *Club
	Name       string
	Address    string
	PostalCode string
	City       string // have a city db
	State      string // have a state db
	County     string // have a county db
	Province   string // have a province db
	Country    string // have a country db
	Latitude   float64
	Longitude  float64
	Timezone   Timezone
}
