package tournamentdb

import (
	"time"

	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

// dsn is the postgres database connection string for development purposes.
var dsn = "host=localhost user=devuser password=devpass dbname=mydb port=5432 sslmode=disable"

// Store manages the set of APIs for the tournament database access.
type Store struct {
	log *logger.Logger
	db  *gorm.DB
}

type tournament struct {
	gorm.Model
	Name         string
	Description  string
	Poster       string
	Rounds       uint8
	Enabled      bool         `gorm:"default:false"`
	Contact      contact      `gorm:"embedded;embeddedPrefix:contact_"`
	Location     location     `gorm:"embedded;embeddedPrefix:location_"`
	Registration registration `gorm:"embedded;embeddedPrefix:registration_"`
	Status       string       `gorm:"type:varchar(10)"`
	// PublicID is the ID of the tournament in the public website.
	PublicID uuid.UUID `gorm:"type:uuid"`
	// HostWebsiteID is the ID of the website that created the tournament.
	HostWebsiteID uuid.UUID `gorm:"type:uuid"`
	// PlayerID is the ID of the player that created the tournament.
	PlayerID uuid.UUID  `gorm:"type:uuid"`
	Schedule []schedule `gorm:"foreignKey:TournamentID"`
}

type schedule struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey"`
	TournamentID uint        `gorm:"not null;index"` // Foreign key to Tournament
	Round        int         `gorm:"not null"`
	ScheduledAt  time.Time   `gorm:"not null;index"` // Index for date queries!
	Matches      []uuid.UUID `gorm:"type:uuid[]"`    // PostgreSQL array
}

type contact struct {
	Name  string `gorm:"type:varchar(50)"`
	Email string `gorm:"type:varchar(50)"`
	Phone string `gorm:"type:varchar(20)"`
	URL   string `gorm:"type:varchar(50)"`
}

type location struct {
	ID       uuid.UUID
	Name     string
	Address  string
	Address2 string
	City     string
	State    string
	Zip      string
	Country  string
}

// Registration represents the registration information for a tournament.
type registration struct {
	Status  string
	Open    string
	Close   string
	Fee     string
	Contact string
}
