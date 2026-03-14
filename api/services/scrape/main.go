package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf/v3"
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
		log.Error(ctx, "Error running scraper: %v", err)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	env := os.Getenv("SCRAPER_ENV")
	if env == "" || env == "development" {
		if err := godotenv.Load("zoltan/compose/.env.dev"); err != nil {
			log.Info(ctx, "startup", "warn", "Warning: .env.dev file not found")
		}
	}

	//-------------- CONFIG --------------
	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:8080"`
			DebugHost          string        `conf:"default:0.0.0.0:8081"`
			CORSAllowedOrigins []string      `conf:"default:*"`
			UserAgent          string        `conf:"default:Mozilla/5.0 (compatible; MapleBot/1.0; +https://maple.chess)"`
			ReqDelay           time.Duration `conf:"default:2s"`
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
			Desc:  "Scraper for Maple",
		},
	}

	const prefix = "SCRAPER"
	help, err := conf.Parse(prefix, cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		return fmt.Errorf("parsing config: %w", err)
	}

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}

	log.Info(ctx, "startup", "config", out)

	expvar.NewString("build").Set(cfg.Build)

	// -------------------------------------------------------------------------
	// DB startup ==============================================================
	// -------------------------------------------------------------------------
	log.Info(ctx, "startup", "status", "intitializing database support", "hostport", cfg.DB.Host)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Info),
	})
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}

	/*
		err = tournamentdb.CreateMigration(db)
		if err != nil {
			return fmt.Errorf("creating migration: %w", err)
		}
	*/

	psql, err := db.DB()
	if err != nil {
		return fmt.Errorf("getting database connection: %w", err)
	}

	psql.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	psql.SetMaxOpenConns(cfg.DB.MaxOpenConns)

	defer psql.Close()

	return nil
}
