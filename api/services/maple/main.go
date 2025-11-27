// Package maple is the main package for the maple architecture.
package maple

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
	"github.com/ctfrancia/maple/business/domain/tournamentbus"
	"github.com/ctfrancia/maple/business/domain/tournamentbus/extensions/tournamentotel"
	"github.com/ctfrancia/maple/business/domain/tournamentbus/stores/tournamentdb"
	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/ctfrancia/maple/foundation/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var build = "develop" // set during build process to be hash of git commit

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

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "MAPLE", traceIDFn, events)

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup failed", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	//-------------- CONFIG --------------
	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:3010"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}
		Auth struct {
			Host string `conf:"default:http://auth-service:6000"`
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
			ServiceName string  `conf:"default:maple"`
			Probability float64 `conf:"default:0.5"`
			// Shouldn't use a high Probability value in non-developer systems.
			// 0.05 should be enough for most systems. Some might want to have
			// this even lower.
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Maple",
		},
	}

	const prefix = "MAPLE"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		return fmt.Errorf("parsing config: %w", err)
	}

	// API startup ===============================================================================
	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}

	log.Info(ctx, "startup", "config", out)

	expvar.NewString("build").Set(cfg.Build)

	// ===========================================================================================
	//  ----------- DATABASE SETUP TODO -----------
	// ===========================================================================================

	log.Info(ctx, "startup", "status", "intitializing database support", "hostport", cfg.DB.Host)

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Info),
	})
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}

	psql, err := db.DB()
	if err != nil {
		return fmt.Errorf("getting database connection: %w", err)
	}

	psql.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	psql.SetMaxOpenConns(cfg.DB.MaxOpenConns)

	defer psql.Close()
	// ===========================================================================================
	//  ----------- BUSINESS  PACKAGES --------------
	// ===========================================================================================

	tournamentOtelExt := tournamentotel.NewExtension()
	tournamentStore := tournamentdb.NewStore(log, db)
	tournamentBus := tournamentbus.NewBusiness(log, tournamentStore, tournamentOtelExt)

	// -------------------------------------------------------------------------
	// Initialize authentication support

	//log.Info(ctx, "startup", "status", "initializing authentication support")

	//authClient := authclient.New(log, cfg.Auth.Host)

	// -------------------------------------------------------------------------
	// Start Tracing Support

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
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()
	// ===========================================================================================
	// Start API Service
	// ===========================================================================================

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		CORSAllowedOrigins: cfg.Web.CORSAllowedOrigins,
		Build:              cfg.Build,
		Tracer:             tracer,
		DB:                 db,
		BusConfig: mux.BusConfig{
			TournamentBus: tournamentBus,
		},
	}

	webAPI := mux.WebAPI(cfgMux)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      webAPI,
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
	// ===========================================================================================
	// SHUTDOWN
	// ===========================================================================================

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
