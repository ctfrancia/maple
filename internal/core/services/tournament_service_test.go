package services

import (
	"context"
	"testing"
	"time"

	"github.com/ctfrancia/maple/internal/adapters/logger"
	"github.com/ctfrancia/maple/internal/adapters/repository/inmemory"
	"github.com/ctfrancia/maple/internal/core/domain"
)

var lggr = logger.NewZapLogger("test")

func TestCreateTournament_ShouldCreateService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wp := NewTournamentWorkerPool(ctx, cancel)
	wp.Start()
	defer wp.Stop()

	repo := inmemory.NewInMemoryTournamentRepository()

	_, err := NewTournamentServicer(lggr, repo, wp)
	if err != nil {
		t.Errorf("error creating service: %v", err)
	}
}

func TestCreateTournament_ShouldCreateTournament(t *testing.T) {
	repo := inmemory.NewInMemoryTournamentRepository()
	// ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	wp := NewTournamentWorkerPool(ctx, cancel)
	wp.Start()
	defer wp.Stop()

	ts, err := NewTournamentServicer(lggr, repo, wp)
	if err != nil {
		t.Errorf("error creating service: %v", err)
	}

	result, err := ts.CreateTournament(ctx, domain.Tournament{Name: "Test Tournament"})
	if err != nil {
		t.Errorf("error creating tournament: %v", err)
	}

	if result.Name != "Test Tournament" {
		t.Errorf("tournament name is not correct")
	}
}

func TestCreateTournament_ShouldListTournaments(t *testing.T) {
	repo := inmemory.NewInMemoryTournamentRepository()
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()

	wp := NewTournamentWorkerPool(ctx, cancel)
	wp.Start()
	defer wp.Stop()

	ts, err := NewTournamentServicer(lggr, repo, wp)
	if err != nil {
		t.Errorf("error creating service: %v", err)
	}

	// Create a tournament
	_, err = ts.CreateTournament(ctx, domain.Tournament{Name: "Test Tournament"})
	if err != nil {
		t.Errorf("error creating tournament: %v", err)
	}

	result, err := ts.ListTournaments(ctx)
	if err != nil {
		t.Errorf("error listing tournaments: %v", err)
	}

	if len(result) == 1 {
		t.Errorf("tournaments list is empty")
	}
}
