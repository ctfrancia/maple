package worker

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/ctfrancia/maple/api/services/scrape/client"
	"github.com/ctfrancia/maple/api/services/scrape/events"
	"github.com/ctfrancia/maple/api/services/scrape/publisher"
	"github.com/ctfrancia/maple/business/domain/scraperbus"
	"github.com/ctfrancia/maple/business/sdk/order"
	"github.com/ctfrancia/maple/business/sdk/page"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/google/uuid"
)

// -------------------------------------------------------------------------
// Fakes

// fakeBus implements scraperbus.ExtBusiness using in-memory state.
type fakeBus struct {
	mu          sync.Mutex
	jobs        map[uuid.UUID]scraperbus.ScrapeJob
	tournaments []scraperbus.ScrapedTournament
	updates     []scraperbus.UpdateScrapeJob // records every UpdateJob call in order
}

func newFakeBus() *fakeBus {
	return &fakeBus{jobs: make(map[uuid.UUID]scraperbus.ScrapeJob)}
}

func (f *fakeBus) addJob(job scraperbus.ScrapeJob) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.jobs[job.ID] = job
}

func (f *fakeBus) CreateJob(_ context.Context, nj scraperbus.NewScrapeJob) (scraperbus.ScrapeJob, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	job := scraperbus.ScrapeJob{
		ID:          uuid.New(),
		Federation:  nj.Federation,
		Status:      scraperbus.StatusPending,
		StartedAt:   time.Now(),
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}
	f.jobs[job.ID] = job
	return job, nil
}

func (f *fakeBus) UpdateJob(_ context.Context, job scraperbus.ScrapeJob, update scraperbus.UpdateScrapeJob) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	existing, ok := f.jobs[job.ID]
	if !ok {
		return errors.New("job not found")
	}
	if update.Status != nil {
		existing.Status = *update.Status
	}
	if update.CompletedAt != nil {
		existing.CompletedAt = update.CompletedAt
	}
	if update.Error != nil {
		existing.Error = update.Error
	}
	if update.Found != nil {
		existing.Found = *update.Found
	}
	if update.Created != nil {
		existing.Created = *update.Created
	}
	if update.Updated != nil {
		existing.Updated = *update.Updated
	}
	f.jobs[job.ID] = existing
	f.updates = append(f.updates, update)
	return nil
}

func (f *fakeBus) QueryJobs(_ context.Context, filter scraperbus.QueryFilter, _ order.By, _ page.Page) ([]scraperbus.ScrapeJob, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var result []scraperbus.ScrapeJob
	for _, j := range f.jobs {
		if filter.Status != nil && j.Status != *filter.Status {
			continue
		}
		result = append(result, j)
	}
	return result, nil
}

func (f *fakeBus) QueryJobByID(_ context.Context, jobID uuid.UUID) (scraperbus.ScrapeJob, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	j, ok := f.jobs[jobID]
	if !ok {
		return scraperbus.ScrapeJob{}, errors.New("job not found")
	}
	return j, nil
}

func (f *fakeBus) UpsertTournament(_ context.Context, tournament scraperbus.ScrapedTournament) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.tournaments = append(f.tournaments, tournament)
	return nil
}

func (f *fakeBus) QueryTournaments(_ context.Context, _ scraperbus.QueryFilter, _ order.By, _ page.Page) ([]scraperbus.ScrapedTournament, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.tournaments, nil
}

// fakePublisher records published events.
type fakePublisher struct {
	mu     sync.Mutex
	events []events.Event
}

func (p *fakePublisher) Publish(_ context.Context, e events.Event) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, e)
	return nil
}

func (p *fakePublisher) Close() error { return nil }

var _ publisher.Publisher = (*fakePublisher)(nil)

// -------------------------------------------------------------------------
// HTML fixtures

