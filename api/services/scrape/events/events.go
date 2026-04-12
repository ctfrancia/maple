package events

import "time"

// EventType categorizes what happened.
type EventType string

const (
	TournamentDiscovered EventType = "tournament.discovered"
	TournamentUpdated    EventType = "tournament.updated"
	TournamentCompleted  EventType = "tournament.completed"
	StandingsUpdated     EventType = "standings.updated"
	ScrapeError          EventType = "scrape.error"
)

// Event is published when the scraper discovers or updates tournament data.
type Event struct {
	Type         EventType `json:"type"`
	Timestamp    time.Time `json:"timestamp"`
	TournamentID string    `json:"tournament_id"`
	Federation   string    `json:"federation,omitempty"`
	Payload      any       `json:"payload,omitempty"`
}
