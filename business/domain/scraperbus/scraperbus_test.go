package scraperbus_test

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/ctfrancia/maple/business/domain/scraperbus"
	"github.com/ctfrancia/maple/business/sdk/order"
	"github.com/ctfrancia/maple/business/sdk/page"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/google/uuid"
)

// -------------------------------------------------------------------------
// In-memory Storer

var errNotFound = errors.New("not found")

type memStore struct {
	mu          sync.RWMutex
	jobs        map[uuid.UUID]scraperbus.ScrapeJob
	tournaments map[string]scraperbus.ScrapedTournament // keyed by ExternalID
}

func newMemStore() *memStore {
	return &memStore{
		jobs:        make(map[uuid.UUID]scraperbus.ScrapeJob),
		tournaments: make(map[string]scraperbus.ScrapedTournament),
	}
}

func (m *memStore) CreateJob(_ context.Context, job scraperbus.ScrapeJob) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobs[job.ID] = job
	return nil
}

func (m *memStore) UpdateJob(_ context.Context, job scraperbus.ScrapeJob, update scraperbus.UpdateScrapeJob) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	existing, ok := m.jobs[job.ID]
	if !ok {
		return errNotFound
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
	existing.DateUpdated = time.Now()
	m.jobs[job.ID] = existing
	return nil
}

func (m *memStore) QueryJobs(_ context.Context, filter scraperbus.QueryFilter, _ order.By, pg page.Page) ([]scraperbus.ScrapeJob, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []scraperbus.ScrapeJob
	for _, j := range m.jobs {
		if filter.Federation != nil && j.Federation != *filter.Federation {
			continue
		}
		if filter.Status != nil && j.Status != *filter.Status {
			continue
		}
		result = append(result, j)
	}
	// Simple pagination: offset = (page-1)*rows
	rows := pg.RowsPerPage()
	offset := (pg.Number() - 1) * rows
	if offset >= len(result) {
		return []scraperbus.ScrapeJob{}, nil
	}
	end := offset + rows
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}

func (m *memStore) QueryJobByID(_ context.Context, jobID uuid.UUID) (scraperbus.ScrapeJob, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	j, ok := m.jobs[jobID]
	if !ok {
		return scraperbus.ScrapeJob{}, errNotFound
	}
	return j, nil
}

func (m *memStore) UpsertTournament(_ context.Context, tournament scraperbus.ScrapedTournament) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// If external ID already exists, update; otherwise insert.
	if existing, ok := m.tournaments[tournament.ExternalID]; ok {
		tournament.ID = existing.ID
		tournament.DateCreated = existing.DateCreated
	}
	tournament.DateUpdated = time.Now()
	tournament.LastScrapedAt = time.Now()
	m.tournaments[tournament.ExternalID] = tournament
	return nil
}

func (m *memStore) QueryTournaments(_ context.Context, filter scraperbus.QueryFilter, _ order.By, pg page.Page) ([]scraperbus.ScrapedTournament, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []scraperbus.ScrapedTournament
	for _, t := range m.tournaments {
		if filter.Federation != nil && t.Federation != *filter.Federation {
			continue
		}
		result = append(result, t)
	}
	rows := pg.RowsPerPage()
	offset := (pg.Number() - 1) * rows
	if offset >= len(result) {
		return []scraperbus.ScrapedTournament{}, nil
	}
	end := offset + rows
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}

func (m *memStore) QueryTournamentByExternalID(_ context.Context, externalID string) (scraperbus.ScrapedTournament, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.tournaments[externalID]
	if !ok {
		return scraperbus.ScrapedTournament{}, errNotFound
	}
	return t, nil
}

// -------------------------------------------------------------------------
// Helpers

func newTestBusiness(t *testing.T) scraperbus.ExtBusiness {
	t.Helper()
	log := logger.New(io.Discard, logger.LevelInfo, "test", func(_ context.Context) string { return "" })
	return scraperbus.NewBusiness(log, newMemStore())
}

// -------------------------------------------------------------------------
// CreateJob

