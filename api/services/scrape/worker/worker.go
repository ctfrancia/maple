package worker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ctfrancia/maple/api/services/scrape/client"
	"github.com/ctfrancia/maple/api/services/scrape/events"
	"github.com/ctfrancia/maple/api/services/scrape/parser"
	"github.com/ctfrancia/maple/api/services/scrape/publisher"
	"github.com/ctfrancia/maple/business/domain/scraperbus"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/google/uuid"
)

// Worker processes scrape jobs
type Worker struct {
	log            *logger.Logger
	scraperBus     scraperbus.ExtBusiness
	client         *client.Client
	publisher      publisher.Publisher
	locationFilter string
}

// Config holds the worker configuration
type Config struct {
	Log            *logger.Logger
	ScraperBus     scraperbus.ExtBusiness
	Client         *client.Client
	Publisher      publisher.Publisher
	LocationFilter string
}

// New creates a new worker
func New(cfg Config) *Worker {
	return &Worker{
		log:            cfg.Log,
		scraperBus:     cfg.ScraperBus,
		client:         cfg.Client,
		publisher:      cfg.Publisher,
		locationFilter: cfg.LocationFilter,
	}
}

// ProcessJob processes a single scrape job
func (w *Worker) ProcessJob(ctx context.Context, jobID uuid.UUID) error {
	w.log.Info(ctx, "processing job", "job_id", jobID.String())

	// Get the job
	job, err := w.scraperBus.QueryJobByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("getting job: %w", err)
	}

	// Update job status to running
	now := time.Now()
	runningStatus := scraperbus.StatusRunning
	if err := w.scraperBus.UpdateJob(ctx, job, scraperbus.UpdateScrapeJob{
		Status: &runningStatus,
	}); err != nil {
		return fmt.Errorf("updating job status: %w", err)
	}

	// Scrape the federation
	stats, err := w.scrapeFederation(ctx, job.Federation)
	if err != nil {
		// Mark job as failed
		failedStatus := scraperbus.StatusFailed
		completedAt := time.Now()
		errMsg := err.Error()
		updateErr := w.scraperBus.UpdateJob(ctx, job, scraperbus.UpdateScrapeJob{
			Status:      &failedStatus,
			CompletedAt: &completedAt,
			Error:       &errMsg,
		})
		if updateErr != nil {
			w.log.Error(ctx, "failed to update job after error", "error", updateErr)
		}

		// Publish error event
		w.publisher.Publish(ctx, events.Event{
			Type:       events.ScrapeError,
			Timestamp:  time.Now(),
			Federation: job.Federation,
			Payload: map[string]interface{}{
				"job_id": jobID.String(),
				"error":  err.Error(),
			},
		})

		return fmt.Errorf("scraping federation: %w", err)
	}

	// Update job as completed
	completedStatus := scraperbus.StatusCompleted
	completedAt := time.Now()
	if err := w.scraperBus.UpdateJob(ctx, job, scraperbus.UpdateScrapeJob{
		Status:      &completedStatus,
		CompletedAt: &completedAt,
		Found:       &stats.Found,
		Created:     &stats.Created,
		Updated:     &stats.Updated,
	}); err != nil {
		return fmt.Errorf("updating job completion: %w", err)
	}

	w.log.Info(ctx, "job completed",
		"job_id", jobID.String(),
		"found", stats.Found,
		"created", stats.Created,
		"updated", stats.Updated,
		"duration", time.Since(now).String(),
	)

	return nil
}

// ScrapeStats holds statistics about a scrape operation
type ScrapeStats struct {
	Found   int
	Created int
	Updated int
}

