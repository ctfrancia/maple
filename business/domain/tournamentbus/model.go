package tournamentbus

import (
	"net/mail"
	"net/url"
	"time"

	"github.com/ctfrancia/maple/business/types/address"
	"github.com/ctfrancia/maple/business/types/city"
	"github.com/ctfrancia/maple/business/types/country"
	"github.com/ctfrancia/maple/business/types/description"
	"github.com/ctfrancia/maple/business/types/name"
	"github.com/ctfrancia/maple/business/types/phone"
	"github.com/ctfrancia/maple/business/types/poster"
	"github.com/ctfrancia/maple/business/types/round"
	"github.com/ctfrancia/maple/business/types/state"
	"github.com/ctfrancia/maple/business/types/zip"

	"github.com/google/uuid"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusClosed    Status = "closed"
)

// Tournament represents information about an individual tournament.
type Tournament struct {
	ID           uuid.UUID
	Name         name.Name
	CreatedBy    uuid.UUID // the website owner (api consumer ID)
	PlayerID     uuid.UUID // the player ID (user in the website's system)
	Enabled      bool
	Description  description.Description
	Poster       poster.Poster
	Location     Location
	Schedule     []Schedule
	Contact      Contact
	Registration Registration
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewTournament represents the information we need to create a new tournament.
type NewTournament struct {
	Name        name.Name
	CreatedBy   uuid.UUID // the website owner (api consumer ID)
	PlayerID    uuid.UUID // the player ID (user in the website's system)
	Description description.Description
}

// UpdateTournament contains information needed to update a user.
type UpdateTournament struct {
	Name        *name.Name
	Description *description.Description
	Poster      *poster.Poster
}

// RequiredTournament represents the minimal information needed
// to publish a tournament on the platform.
type RequiredTournament struct {
}

// Schedule represents a schedule for a tournament.
type Schedule struct {
	ID      uuid.UUID
	Rounds  round.Round
	Dates   []time.Time
	Matches []uuid.UUID
}

// Contact represents the contact information for a tournament.
type Contact struct {
	Name  name.Name
	Email mail.Address
	Phone phone.Phone
	URL   url.URL
}

// Location represents the location of a tournament.
type Location struct {
	ID       uuid.UUID
	Name     name.Name
	Address  address.Address
	Address2 address.Address
	City     city.City
	State    state.State
	Zip      zip.Zip
	Country  country.Country
}

// Registration represents the registration information for a tournament.
type Registration struct {
	Status  string
	Open    time.Time
	Close   time.Time
	Fee     uint32
	Contact Contact
}
