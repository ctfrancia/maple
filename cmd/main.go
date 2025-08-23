package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	rest "github.com/ctfrancia/maple/internal/adapters/http"
	"github.com/ctfrancia/maple/internal/adapters/logger"
	"github.com/ctfrancia/maple/internal/adapters/persistence/inmemory"
	"github.com/ctfrancia/maple/internal/adapters/system"
	"github.com/ctfrancia/maple/internal/application/services"
	"github.com/ctfrancia/maple/internal/core/ports"
)

var (
	env                  = os.Getenv("ENV")
	listenAddress        = os.Getenv("LISTEN_ADDRESS")
	workerCount          = os.Getenv("WORKER_COUNT")
	sigServiceTimeout    = os.Getenv("SIG_SERVICE_TIMEOUT")
	readTimeout          = os.Getenv("READ_TIMEOUT")
	writeTimeout         = os.Getenv("WRITE_TIMEOUT")
	idleTimeout          = os.Getenv("IDLE_TIMEOUT")
	log                  ports.Logger
	tournamentRepository ports.TournamentRepository
	repoProvider         ports.TournamentRepositoryProvider
)

func main() {
	var rt, wt, it time.Duration
	if listenAddress == "" {
		panic("LISTEN_ADDRESS is not set")
	}
	if workerCount == "" {
		panic("WORKER_COUNT is not set")
	}
	if sigServiceTimeout == "" {
		panic("SIG_SERVICE_TIMEOUT is not set")
	}

	// Parse timeout first
	timeout, err := time.ParseDuration(sigServiceTimeout)
	if err != nil {
		panic(err)
	}

	// Create main application context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// environment specific setup
	switch env {
	case "prod":
		fmt.Println("using production environment")
		log = logger.NewZapLogger(env)
	case "dev", "test":
		fmt.Println("using dev|test environment")
		log = logger.NewZapLogger(env)
		tournamentRepository = inmemory.NewInMemoryTournamentRepository()
		repoProvider = inmemory.NewTournamentRepositoryProvider(tournamentRepository)
		rt = 15 * time.Second
		wt = 15 * time.Second
		it = 60 * time.Second
	default:
		log = logger.NewZapLogger("dev")
		fmt.Println("reached default using dev logger")
	}

	log.Info(context.Background(), "Starting server")

	// Adapters
	sa := system.NewSystemAdapter()

	// Services - use the main context
	wp := services.NewTournamentWorkerPool(ctx, cancel)
	wp.Start()
	defer wp.Stop()
	shs := services.NewSystemHealthServicer(sa, nil, nil)

	ts, err := services.NewTournamentServicer(log, repoProvider, wp)
	if err != nil {
		log.Error(context.Background(), "Tournament service creation failed", ports.Error("error", err))
		os.Exit(1)
	}

	// Create a new router
	// TODO: this will be moved to server.go file
	router := rest.NewRouter(log, shs, ts)
	srv := &http.Server{
		Addr:         listenAddress,
		Handler:      router,
		ReadTimeout:  rt,
		WriteTimeout: wt,
		IdleTimeout:  it,
	}

	// Channel to listen for interrupt signal to trigger shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine so it doesn't block
	go func() {
		log.Info(context.Background(), fmt.Sprintf("Server starting on %s", listenAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(context.Background(), "Failed to start server", ports.Error("error", err))
			cancel()
		}
	}()

	// Wait for either shutdown signal or context cancellation
	select {
	case <-quit:
		log.Info(context.Background(), "Shutdown signal received")
	case <-ctx.Done():
		log.Info(context.Background(), "Context cancelled, shutting down")
	}

	log.Info(context.Background(), "Shutting down server...")

	// Cancel the main context to signal all services to stop
	cancel()

	// Creates a deadline for shutdown (use the parsed timeout)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error(context.Background(), "Server forced to shutdown", ports.Error("error", err))
		os.Exit(1)
	}

	log.Info(context.Background(), "Server exited")
}
