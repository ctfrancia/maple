package services

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/ctfrancia/maple/internal/core/ports"
	"github.com/google/uuid"
)

type TaskType string

const (
	TaskTypeCreateTournament TaskType = "create_tournament"
	TaskTypeFindTournament   TaskType = "find_tournament"
	TaskTypeListTournaments  TaskType = "list_tournaments"
)

type TournamentWorkerPool struct {
	workers   int
	taskQueue chan TournamentTask
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	started   bool
	mu        sync.RWMutex
}

// TaskResult represents the result of a task execution
type TaskResult struct {
	Data  any
	Error error
}

type TournamentTask struct {
	ID         uuid.UUID
	Type       TaskType
	Data       any // will be converted to a specific task type in handler
	ResultCh   chan TaskResult
	Repository ports.TournamentRepository
	Context    context.Context
}

type CreateTournamentTask struct {
	Tournament domain.Tournament
	// Repository ports.TournamentRepository
}

type FindTournamentTask struct {
	TournamentID string
}

type ListTournamentsTask struct {
	Tournaments []domain.Tournament
}

func NewTournamentWorkerPool(ctx context.Context, cancel context.CancelFunc) *TournamentWorkerPool {
	return &TournamentWorkerPool{
		workers:   runtime.NumCPU() * 2,
		taskQueue: make(chan TournamentTask),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (twp *TournamentWorkerPool) Start() {
	twp.mu.Lock()
	defer twp.mu.Unlock()

	if twp.started {
		return
	}

	for i := 0; i < twp.workers; i++ {
		twp.wg.Add(1)
		go twp.worker()
	}
	twp.started = true
}

func (twp *TournamentWorkerPool) worker() {
	for {
		select {
		case task := <-twp.taskQueue:
			switch task.Type {
			case TaskTypeCreateTournament:
				result := twp.createTournament(task)
				task.ResultCh <- result
			case TaskTypeFindTournament:
				// twp.findTournament(task)
			case TaskTypeListTournaments:
				// twp.listTournaments(task)
			}
		case <-twp.ctx.Done():
			twp.wg.Done()
			return // Context's Done() called, so we can close the worker
		}
	}
}

func (twp *TournamentWorkerPool) Stop() {
	// TODO: stop the workers
}

func (twp *TournamentWorkerPool) SubmitTask(task TournamentTask) <-chan TaskResult {
	twp.mu.RLock()
	defer twp.mu.RUnlock()

	select {
	case twp.taskQueue <- task:
		return task.ResultCh
	case <-twp.ctx.Done():
		resultCh := make(chan TaskResult, 1)
		resultCh <- TaskResult{Error: fmt.Errorf("worker pool shutting down")}
		close(resultCh)
		return resultCh
	default:
		resultCh := make(chan TaskResult, 1)
		resultCh <- TaskResult{Error: fmt.Errorf("worker pool queue full")}
		close(resultCh)
		return resultCh
	}
}

func (twp *TournamentWorkerPool) createTournament(task TournamentTask) TaskResult {
	t, ok := task.Data.(CreateTournamentTask)
	if !ok {
		task.ResultCh <- TaskResult{Error: fmt.Errorf("invalid task data")}
		return TaskResult{Error: fmt.Errorf("invalid task data")}
	}

	fmt.Printf("task %#v", t)
	result, err := task.Repository.CreateTournament(t.Tournament)
	if err != nil {
		// task.ResultCh <- TaskResult{Error: err}
		return TaskResult{Error: err}
	}

	// task.ResultCh <- TaskResult{Data: result}
	return TaskResult{Data: result}
}
