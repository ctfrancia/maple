package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ctfrancia/maple/business/domain/scraperbus"
	"github.com/ctfrancia/maple/business/sdk/order"
	"github.com/ctfrancia/maple/business/sdk/page"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/google/uuid"
)

// Orchestrator manages scrape job processing
type Orchestrator struct {
	log        *logger.Logger
	scraperBus scraperbus.ExtBusiness
	worker     *Worker
	running    bool
	mu         sync.Mutex
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// OrchestratorConfig holds the orchestrator configuration
type OrchestratorConfig struct {
	Log        *logger.Logger
	ScraperBus scraperbus.ExtBusiness
	Worker     *Worker
}

// NewOrchestrator creates a new orchestrator
func NewOrchestrator(cfg OrchestratorConfig) *Orchestrator {
	return &Orchestrator{
		log:        cfg.Log,
		scraperBus: cfg.ScraperBus,
		worker:     cfg.Worker,
		stopCh:     make(chan struct{}),
	}
}

// Start begins processing jobs
func (o *Orchestrator) Start(ctx context.Context) error {
	o.mu.Lock()
	if o.running {
		o.mu.Unlock()
		return fmt.Errorf("orchestrator already running")
	}
	o.running = true
	o.mu.Unlock()

	o.log.Info(ctx, "orchestrator started")

	o.wg.Add(1)
	go o.processLoop(ctx)

	return nil
}

// Stop stops processing jobs
func (o *Orchestrator) Stop(ctx context.Context) error {
	o.mu.Lock()
	if !o.running {
		o.mu.Unlock()
		return nil
	}
	o.running = false
	o.mu.Unlock()

	o.log.Info(ctx, "stopping orchestrator")
	close(o.stopCh)
	o.wg.Wait()
	o.log.Info(ctx, "orchestrator stopped")

	return nil
}

// processLoop continuously processes pending jobs
func (o *Orchestrator) processLoop(ctx context.Context) {
	defer o.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-o.stopCh:
			return
		case <-ticker.C:
			o.processPendingJobs(ctx)
		}
	}
}

// processPendingJobs finds and processes all pending jobs
func (o *Orchestrator) processPendingJobs(ctx context.Context) {
	// Query for pending jobs
	pendingStatus := scraperbus.StatusPending
	filter := scraperbus.QueryFilter{
		Status: &pendingStatus,
	}

	jobs, err := o.scraperBus.QueryJobs(ctx, filter, order.By{Field: "date_created", Direction: order.ASC}, page.MustParse("1", "10"))
	if err != nil {
		o.log.Error(ctx, "failed to query pending jobs", "error", err)
		return
	}

	if len(jobs) == 0 {
		return
	}

	o.log.Info(ctx, "found pending jobs", "count", len(jobs))

	// Process each job
	for _, job := range jobs {
		if err := o.worker.ProcessJob(ctx, job.ID); err != nil {
			o.log.Error(ctx, "failed to process job", "job_id", job.ID.String(), "error", err)
		}
	}
}

// TriggerScrape creates and immediately processes a scrape job for a federation
func (o *Orchestrator) TriggerScrape(ctx context.Context, federation string) (uuid.UUID, error) {
	o.log.Info(ctx, "triggering manual scrape", "federation", federation)

	// Create a new job
	job, err := o.scraperBus.CreateJob(ctx, scraperbus.NewScrapeJob{
		Federation: federation,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("creating job: %w", err)
	}

	// Process it immediately in a goroutine
	go func() {
		jobCtx := context.Background()
		if err := o.worker.ProcessJob(jobCtx, job.ID); err != nil {
			o.log.Error(jobCtx, "failed to process triggered job", "job_id", job.ID.String(), "error", err)
		}
	}()

	return job.ID, nil
}
