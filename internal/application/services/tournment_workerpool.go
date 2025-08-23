package services

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	commands "github.com/ctfrancia/maple/internal/application/commands/tournament"
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
	Repository ports.TournamentRepositoryProvider

	Context context.Context
}

type CreateTournamentTask struct {
	Tournament commands.CreateTournamentCommand
}

type FindTournamentTask struct {
	TournamentID uuid.UUID
}

type ListTournamentsTask struct{}

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
			var result TaskResult

			switch task.Type {
			case TaskTypeCreateTournament:
				result = twp.createTournament(task)

			case TaskTypeFindTournament:
				result = twp.findTournament(task)

			case TaskTypeListTournaments:
				result = twp.listTournaments(task)

			default:
				result = TaskResult{Error: fmt.Errorf("invalid task type")}
			}

			task.ResultCh <- result // ← Only one send here

		case <-twp.ctx.Done():
			twp.wg.Done()
			return
		}
	}
}

func (twp *TournamentWorkerPool) Stop() {
	twp.mu.Lock()
	defer twp.mu.Unlock()

	if !twp.started {
		return
	}

	twp.cancel()
	close(twp.taskQueue)
	twp.wg.Wait()
	twp.started = false
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
	var result domain.Tournament
	var err error
	t, ok := task.Data.(CreateTournamentTask)
	if !ok {
		task.ResultCh <- TaskResult{Error: fmt.Errorf("invalid task data")}
		return TaskResult{Error: fmt.Errorf("invalid task data")}
	}

	tournament := domain.NewTournament(t.Tournament.Name, t.Tournament.Description)

	err = task.Repository.WriteTx(func(repo ports.TournamentRepository) error {
		result, err = repo.CreateTournament(*tournament)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return TaskResult{Error: err}
	}

	return TaskResult{Data: result}
}

func (twp *TournamentWorkerPool) findTournament(task TournamentTask) TaskResult {
	var result domain.Tournament
	var err error
	t, ok := task.Data.(FindTournamentTask)
	if !ok {
		return TaskResult{Error: fmt.Errorf("invalid task data")}
	}

	err = task.Repository.ReadTx(func(repo ports.TournamentRepository) error {
		result, err = repo.FindTournament(t.TournamentID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return TaskResult{Error: fmt.Errorf("error finding tournament: %v", err)}
	}

	return TaskResult{Data: result}
}

func (twp *TournamentWorkerPool) listTournaments(task TournamentTask) TaskResult {
	var results []domain.Tournament
	var err error
	_, ok := task.Data.(ListTournamentsTask)
	if !ok {
		return TaskResult{Error: fmt.Errorf("invalid task data")}
	}

	err = task.Repository.ReadTx(func(repo ports.TournamentRepository) error {
		results, err = repo.ListTournaments(nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return TaskResult{Error: fmt.Errorf("error listing tournaments: %v", err)}
	}

	return TaskResult{Data: results}
}
