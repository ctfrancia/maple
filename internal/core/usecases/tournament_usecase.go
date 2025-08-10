package usecases

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type JobType string

const (
	JobTypeCreateTournament JobType = "create_tournament"
	JobTypeFindTournament   JobType = "find_tournament"
	JobTypeFindTournaments  JobType = "find_tournaments"
)

// Specific result types for type safety
type CreateTournamentResult struct {
	Tournament Tournament
	Error      error
}

type FindTournamentResult struct {
	Tournament Tournament
	Error      error
}

type FindTournamentsResult struct {
	Tournaments []Tournament
	Error       error
}

// Job payloads - type safe
type CreateTournamentPayload struct {
	Tournament Tournament
}

type FindTournamentPayload struct {
	ID int
}

type FindTournamentsPayload struct {
	Limit  int
	Offset int
}

// Generic job structure
type TournamentJob struct {
	ID       uuid.UUID
	Type     JobType
	Payload  any
	Context  context.Context
	resultCh chan any // Will contain specific result types
}

type Tournament struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Repository interface for better testability
type TournamentRepository interface {
	Create(ctx context.Context, tournament Tournament) (Tournament, error)
	FindByID(ctx context.Context, id int) (Tournament, error)
	FindAll(ctx context.Context, limit, offset int) ([]Tournament, error)
}

// Logger interface
type Logger interface {
	Info(ctx context.Context, msg string)
	Error(ctx context.Context, msg string)
}

type TournamentUseCase struct {
	repository  TournamentRepository
	logger      Logger
	wg          sync.WaitGroup
	workerCount int
	started     bool
	quit        chan struct{}
	jobs        chan TournamentJob
	mu          sync.RWMutex
}

func NewTournamentUseCase(repository TournamentRepository, logger Logger, workerCount int) *TournamentUseCase {
	if workerCount <= 0 {
		workerCount = 4 // Default worker count
	}

	return &TournamentUseCase{
		repository:  repository,
		logger:      logger,
		workerCount: workerCount,
		jobs:        make(chan TournamentJob, workerCount*2), // Buffer for bursts
		quit:        make(chan struct{}),
	}
}

// Start begins the worker pool
func (tuc *TournamentUseCase) Start() {
	tuc.mu.Lock()
	defer tuc.mu.Unlock()

	if tuc.started {
		return
	}

	// Start worker goroutines
	for i := 0; i < tuc.workerCount; i++ {
		tuc.wg.Add(1)
		go tuc.worker(i)
	}

	tuc.started = true
	if tuc.logger != nil {
		tuc.logger.Info(context.Background(), fmt.Sprintf("TournamentUseCase started with %d workers", tuc.workerCount))
	}
}

// Stop gracefully shuts down the worker pool
func (tuc *TournamentUseCase) Stop(timeout time.Duration) error {
	tuc.mu.Lock()
	if !tuc.started {
		tuc.mu.Unlock()
		return nil
	}
	tuc.mu.Unlock()

	// Signal workers to stop
	close(tuc.quit)

	// Wait for workers with timeout
	done := make(chan struct{})
	go func() {
		defer close(done)
		tuc.wg.Wait()
	}()

	select {
	case <-done:
		if tuc.logger != nil {
			tuc.logger.Info(context.Background(), "TournamentUseCase stopped gracefully")
		}
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for workers to stop")
	}
}

func (tuc *TournamentUseCase) worker(id int) {
	defer tuc.wg.Done()

	if tuc.logger != nil {
		tuc.logger.Info(context.Background(), fmt.Sprintf("Worker %d started", id))
	}

	for {
		select {
		case job := <-tuc.jobs:
			tuc.processJob(id, job)
		case <-tuc.quit:
			if tuc.logger != nil {
				tuc.logger.Info(context.Background(), fmt.Sprintf("Worker %d stopping", id))
			}
			return
		}
	}
}