// scrapeFederation scrapes all tournaments for a given federation
func (w *Worker) scrapeFederation(ctx context.Context, federation string) (*ScrapeStats, error) {
	w.log.Info(ctx, "scraping federation", "federation", federation)

	stats := &ScrapeStats{}

	// Fetch the federation page
	federationHTML, err := w.client.Get(ctx, client.FederationURL(federation))
	if err != nil {
		return nil, fmt.Errorf("fetching federation page: %w", err)
	}

	// Parse tournament listings
	entries, err := parser.ParseFederationTournaments(federationHTML)
	if err != nil {
		return nil, fmt.Errorf("parsing federation tournaments: %w", err)
	}

	w.log.Info(ctx, "found tournaments in federation", "count", len(entries), "federation", federation)

	// Filter by location if configured
	var filteredEntries []parser.FederationTournamentEntry
	if w.locationFilter != "" {
		w.log.Info(ctx, "applying location filter", "filter", w.locationFilter)
		for _, entry := range entries {
			if matchesLocation(entry.Location, w.locationFilter) || matchesLocation(entry.Name, w.locationFilter) {
				filteredEntries = append(filteredEntries, entry)
			}
		}
		w.log.Info(ctx, "tournaments after location filter",
			"before", len(entries),
			"after", len(filteredEntries),
			"filter", w.locationFilter,
		)
	} else {
		filteredEntries = entries
	}

	stats.Found = len(filteredEntries)
	w.log.Info(ctx, "processing tournaments", "count", stats.Found, "federation", federation)

	// Process each tournament
	for i, entry := range filteredEntries {
		if ctx.Err() != nil {
			return stats, ctx.Err()
		}

		w.log.Info(ctx, "processing tournament",
			"index", i+1,
			"total", len(entries),
			"external_id", entry.ID,
			"name", entry.Name,
		)

		// Fetch tournament details
		tournamentHTML, err := w.client.Get(ctx, client.TournamentURL(entry.ID, 0))
		if err != nil {
			w.log.Error(ctx, "failed to fetch tournament", "external_id", entry.ID, "error", err)
			continue
		}

		// Parse tournament info
		tournamentInfo, err := parser.ParseTournamentInfo(tournamentHTML, entry.ID)
		if err != nil {
			w.log.Error(ctx, "failed to parse tournament", "external_id", entry.ID, "error", err)
			continue
		}

		// Check if tournament already exists (we'll determine this in the database layer via upsert)
		// For now, we'll just track stats based on the upsert result
		isNew := true // This will be refined when we add proper checking

		// Upsert tournament
		now := time.Now()
		scrapedTournament := scraperbus.ScrapedTournament{
			ID:            uuid.New(),
			ExternalID:    entry.ID,
			Name:          tournamentInfo.Name,
			URL:           tournamentInfo.URL,
			Federation:    tournamentInfo.Federation,
			Organizer:     stringPtr(tournamentInfo.Organizer),
			Arbiter:       stringPtr(tournamentInfo.Arbiter),
			Location:      stringPtr(tournamentInfo.Location),
			StartDate:     stringPtr(tournamentInfo.StartDate),
			Players:       tournamentInfo.Players,
			Rounds:        tournamentInfo.Rounds,
			TimeControl:   stringPtr(tournamentInfo.TimeControl),
			LastScrapedAt: now,
			DateCreated:   now,
			DateUpdated:   now,
		}

		if err := w.scraperBus.UpsertTournament(ctx, scrapedTournament); err != nil {
			w.log.Error(ctx, "failed to upsert tournament", "external_id", entry.ID, "error", err)
			continue
		}

		// Update stats
		if isNew {
			stats.Created++
			// Publish discovery event
			w.publisher.Publish(ctx, events.Event{
				Type:         events.TournamentDiscovered,
				Timestamp:    time.Now(),
				TournamentID: entry.ID,
				Federation:   federation,
				Payload:      scrapedTournament,
			})
		} else {
			stats.Updated++
			// Publish update event
			w.publisher.Publish(ctx, events.Event{
				Type:         events.TournamentUpdated,
				Timestamp:    time.Now(),
				TournamentID: entry.ID,
				Federation:   federation,
				Payload:      scrapedTournament,
			})
		}
	}

	return stats, nil
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// matchesLocation checks if a location string matches the filter
func matchesLocation(location, filter string) bool {
	if location == "" {
		return false
	}
	// Case-insensitive substring match
	locationLower := strings.ToLower(location)
	filterLower := strings.ToLower(filter)
	return strings.Contains(locationLower, filterLower)
}
