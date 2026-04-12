package worker

import (
	"context"
	"sync"
	"time"

	"github.com/ctfrancia/maple/business/domain/scraperbus"
	"github.com/ctfrancia/maple/foundation/logger"
)

// Scheduler creates periodic scrape jobs
type Scheduler struct {
	log          *logger.Logger
	scraperBus   scraperbus.ExtBusiness
	federations  []string
	interval     time.Duration
	running      bool
	mu           sync.Mutex
	stopCh       chan struct{}
	wg           sync.WaitGroup
}

// SchedulerConfig holds the scheduler configuration
type SchedulerConfig struct {
	Log         *logger.Logger
	ScraperBus  scraperbus.ExtBusiness
	Federations []string
	Interval    time.Duration
}

// NewScheduler creates a new scheduler
func NewScheduler(cfg SchedulerConfig) *Scheduler {
	return &Scheduler{
		log:         cfg.Log,
		scraperBus:  cfg.ScraperBus,
		federations: cfg.Federations,
		interval:    cfg.Interval,
		stopCh:      make(chan struct{}),
	}
}

// Start begins the scheduler
func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = true
	s.mu.Unlock()

	s.log.Info(ctx, "scheduler started",
		"federations", s.federations,
		"interval", s.interval.String(),
	)

	s.wg.Add(1)
	go s.scheduleLoop(ctx)

	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop(ctx context.Context) error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = false
	s.mu.Unlock()

	s.log.Info(ctx, "stopping scheduler")
	close(s.stopCh)
	s.wg.Wait()
	s.log.Info(ctx, "scheduler stopped")

	return nil
}

// scheduleLoop creates periodic scrape jobs
func (s *Scheduler) scheduleLoop(ctx context.Context) {
	defer s.wg.Done()

	// Create initial jobs immediately
	s.createJobs(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.createJobs(ctx)
		}
	}
}

// createJobs creates a scrape job for each configured federation
func (s *Scheduler) createJobs(ctx context.Context) {
	s.log.Info(ctx, "creating scheduled scrape jobs", "federations", len(s.federations))

	for _, federation := range s.federations {
		job, err := s.scraperBus.CreateJob(ctx, scraperbus.NewScrapeJob{
			Federation: federation,
		})
		if err != nil {
			s.log.Error(ctx, "failed to create scheduled job",
				"federation", federation,
				"error", err,
			)
			continue
		}

		s.log.Info(ctx, "created scheduled job",
			"job_id", job.ID.String(),
			"federation", federation,
		)
	}
}
