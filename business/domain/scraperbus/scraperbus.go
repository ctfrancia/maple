package scraperbus

import (
	"context"
	"fmt"
	"time"

	"github.com/ctfrancia/maple/business/sdk/order"
	"github.com/ctfrancia/maple/business/sdk/page"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/ctfrancia/maple/foundation/otel"
	"github.com/google/uuid"
)

// Storer is the interface the business requires to access data.
type Storer interface {
	CreateJob(ctx context.Context, job ScrapeJob) error
	UpdateJob(ctx context.Context, job ScrapeJob, update UpdateScrapeJob) error
	QueryJobs(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]ScrapeJob, error)
	QueryJobByID(ctx context.Context, jobID uuid.UUID) (ScrapeJob, error)

	UpsertTournament(ctx context.Context, tournament ScrapedTournament) error
	QueryTournaments(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]ScrapedTournament, error)
	QueryTournamentByExternalID(ctx context.Context, externalID string) (ScrapedTournament, error)
}

// ExtBusiness defines the external interface for the scraper business logic.
type ExtBusiness interface {
	CreateJob(ctx context.Context, nj NewScrapeJob) (ScrapeJob, error)
	UpdateJob(ctx context.Context, job ScrapeJob, update UpdateScrapeJob) error
	QueryJobs(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]ScrapeJob, error)
	QueryJobByID(ctx context.Context, jobID uuid.UUID) (ScrapeJob, error)

	UpsertTournament(ctx context.Context, tournament ScrapedTournament) error
	QueryTournaments(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]ScrapedTournament, error)
}

// Extension wraps additional business logic around an existing one.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for scraper access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a scraper business API for use.
func NewBusiness(log *logger.Logger, storer Storer, extensions ...Extension) ExtBusiness {
	b := ExtBusiness(&Business{
		log:    log,
		storer: storer,
	})

	for i := len(extensions) - 1; i >= 0; i-- {
		ext := extensions[i]
		if ext != nil {
			b = ext(b)
		}
	}

	return b
}

// CreateJob creates a new scrape job.
func (b *Business) CreateJob(ctx context.Context, nj NewScrapeJob) (ScrapeJob, error) {
	ctx, span := otel.AddSpan(ctx, "business.scraperbus.CreateJob")
	defer span.End()

	now := time.Now()
	job := ScrapeJob{
		ID:          uuid.New(),
		Federation:  nj.Federation,
		Status:      StatusPending,
		StartedAt:   now,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.CreateJob(ctx, job); err != nil {
		return ScrapeJob{}, fmt.Errorf("createJob: %w", err)
	}

	return job, nil
}

// UpdateJob updates an existing scrape job.
func (b *Business) UpdateJob(ctx context.Context, job ScrapeJob, update UpdateScrapeJob) error {
	ctx, span := otel.AddSpan(ctx, "business.scraperbus.UpdateJob")
	defer span.End()

	if err := b.storer.UpdateJob(ctx, job, update); err != nil {
		return fmt.Errorf("updateJob: %w", err)
	}

	return nil
}

// QueryJobs retrieves a list of scrape jobs based on the provided filter.
func (b *Business) QueryJobs(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]ScrapeJob, error) {
	ctx, span := otel.AddSpan(ctx, "business.scraperbus.QueryJobs")
	defer span.End()

	jobs, err := b.storer.QueryJobs(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("queryJobs: %w", err)
	}

	return jobs, nil
}

// QueryJobByID retrieves a scrape job by ID.
func (b *Business) QueryJobByID(ctx context.Context, jobID uuid.UUID) (ScrapeJob, error) {
	ctx, span := otel.AddSpan(ctx, "business.scraperbus.QueryJobByID")
	defer span.End()

	job, err := b.storer.QueryJobByID(ctx, jobID)
	if err != nil {
		return ScrapeJob{}, fmt.Errorf("queryJobByID: %w", err)
	}

	return job, nil
}

// UpsertTournament creates or updates a scraped tournament.
func (b *Business) UpsertTournament(ctx context.Context, tournament ScrapedTournament) error {
	ctx, span := otel.AddSpan(ctx, "business.scraperbus.UpsertTournament")
	defer span.End()

	if err := b.storer.UpsertTournament(ctx, tournament); err != nil {
		return fmt.Errorf("upsertTournament: %w", err)
	}

	return nil
}

// QueryTournaments retrieves a list of scraped tournaments based on the provided filter.
func (b *Business) QueryTournaments(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]ScrapedTournament, error) {
	ctx, span := otel.AddSpan(ctx, "business.scraperbus.QueryTournaments")
	defer span.End()

	tournaments, err := b.storer.QueryTournaments(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("queryTournaments: %w", err)
	}

	return tournaments, nil
}