const federationHTML = `<html><body>
<table class="CRs1">
<tr><td>Name</td><td>Location</td><td>Date</td></tr>
<tr>
  <td><a href="tnr111111.aspx?lan=1">Copa Catalana 2024</a></td>
  <td>Barcelona, Spain</td><td>15/02/2024</td>
</tr>
<tr>
  <td><a href="tnr222222.aspx?lan=1">Open Girona 2024</a></td>
  <td>Girona, Spain</td><td>01/03/2024</td>
</tr>
</table>
</body></html>`

const tournamentHTML = `<html>
<head><title>Chess-Results Server - Copa Catalana 2024</title></head>
<body>
<table class="CRs1">
<tr><td>Tournament name</td><td>Copa Catalana 2024</td></tr>
<tr><td>Federation</td><td>CAT</td></tr>
<tr><td>Location</td><td>Barcelona, Spain</td></tr>
<tr><td>Date</td><td>15/02/2024</td></tr>
<tr><td>Players</td><td>128</td></tr>
<tr><td>Rounds</td><td>9</td></tr>
<tr><td>Time control</td><td>90min+30sec</td></tr>
<tr><td>Organizer</td><td>Federació Catalana</td></tr>
</table>
</body></html>`

// -------------------------------------------------------------------------
// Helpers

func newTestLogger() *logger.Logger {
	return logger.New(io.Discard, logger.LevelInfo, "test", func(_ context.Context) string { return "" })
}

func newTestWorker(bus *fakeBus, pub *fakePublisher, serverURL string, locationFilter string) *Worker {
	httpClient := client.New(
		client.WithBaseURL(serverURL),
		client.WithDelay(0),
	)
	return New(Config{
		Log:            newTestLogger(),
		ScraperBus:     bus,
		Client:         httpClient,
		Publisher:      pub,
		LocationFilter: locationFilter,
	})
}

func pendingJob(federation string) scraperbus.ScrapeJob {
	now := time.Now()
	return scraperbus.ScrapeJob{
		ID:          uuid.New(),
		Federation:  federation,
		Status:      scraperbus.StatusPending,
		StartedAt:   now,
		DateCreated: now,
		DateUpdated: now,
	}
}

// -------------------------------------------------------------------------
// matchesLocation

func TestMatchesLocation(t *testing.T) {
	tests := []struct {
		location string
		filter   string
		want     bool
	}{
		{"Barcelona, Spain", "Barcelona", true},
		{"barcelona, spain", "BARCELONA", true}, // case-insensitive
		{"Girona, Spain", "Barcelona", false},
		{"", "Barcelona", false},
		{"L'Hospitalet de Llobregat, Barcelona", "Barcelona", true},
		{"Open de Barcelona 2024", "Barcelona", true},
	}

	for _, tt := range tests {
		got := matchesLocation(tt.location, tt.filter)
		if got != tt.want {
			t.Errorf("matchesLocation(%q, %q) = %v, want %v", tt.location, tt.filter, got, tt.want)
		}
	}
}

// -------------------------------------------------------------------------
// ProcessJob

func TestProcessJob_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/fed.aspx":
			w.Write([]byte(federationHTML))
		default:
			// All /tnr*.aspx requests get tournament detail HTML
			w.Write([]byte(tournamentHTML))
		}
	}))
	defer srv.Close()

	bus := newFakeBus()
	pub := &fakePublisher{}
	job := pendingJob("CAT")
	bus.addJob(job)

	w := newTestWorker(bus, pub, srv.URL, "Barcelona")

	if err := w.ProcessJob(context.Background(), job.ID); err != nil {
		t.Fatalf("ProcessJob: %v", err)
	}

	// Job should be in completed state
	final, _ := bus.QueryJobByID(context.Background(), job.ID)
	if final.Status != scraperbus.StatusCompleted {
		t.Errorf("job status: got %q, want %q", final.Status, scraperbus.StatusCompleted)
	}

	// Only the Barcelona tournament should have been upserted (Girona filtered out)
	if len(bus.tournaments) != 1 {
		t.Errorf("tournaments upserted: got %d, want 1", len(bus.tournaments))
	}
	if bus.tournaments[0].Name != "Copa Catalana 2024" {
		t.Errorf("tournament name: got %q, want %q", bus.tournaments[0].Name, "Copa Catalana 2024")
	}

	// At least one event should have been published
	pub.mu.Lock()
	eventCount := len(pub.events)
	pub.mu.Unlock()
	if eventCount == 0 {
		t.Error("expected at least one published event")
	}
}

