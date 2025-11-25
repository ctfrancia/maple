package tournamentbus

import (
	"net/mail"
	"net/url"
	"time"

	"github.com/ctfrancia/maple/business/types/address"
	"github.com/ctfrancia/maple/business/types/arbiter"
	"github.com/ctfrancia/maple/business/types/category"
	"github.com/ctfrancia/maple/business/types/city"
	"github.com/ctfrancia/maple/business/types/country"
	"github.com/ctfrancia/maple/business/types/description"
	"github.com/ctfrancia/maple/business/types/federation"
	"github.com/ctfrancia/maple/business/types/fide"
	"github.com/ctfrancia/maple/business/types/location"
	"github.com/ctfrancia/maple/business/types/name"
	"github.com/ctfrancia/maple/business/types/pgn"
	"github.com/ctfrancia/maple/business/types/phone"
	"github.com/ctfrancia/maple/business/types/poster"
	"github.com/ctfrancia/maple/business/types/qrcode"
	"github.com/ctfrancia/maple/business/types/round"
	"github.com/ctfrancia/maple/business/types/state"
	"github.com/ctfrancia/maple/business/types/timecontrol"
	"github.com/ctfrancia/maple/business/types/zip"

	"github.com/google/uuid"
)

type Status int

const (
	StatusClosed Status = iota
	StatusDraft
	StatusPublished
)

type TournamentType int

const (
	TournamentTypeSwissSystem TournamentType = iota
	TournamentTypeRoundRobin
	TournamentTypeTeamSwissSystem
	TournamentTypeTeamRoundRobin
)

type RatingCalculation int

const (
	RatingCalculationNone RatingCalculation = iota
	RatingCalculationInternational
)

type Tournament struct {
	ID                  uuid.UUID
	OrganizerID         uuid.UUID               // this will be the chess club
	Name                name.Name               // name of the tournament
	Description         description.Description // description of the tournament
	Poster              poster.Poster           // poster of the tournament
	FIDE                fide.Fide               // any fide specific info
	Federation          federation.Federation   // any federation specific info
	QRCode              qrcode.QRCode           // linking to /tournament/{ID}
	Categories          []category.Category
	Links               []url.URL
	Director            arbiter.Arbiter
	ChiefArbiter        arbiter.Arbiter
	DeputryChiefArbiter arbiter.Arbiter
	Arbiters            []arbiter.Arbiter
	TimeControl         timecontrol.TimeControl
	Location            location.Location
	StartDate           time.Time
	EndDate             time.Time // if no end date, then it's only a one-day tournament
	RatingCalculation   RatingCalculation
	Schedule            []Schedule
	Rounds              round.Round // can get by the schedule
	Prizes              []Prize
	Type                TournamentType
	PairingProgram      []url.URL
	Contact             Contact
	CreatedBy           uuid.UUID
	UpdatedBy           uuid.UUID
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Status              Status
	Version             uint16
}

type Prize struct {
	Name  name.Name // name of the prize ("1", "2", "3", etc.)
	Type  string    // type of the prize ("monitary", "physical", etc.)
	Value string    // value of the prize ("100", "1000", "name of book", "lesson with..." etc.)
}

// NewTournament represents the initial information needed to create and build a new tournament.
type NewTournament struct {
	Name        name.Name
	CreatedBy   uuid.UUID // the website owner (api consumer ID)
	PlayerID    uuid.UUID // the player ID (user in the website's system)
	Description description.Description
}

// UpdateTournament contains information needed to update a user.
// all fields are optional.
type UpdateTournament struct {
	ID          uuid.UUID
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
	ID           uuid.UUID
	TournamentID uuid.UUID
	Dates        []time.Time // []2024-12-02 15:00, 2024-12-03 15:00
	Pairings     []Pairing
}

type Pairing struct {
	ID           uuid.UUID
	RoundID      uuid.UUID
	TournamentID uuid.UUID
	Date         time.Time
	Players      [2]uuid.UUID // 0 = white, 1 = black
	PGN          pgn.PGN
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
