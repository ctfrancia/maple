package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ctfrancia/maple/internal/adapters/logger"
	"github.com/ctfrancia/maple/internal/adapters/repository/inmemory"
	"github.com/ctfrancia/maple/internal/adapters/rest"
	"github.com/ctfrancia/maple/internal/adapters/system"
	"github.com/ctfrancia/maple/internal/core/ports"
	"github.com/ctfrancia/maple/internal/core/services"
)

var (
	env                  = os.Getenv("ENV")
	listenAddress        = os.Getenv("LISTEN_ADDRESS")
	workerCount          = os.Getenv("WORKER_COUNT")
	sigServiceTimeout    = os.Getenv("SIG_SERVICE_TIMEOUT")
	log                  ports.Logger
	tournamentRepository ports.TournamentRepository
)

func main() {
	if listenAddress == "" {
		panic("LISTEN_ADDRESS is not set")
	}
	if workerCount == "" {
		panic("WORKER_COUNT is not set")
	}
	if sigServiceTimeout == "" {
		panic("SIG_SERVICE_TIMEOUT is not set")
	}

	// environment specific setup
	switch env {
	case "prod":
		fmt.Println("using production environment")
		log = logger.NewZapLogger(env)

	case "dev", "test":
		fmt.Println("using dev|test environment")
		log = logger.NewZapLogger(env)
		tournamentRepository = inmemory.NewInMemoryTournamentRepository()

	default:
		log = logger.NewZapLogger("dev")
		fmt.Println("reached default using dev logger")
	}

	log.Info(context.Background(), "Starting server")

	// Adapters
	sa := system.NewSystemAdapter()

	// Services
	shs := services.NewSystemHealthServicer(sa, nil, nil)
	ts, err := services.NewTournamentServicer(log, tournamentRepository, nil)
	if err != nil {
		log.Error(context.Background(), "Tournament service creation failed", ports.Error("error", err))
		os.Exit(1)
	}

	// Create a new router
	router := rest.NewRouter(log, shs, ts)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for interrupt signal to trigger shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine so it doesn't block
	go func() {
		log.Info(context.Background(), "Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(context.Background(), "Failed to start server", ports.Error("error", err))
			os.Exit(1)
		}
	}()

	<-quit
	log.Info(context.Background(), "Shutting down server...")

	// Creates a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Error(context.Background(), "Server forced to shutdown", ports.Error("error", err))
		os.Exit(1)
	}

	log.Info(context.Background(), "Server exited")
}