func TestCreateJob_SetsDefaults(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	job, err := bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "CAT"})
	if err != nil {
		t.Fatalf("CreateJob: %v", err)
	}

	if job.ID == uuid.Nil {
		t.Error("expected non-nil UUID")
	}
	if job.Federation != "CAT" {
		t.Errorf("federation: got %q, want %q", job.Federation, "CAT")
	}
	if job.Status != scraperbus.StatusPending {
		t.Errorf("status: got %q, want %q", job.Status, scraperbus.StatusPending)
	}
	if job.StartedAt.IsZero() {
		t.Error("StartedAt should be set")
	}
}

func TestCreateJob_MultipleJobsGetUniqueIDs(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	j1, _ := bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "CAT"})
	j2, _ := bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "ESP"})

	if j1.ID == j2.ID {
		t.Error("expected unique IDs for different jobs")
	}
}

// -------------------------------------------------------------------------
// QueryJobByID

func TestQueryJobByID_Found(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	created, _ := bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "CAT"})

	got, err := bus.QueryJobByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("QueryJobByID: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID: got %v, want %v", got.ID, created.ID)
	}
}

func TestQueryJobByID_NotFound(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	_, err := bus.QueryJobByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for unknown ID, got nil")
	}
}

// -------------------------------------------------------------------------
// UpdateJob

func TestUpdateJob_StatusTransitions(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	job, _ := bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "CAT"})

	running := scraperbus.StatusRunning
	if err := bus.UpdateJob(ctx, job, scraperbus.UpdateScrapeJob{Status: &running}); err != nil {
		t.Fatalf("UpdateJob to running: %v", err)
	}

	updated, _ := bus.QueryJobByID(ctx, job.ID)
	if updated.Status != scraperbus.StatusRunning {
		t.Errorf("status after update: got %q, want %q", updated.Status, scraperbus.StatusRunning)
	}

	completed := scraperbus.StatusCompleted
	now := time.Now()
	found, created, upd := 10, 8, 2
	if err := bus.UpdateJob(ctx, job, scraperbus.UpdateScrapeJob{
		Status:      &completed,
		CompletedAt: &now,
		Found:       &found,
		Created:     &created,
		Updated:     &upd,
	}); err != nil {
		t.Fatalf("UpdateJob to completed: %v", err)
	}

	final, _ := bus.QueryJobByID(ctx, job.ID)
	if final.Status != scraperbus.StatusCompleted {
		t.Errorf("final status: got %q, want %q", final.Status, scraperbus.StatusCompleted)
	}
	if final.Found != 10 || final.Created != 8 || final.Updated != 2 {
		t.Errorf("stats: got found=%d created=%d updated=%d", final.Found, final.Created, final.Updated)
	}
	if final.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}
}

func TestUpdateJob_ErrorField(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	job, _ := bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "CAT"})

	failed := scraperbus.StatusFailed
	errMsg := "connection refused"
	if err := bus.UpdateJob(ctx, job, scraperbus.UpdateScrapeJob{
		Status: &failed,
		Error:  &errMsg,
	}); err != nil {
		t.Fatalf("UpdateJob: %v", err)
	}

	got, _ := bus.QueryJobByID(ctx, job.ID)
	if got.Status != scraperbus.StatusFailed {
		t.Errorf("status: got %q, want %q", got.Status, scraperbus.StatusFailed)
	}
	if got.Error == nil || *got.Error != errMsg {
		t.Errorf("error field: got %v, want %q", got.Error, errMsg)
	}
}

// -------------------------------------------------------------------------
// QueryJobs

func TestQueryJobs_FilterByStatus(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	j1, _ := bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "CAT"})
	bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "ESP"})

	running := scraperbus.StatusRunning
	bus.UpdateJob(ctx, j1, scraperbus.UpdateScrapeJob{Status: &running})

	pendingStatus := scraperbus.StatusPending
	jobs, err := bus.QueryJobs(ctx, scraperbus.QueryFilter{Status: &pendingStatus}, order.By{}, page.MustParse("1", "100"))
	if err != nil {
		t.Fatalf("QueryJobs: %v", err)
	}
	if len(jobs) != 1 {
		t.Errorf("pending jobs: got %d, want 1", len(jobs))
	}
}

