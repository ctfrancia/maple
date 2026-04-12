package scraperapp

import (
	"time"

	"github.com/ctfrancia/maple/business/domain/scraperbus"
)

// ScrapeJob represents a scrape job in the API
type ScrapeJob struct {
	ID          string     `json:"id"`
	Federation  string     `json:"federation"`
	Status      string     `json:"status"`
	StartedAt   string     `json:"started_at"`
	CompletedAt *string    `json:"completed_at,omitempty"`
	Error       *string    `json:"error,omitempty"`
	Found       int        `json:"found"`
	Created     int        `json:"created"`
	Updated     int        `json:"updated"`
	DateCreated string     `json:"date_created"`
	DateUpdated string     `json:"date_updated"`
}

// NewScrapeJob represents the data needed to create a new scrape job
type NewScrapeJob struct {
	Federation string `json:"federation"`
}

// ScrapedTournament represents a scraped tournament in the API
type ScrapedTournament struct {
	ID            string  `json:"id"`
	ExternalID    string  `json:"external_id"`
	Name          string  `json:"name"`
	URL           string  `json:"url"`
	Federation    string  `json:"federation"`
	Organizer     *string `json:"organizer,omitempty"`
	Arbiter       *string `json:"arbiter,omitempty"`
	Location      *string `json:"location,omitempty"`
	StartDate     *string `json:"start_date,omitempty"`
	Players       int     `json:"players"`
	Rounds        int     `json:"rounds"`
	TimeControl   *string `json:"time_control,omitempty"`
	LastScrapedAt string  `json:"last_scraped_at"`
	DateCreated   string  `json:"date_created"`
	DateUpdated   string  `json:"date_updated"`
}

// toAppScrapeJob converts a business scrape job to an app scrape job
func toAppScrapeJob(job scraperbus.ScrapeJob) ScrapeJob {
	var completedAt *string
	if job.CompletedAt != nil {
		s := job.CompletedAt.Format(time.RFC3339)
		completedAt = &s
	}

	return ScrapeJob{
		ID:          job.ID.String(),
		Federation:  job.Federation,
		Status:      job.Status,
		StartedAt:   job.StartedAt.Format(time.RFC3339),
		CompletedAt: completedAt,
		Error:       job.Error,
		Found:       job.Found,
		Created:     job.Created,
		Updated:     job.Updated,
		DateCreated: job.DateCreated.Format(time.RFC3339),
		DateUpdated: job.DateUpdated.Format(time.RFC3339),
	}
}

// toAppScrapedTournament converts a business scraped tournament to an app scraped tournament
func toAppScrapedTournament(t scraperbus.ScrapedTournament) ScrapedTournament {
	return ScrapedTournament{
		ID:            t.ID.String(),
		ExternalID:    t.ExternalID,
		Name:          t.Name,
		URL:           t.URL,
		Federation:    t.Federation,
		Organizer:     t.Organizer,
		Arbiter:       t.Arbiter,
		Location:      t.Location,
		StartDate:     t.StartDate,
		Players:       t.Players,
		Rounds:        t.Rounds,
		TimeControl:   t.TimeControl,
		LastScrapedAt: t.LastScrapedAt.Format(time.RFC3339),
		DateCreated:   t.DateCreated.Format(time.RFC3339),
		DateUpdated:   t.DateUpdated.Format(time.RFC3339),
	}
}
