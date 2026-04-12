package scraperbus

import (
	"time"

	"github.com/google/uuid"
)

// ScrapeJob represents a scraping job
type ScrapeJob struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Federation  string    `gorm:"type:varchar(10);not null;index"`
	Status      string    `gorm:"type:varchar(20);not null;index"`
	StartedAt   time.Time `gorm:"not null"`
	CompletedAt *time.Time
	Error       *string
	Found       int `gorm:"default:0"`
	Created     int `gorm:"default:0"`
	Updated     int `gorm:"default:0"`
	DateCreated time.Time `gorm:"not null;index"`
	DateUpdated time.Time `gorm:"not null"`
}

// ScrapedTournament represents tournament data scraped from chess-results.com
type ScrapedTournament struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	ExternalID    string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Name          string    `gorm:"type:varchar(255);not null"`
	URL           string    `gorm:"type:varchar(500)"`
	Federation    string    `gorm:"type:varchar(10);index"`
	Organizer     *string   `gorm:"type:varchar(255)"`
	Arbiter       *string   `gorm:"type:varchar(255)"`
	Location      *string   `gorm:"type:varchar(255)"`
	StartDate     *string   `gorm:"type:varchar(50)"`
	Players       int       `gorm:"default:0"`
	Rounds        int       `gorm:"default:0"`
	TimeControl   *string   `gorm:"type:varchar(255)"`
	LastScrapedAt time.Time `gorm:"not null;index"`
	DateCreated   time.Time `gorm:"not null;index"`
	DateUpdated   time.Time `gorm:"not null"`
}

// NewScrapeJob represents the data needed to create a new scrape job
type NewScrapeJob struct {
	Federation string
}

// UpdateScrapeJob represents the data that can be updated for a scrape job
type UpdateScrapeJob struct {
	Status      *string
	CompletedAt *time.Time
	Error       *string
	Found       *int
	Created     *int
	Updated     *int
}

// ScrapeJobStatus constants
const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

// TableName overrides the table name for ScrapeJob
func (ScrapeJob) TableName() string {
	return "scrape_jobs"
}

// TableName overrides the table name for ScrapedTournament
func (ScrapedTournament) TableName() string {
	return "scraped_tournaments"
}