func TestQueryJobs_FilterByFederation(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "CAT"})
	bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "CAT"})
	bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "ESP"})

	fed := "CAT"
	jobs, err := bus.QueryJobs(ctx, scraperbus.QueryFilter{Federation: &fed}, order.By{}, page.MustParse("1", "100"))
	if err != nil {
		t.Fatalf("QueryJobs: %v", err)
	}
	if len(jobs) != 2 {
		t.Errorf("CAT jobs: got %d, want 2", len(jobs))
	}
}

func TestQueryJobs_NoFilter_ReturnsAll(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "CAT"})
	bus.CreateJob(ctx, scraperbus.NewScrapeJob{Federation: "ESP"})

	jobs, err := bus.QueryJobs(ctx, scraperbus.QueryFilter{}, order.By{}, page.MustParse("1", "100"))
	if err != nil {
		t.Fatalf("QueryJobs: %v", err)
	}
	if len(jobs) != 2 {
		t.Errorf("all jobs: got %d, want 2", len(jobs))
	}
}

// -------------------------------------------------------------------------
// UpsertTournament / QueryTournaments

func TestUpsertTournament_Create(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	loc := "Barcelona, Spain"
	tc := "90min+30sec"
	tournament := scraperbus.ScrapedTournament{
		ID:          uuid.New(),
		ExternalID:  "123456",
		Name:        "Copa Catalana 2024",
		URL:         "https://chess-results.com/tnr123456.aspx",
		Federation:  "CAT",
		Location:    &loc,
		Players:     128,
		Rounds:      9,
		TimeControl: &tc,
	}

	if err := bus.UpsertTournament(ctx, tournament); err != nil {
		t.Fatalf("UpsertTournament: %v", err)
	}

	all, err := bus.QueryTournaments(ctx, scraperbus.QueryFilter{}, order.By{}, page.MustParse("1", "100"))
	if err != nil {
		t.Fatalf("QueryTournaments: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("tournaments: got %d, want 1", len(all))
	}
	if all[0].Name != "Copa Catalana 2024" {
		t.Errorf("name: got %q, want %q", all[0].Name, "Copa Catalana 2024")
	}
}

func TestUpsertTournament_UpdatesExisting(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	loc := "Barcelona"
	bus.UpsertTournament(ctx, scraperbus.ScrapedTournament{
		ID:         uuid.New(),
		ExternalID: "123456",
		Name:       "Old Name",
		Federation: "CAT",
		Location:   &loc,
		Players:    64,
	})

	newLoc := "Barcelona, Spain"
	bus.UpsertTournament(ctx, scraperbus.ScrapedTournament{
		ID:         uuid.New(), // different ID, same ExternalID
		ExternalID: "123456",
		Name:       "Copa Catalana 2024",
		Federation: "CAT",
		Location:   &newLoc,
		Players:    128,
	})

	all, _ := bus.QueryTournaments(ctx, scraperbus.QueryFilter{}, order.By{}, page.MustParse("1", "100"))
	if len(all) != 1 {
		t.Fatalf("expected 1 tournament after upsert, got %d", len(all))
	}
	if all[0].Name != "Copa Catalana 2024" {
		t.Errorf("name after upsert: got %q, want %q", all[0].Name, "Copa Catalana 2024")
	}
	if all[0].Players != 128 {
		t.Errorf("players after upsert: got %d, want 128", all[0].Players)
	}
}

func TestQueryTournaments_FilterByFederation(t *testing.T) {
	bus := newTestBusiness(t)
	ctx := context.Background()

	bus.UpsertTournament(ctx, scraperbus.ScrapedTournament{ID: uuid.New(), ExternalID: "1", Name: "T1", Federation: "CAT"})
	bus.UpsertTournament(ctx, scraperbus.ScrapedTournament{ID: uuid.New(), ExternalID: "2", Name: "T2", Federation: "CAT"})
	bus.UpsertTournament(ctx, scraperbus.ScrapedTournament{ID: uuid.New(), ExternalID: "3", Name: "T3", Federation: "ESP"})

	fed := "CAT"
	results, err := bus.QueryTournaments(ctx, scraperbus.QueryFilter{Federation: &fed}, order.By{}, page.MustParse("1", "100"))
	if err != nil {
		t.Fatalf("QueryTournaments: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("CAT tournaments: got %d, want 2", len(results))
	}
}
