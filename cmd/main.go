package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ctfrancia/maple/internal/adapters/logger"
	"github.com/ctfrancia/maple/internal/adapters/rest"
	"github.com/ctfrancia/maple/internal/adapters/system"
	"github.com/ctfrancia/maple/internal/core/services"
)

func main() {
	// create ancillary services
	logger := logger.NewZapLogger("dev")

	// Adapters
	sa := system.NewSystemAdapter()

	// Services
	shs := services.NewSystemHealthServicer(sa)

	// Create router
	router := rest.NewRouter(shs, logger)

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port, // TODO: come from config
		Handler:      router,
		ReadTimeout:  15 * time.Second, // TODO: come from config
		WriteTimeout: 15 * time.Second, // TODO: come from config
		IdleTimeout:  60 * time.Second, // TODO: come from config
	}

	// Start server in a goroutine so it doesn't block
	go func() {
		logger.Info(nil, "Starting server on port "+port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(nil, "Server failed to start: "+err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(nil, "Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(ctx, "Server forced to shutdown: "+err.Error())
	}

	logger.Info(ctx, "Server stopped gracefully")
}
