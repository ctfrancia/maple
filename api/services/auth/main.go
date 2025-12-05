package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/ctfrancia/maple/app/sdk/debug"
	"github.com/ctfrancia/maple/app/sdk/mux"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/ctfrancia/maple/foundation/otel"
)

var build = "develop"

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "****** SEND ALERT ******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return otel.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "AUTH", traceIDFn, events)

	// -----------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// CONFIG

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:10s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:6000"`
			DebugHost          string        `conf:"default:0.0.0.0:6010"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}
		Auth struct {
			KeysJSON   string `conf:"mask"`
			KeysFolder string `conf:"default:zoltan/keys/"`
			ActiveKID  string `conf:"default:KID_1"`
			Issuer     string `conf:"default:maple"`
		}
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres,mask"`
			Host         string `conf:"default:database-service"`
			Name         string `conf:"default:postgres"`
			MaxIdleConns int    `conf:"default:0"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
		Tempo struct {
			Host        string  `conf:"default:tempo:4317"`
			ServiceName string  `conf:"default:auth"`
			Probability float64 `conf:"default:0.05"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Auth",
		},
	}

	const prefix = "AUTH"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	log.BuildInfo(ctx)

	expvar.NewString("build").Set(cfg.Build)

	// -------------------------------------------------------------------------
	// DB Setup

	log.Info(ctx, "startup", "status", "unintitializing database support - !TODO!", "hostport", nil)

	// -------------------------------------------------------------------------
	// Create Business Packages - !TODO!

	/*
		here will be the code related to the business of the application.
		Here will be the users? I think so because a "user" in the system is someone
		who is using the system which is different then a "player" in the system.

		I'll need to do some thinking and research this as I am unsure how I want to do this.
	*/

	// -------------------------------------------------------------------------
	// TRACING

	log.Info(ctx, "startup", "status", "initializing tracing support")

	traceProvider, teardown, err := otel.InitTracing(log, otel.Config{
		ServiceName: cfg.Tempo.ServiceName,
		Host:        cfg.Tempo.Host,
		ExcludedRoutes: map[string]struct{}{
			"/v1/liveness":  {},
			"/v1/readiness": {},
		},
		Probability: cfg.Tempo.Probability,
	})
	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}
	defer teardown(context.Background())

	tracer := traceProvider.Tracer(cfg.Tempo.ServiceName)

	// -------------------------------------------------------------------------
	// Debug

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		Build:     build,
		Log:       log,
		Tracer:    tracer,
		BusConfig: mux.BusConfig{
			//TournamentBus: tournamentBus,
		},
		AuthConfig: mux.AuthConfig{
			//Auth: ath,
		},
	}

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      mux.WebAPI(cfgMux),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()
	// -------------------------------------------------------------------------
	// SHUTDOWN

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
