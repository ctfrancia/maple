package services

import (
	"context"
	"fmt"

	// "github.com/ctfrancia/maple/internal/adapters/repository/inmemory"
	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/ctfrancia/maple/internal/core/ports"

	"github.com/google/uuid"
)

type TournamentServicer struct {
	logger     ports.Logger
	repository ports.TournamentRepository
	workerPool *TournamentWorkerPool // make this port
}

func NewTournamentServicer(log ports.Logger, tr ports.TournamentRepository, wp *TournamentWorkerPool) (ports.TournamentServicer, error) {
	if log == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if tr == nil {
		return nil, fmt.Errorf("tournament repository is required")
	}
	if wp == nil {
		return nil, fmt.Errorf("tournament worker pool is required")
	}
	return &TournamentServicer{
		logger:     log,
		repository: tr,
		workerPool: wp,
	}, nil
}

func (ts *TournamentServicer) CreateTournament(ctx context.Context, tournament domain.Tournament) (domain.Tournament, error) {
	task := TournamentTask{
		ID:         uuid.New(),
		Type:       TaskTypeCreateTournament,
		Data:       CreateTournamentTask{Tournament: tournament},
		Repository: ts.repository,
		ResultCh:   make(chan TaskResult, 1),
		Context:    ctx,
	}

	resultCh := ts.workerPool.SubmitTask(task)

	select {
	case result := <-resultCh:
		if result.Error != nil {
			// error handling here
			return domain.Tournament{}, result.Error
		}
		return result.Data.(domain.Tournament), nil

	case <-ctx.Done():
		return domain.Tournament{}, ctx.Err()
	}
}

func (ts *TournamentServicer) ListTournaments(ctx context.Context) ([]domain.Tournament, error) {
	task := TournamentTask{
		ID:         uuid.New(),
		Type:       TaskTypeListTournaments,
		Data:       ListTournamentsTask{},
		Repository: ts.repository,
		ResultCh:   make(chan TaskResult, 1),
		Context:    ctx,
	}

	resultCh := ts.workerPool.SubmitTask(task)

	select {
	case result := <-resultCh:
		if result.Error != nil {
			// error handling here
			return nil, result.Error
		}
		return result.Data.([]domain.Tournament), nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
