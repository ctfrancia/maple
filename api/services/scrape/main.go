package main

import (
	"context"
	"os"

	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/ctfrancia/maple/foundation/otel"
	"github.com/joho/godotenv"
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

	return nil
}