func (tuc *TournamentUseCase) processJob(workerID int, job TournamentJob) {
	defer func() {
		// Ensure result channel is always closed
		close(job.resultCh)
	}()

	var result any

	switch job.Type {
	case JobTypeCreateTournament:
		payload, ok := job.Payload.(CreateTournamentPayload)
		if !ok {
			result = CreateTournamentResult{
				Error: fmt.Errorf("invalid payload type for create tournament"),
			}
			break
		}
		tournament, err := tuc.repository.Create(job.Context, payload.Tournament)
		result = CreateTournamentResult{
			Tournament: tournament,
			Error:      err,
		}

	case JobTypeFindTournament:
		payload, ok := job.Payload.(FindTournamentPayload)
		if !ok {
			result = FindTournamentResult{
				Error: fmt.Errorf("invalid payload type for find tournament"),
			}
			break
		}
		tournament, err := tuc.repository.FindByID(job.Context, payload.ID)
		result = FindTournamentResult{
			Tournament: tournament,
			Error:      err,
		}

	case JobTypeFindTournaments:
		payload, ok := job.Payload.(FindTournamentsPayload)
		if !ok {
			result = FindTournamentsResult{
				Error: fmt.Errorf("invalid payload type for find tournaments"),
			}
			break
		}
		tournaments, err := tuc.repository.FindAll(job.Context, payload.Limit, payload.Offset)
		result = FindTournamentsResult{
			Tournaments: tournaments,
			Error:       err,
		}

	default:
		result = CreateTournamentResult{
			Error: fmt.Errorf("unknown job type: %s", job.Type),
		}
	}

	// Send result back through result channel with context cancellation check
	select {
	case job.resultCh <- result:
	case <-job.Context.Done():
		if tuc.logger != nil {
			tuc.logger.Info(job.Context, fmt.Sprintf("Job %s context cancelled before result could be sent", job.ID))
		}
	}
}

// submitJob submits a job to the queue and waits for the result
func (tuc *TournamentUseCase) submitJob(ctx context.Context, jobType JobType, payload any) (any, error) {
	tuc.mu.RLock()
	if !tuc.started {
		tuc.mu.RUnlock()
		return nil, fmt.Errorf("tournament use case not started")
	}

	tuc.mu.RUnlock()

	// Create result channel (unbuffered for simplicity)
	resultCh := make(chan any)

	job := TournamentJob{
		ID:       uuid.New(),
		Type:     jobType,
		Payload:  payload,
		resultCh: resultCh,
		Context:  ctx,
	}

	// Submit job
	select {
	case tuc.jobs <- job:
		// Job submitted successfully
	case <-ctx.Done():
		close(resultCh) // Clean up
		return nil, ctx.Err()
	}

	// Wait for result
	select {
	case result := <-resultCh:
		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Public API methods

func (tuc *TournamentUseCase) CreateTournament(ctx context.Context, tournament Tournament) (Tournament, error) {
	result, err := tuc.submitJob(ctx, JobTypeCreateTournament, CreateTournamentPayload{
		Tournament: tournament,
	})
	if err != nil {
		return Tournament{}, err
	}

	createResult, ok := result.(CreateTournamentResult)
	if !ok {
		return Tournament{}, fmt.Errorf("unexpected result type")
	}

	return createResult.Tournament, createResult.Error
}

func (tuc *TournamentUseCase) FindTournament(ctx context.Context, id int) (Tournament, error) {
	result, err := tuc.submitJob(ctx, JobTypeFindTournament, FindTournamentPayload{
		ID: id,
	})
	if err != nil {
		return Tournament{}, err
	}

	findResult, ok := result.(FindTournamentResult)
	if !ok {
		return Tournament{}, fmt.Errorf("unexpected result type")
	}

	return findResult.Tournament, findResult.Error
}

func (tuc *TournamentUseCase) FindTournaments(ctx context.Context, limit, offset int) ([]Tournament, error) {
	result, err := tuc.submitJob(ctx, JobTypeFindTournaments, FindTournamentsPayload{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, err
	}

	findResult, ok := result.(FindTournamentsResult)
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	return findResult.Tournaments, findResult.Error
}

/*
// Synchronous methods (bypassing worker pool for simpler use cases)
func (tuc *TournamentUseCase) CreateTournamentSync(ctx context.Context, tournament Tournament) (Tournament, error) {
	if tuc.repository == nil {
		return Tournament{}, fmt.Errorf("repository not initialized")
	}
	return tuc.repository.Create(ctx, tournament)
}

func (tuc *TournamentUseCase) FindTournamentSync(ctx context.Context, id int) (Tournament, error) {
	if tuc.repository == nil {
		return Tournament{}, fmt.Errorf("repository not initialized")
	}
	return tuc.repository.FindByID(ctx, id)
}

func (tuc *TournamentUseCase) FindTournamentsSync(ctx context.Context, limit, offset int) ([]Tournament, error) {
	if tuc.repository == nil {
		return nil, fmt.Errorf("repository not initialized")
	}
	return tuc.repository.FindAll(ctx, limit, offset)
}
*/
