package main

import (
	"context"
	"expvar"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ctfrancia/maple/api/services/scrape/client"
	"github.com/ctfrancia/maple/api/services/scrape/config"
	"github.com/ctfrancia/maple/api/services/scrape/publisher"
	"github.com/ctfrancia/maple/api/services/scrape/worker"
	"github.com/ctfrancia/maple/app/domain/scraperapp"
	"github.com/ctfrancia/maple/app/sdk/debug"
	"github.com/ctfrancia/maple/app/sdk/mux"
	"github.com/ctfrancia/maple/business/domain/scraperbus"
	"github.com/ctfrancia/maple/business/domain/scraperbus/stores/scraperdb"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/ctfrancia/maple/foundation/otel"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var build = "develop"

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "**************** SEND ALERT ****************")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return otel.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SCRAPER", traceIDFn, events)

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup failed", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	env := os.Getenv("APP_ENV")
	if env == "" || env == "development" {
		if err := godotenv.Load("zoltan/compose/.env.dev"); err != nil {
			log.Info(ctx, "startup", "warn", "Warning: .env.dev file not found")
		}
	}

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	log.Info(ctx, "starting service", "version", build)
	defer log.Info(ctx, "shutdown complete")

	expvar.NewString("build").Set(build)

	// -------------------------------------------------------------------------
	// Database Setup

	log.Info(ctx, "startup", "status", "initializing database support")

	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Info),
	})
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}

	// Run migrations
	if err := scraperdb.CreateMigration(db); err != nil {
		return fmt.Errorf("creating migration: %w", err)
	}

	psql, err := db.DB()
	if err != nil {
		return fmt.Errorf("getting database connection: %w", err)
	}

	defer psql.Close()

	// -------------------------------------------------------------------------
	// Business Packages

	scraperStore, err := scraperdb.NewStore(log, db)
	if err != nil {
		return fmt.Errorf("creating scraper store: %w", err)
	}
	scraperBus := scraperbus.NewBusiness(log, scraperStore)

	// -------------------------------------------------------------------------
	// Scraper Components

	log.Info(ctx, "startup", "status", "initializing scraper components")

	// Create HTTP client for chess-results.com
	httpClient := client.New(
		client.WithBaseURL(cfg.BaseURL),
		client.WithUserAgent(cfg.UserAgent),
		client.WithDelay(cfg.ReqDelay),
		client.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
	)

	// Create publisher
	pub := publisher.NewLogPublisher(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Create worker
	workerInstance := worker.New(worker.Config{
		Log:            log,
		ScraperBus:     scraperBus,
		Client:         httpClient,
		Publisher:      pub,
		LocationFilter: cfg.LocationFilter,
	})

	log.Info(ctx, "worker configured",
		"federations", cfg.Federations,
		"location_filter", cfg.LocationFilter,
		"interval", cfg.DiscoveryInterval.String(),
	)

	// Create orchestrator
	orchestrator := worker.NewOrchestrator(worker.OrchestratorConfig{
		Log:        log,
		ScraperBus: scraperBus,
		Worker:     workerInstance,
	})

	// Create scheduler
	scheduler := worker.NewScheduler(worker.SchedulerConfig{
		Log:         log,
		ScraperBus:  scraperBus,
		Federations: cfg.Federations,
		Interval:    cfg.DiscoveryInterval,
	})

	// Start orchestrator
	if err := orchestrator.Start(ctx); err != nil {
		return fmt.Errorf("starting orchestrator: %w", err)
	}

	// Start scheduler
	if err := scheduler.Start(ctx); err != nil {
		return fmt.Errorf("starting scheduler: %w", err)
	}

	// -------------------------------------------------------------------------
	// Start Tracing Support

	log.Info(ctx, "startup", "status", "initializing tracing support")

	traceProvider, teardown, err := otel.InitTracing(log, otel.Config{
		ServiceName: "scraper",
		Host:        "tempo:4317",
		ExcludedRoutes: map[string]struct{}{
			"/v1/liveness":  {},
			"/v1/readiness": {},
		},
		Probability: 0.5,
	})
	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}

	defer teardown(context.Background())

	tracer := traceProvider.Tracer("scraper")

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.DebugHost)

		if err := http.ListenAndServe(cfg.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Configure routes - create a minimal mux without scraper routes
	cfgMux := mux.Config{
		CORSAllowedOrigins: cfg.CORSAllowedOrigins,
		Build:              build,
		Log:                log,
		Tracer:             tracer,
		DB:                 db,
		BusConfig: mux.BusConfig{
			// Don't set ScraperBus here, we'll mount it manually with orchestrator
		},
	}

	webAPI := mux.WebAPI(cfgMux)

	// Mount scraper routes with orchestrator
	scraperAppCfg := scraperapp.Config{
		Log:          log,
		ScraperBus:   scraperBus,
		Orchestrator: orchestrator,
	}

	// Type assert to chi.Router to mount routes
	if router, ok := webAPI.(interface{ Mount(string, http.Handler) }); ok {
		router.Mount("/api", scraperapp.V1Routes(scraperAppCfg))
	}

	api := http.Server{
		Addr:         cfg.APIHost,
		Handler:      webAPI,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		// Stop scheduler first
		if err := scheduler.Stop(shutdownCtx); err != nil {
			log.Error(shutdownCtx, "error stopping scheduler", "error", err)
		}

		// Stop orchestrator
		if err := orchestrator.Stop(shutdownCtx); err != nil {
			log.Error(shutdownCtx, "error stopping orchestrator", "error", err)
		}

		// Shutdown API server
		if err := api.Shutdown(shutdownCtx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
