package main

import "time"

// --- Event types for the publisher ---

// EventType categorizes what happened.
type EventType string

const (
	EventTournamentDiscovered EventType = "tournament.discovered"
	EventTournamentUpdated    EventType = "tournament.updated"
	EventTournamentCompleted  EventType = "tournament.completed"
	EventStandingsUpdated     EventType = "standings.updated"
	EventScrapeError          EventType = "scrape.error"
)

// Event is published when the scraper discovers or updates tournament data.
type Event struct {
	Type         EventType   `json:"type"`
	Timestamp    time.Time   `json:"timestamp"`
	TournamentID string      `json:"tournament_id"`
	Federation   string      `json:"federation,omitempty"`
	Payload      interface{} `json:"payload,omitempty"`
}