func TestProcessJob_NoFilter_ScrapesBothTournaments(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/fed.aspx":
			w.Write([]byte(federationHTML))
		default:
			w.Write([]byte(tournamentHTML))
		}
	}))
	defer srv.Close()

	bus := newFakeBus()
	pub := &fakePublisher{}
	job := pendingJob("CAT")
	bus.addJob(job)

	w := newTestWorker(bus, pub, srv.URL, "") // no location filter

	if err := w.ProcessJob(context.Background(), job.ID); err != nil {
		t.Fatalf("ProcessJob: %v", err)
	}

	if len(bus.tournaments) != 2 {
		t.Errorf("tournaments upserted: got %d, want 2", len(bus.tournaments))
	}
}

func TestProcessJob_JobNotFound(t *testing.T) {
	bus := newFakeBus()
	pub := &fakePublisher{}
	w := newTestWorker(bus, pub, "http://localhost", "Barcelona")

	err := w.ProcessJob(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for unknown job ID, got nil")
	}
}

func TestProcessJob_FetchFailure_MarksJobFailed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	bus := newFakeBus()
	pub := &fakePublisher{}
	job := pendingJob("CAT")
	bus.addJob(job)

	// Use a very short retry delay to keep the test fast
	httpClient := client.New(
		client.WithBaseURL(srv.URL),
		client.WithDelay(0),
	)
	w := New(Config{
		Log:            newTestLogger(),
		ScraperBus:     bus,
		Client:         httpClient,
		Publisher:      pub,
		LocationFilter: "Barcelona",
	})

	err := w.ProcessJob(context.Background(), job.ID)
	if err == nil {
		t.Fatal("expected error from fetch failure, got nil")
	}

	final, _ := bus.QueryJobByID(context.Background(), job.ID)
	if final.Status != scraperbus.StatusFailed {
		t.Errorf("job status: got %q, want %q", final.Status, scraperbus.StatusFailed)
	}
	if final.Error == nil {
		t.Error("expected error message to be set on failed job")
	}

	// A scrape.error event should have been published
	pub.mu.Lock()
	published := pub.events
	pub.mu.Unlock()
	if len(published) == 0 {
		t.Error("expected scrape.error event to be published")
	}
	if published[0].Type != events.ScrapeError {
		t.Errorf("event type: got %q, want %q", published[0].Type, events.ScrapeError)
	}
}

func TestProcessJob_StatusTransitions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/fed.aspx":
			w.Write([]byte(federationHTML))
		default:
			w.Write([]byte(tournamentHTML))
		}
	}))
	defer srv.Close()

	bus := newFakeBus()
	pub := &fakePublisher{}
	job := pendingJob("CAT")
	bus.addJob(job)

	w := newTestWorker(bus, pub, srv.URL, "")

	if err := w.ProcessJob(context.Background(), job.ID); err != nil {
		t.Fatalf("ProcessJob: %v", err)
	}

	// Verify the sequence: first update sets running, second sets completed
	bus.mu.Lock()
	updates := bus.updates
	bus.mu.Unlock()

	if len(updates) < 2 {
		t.Fatalf("expected at least 2 UpdateJob calls, got %d", len(updates))
	}
	if updates[0].Status == nil || *updates[0].Status != scraperbus.StatusRunning {
		t.Errorf("first update should set status=running")
	}
	if updates[1].Status == nil || *updates[1].Status != scraperbus.StatusCompleted {
		t.Errorf("second update should set status=completed")
	}
}
