package playerdb

import (
	"context"
	"fmt"
	"github.com/ctfrancia/maple/foundation/logger"
	"time"

	"gorm.io/gorm"
)

type Storer interface {
	// Related to scraping
	UpsertPlayers(ctx context.Context, tournamentID string, players []models.Player) error
}

type Store struct {
	log *logger.Logger
	db  *gorm.DB
}

// NewStore contructs the API for data access.
func NewStore(log *logger.Logger, db *gorm.DB) (*Store, error) {
	psql, err := db.DB()
	if err != nil {
		return nil, nil
	}

	psql.SetMaxOpenConns(10)
	psql.SetMaxIdleConns(5)
	psql.SetConnMaxLifetime(5 * time.Minute)

	if err := psql.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &Store{
		log: log,
		db:  db,
	}, nil
}

func (s *Store) UpsertPlayers(ctx context.Context, tournamentID string, players []models.Player) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM players WHERE tournament_id = $1`, tournamentID); err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO players (tournament_id, rank, title, name, fide_id, federation, rating, points, rating_perf)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range players {
		if _, err := stmt.ExecContext(ctx, tournamentID, p.Rank, p.Title, p.Name, p.FideID, p.Federation, p.Rating, p.Points, p.RatingPerf); err != nil {
			return err
		}
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE tournaments SET status='scraped', last_scraped=NOW(), updated_at=NOW() WHERE id=$1`,
		tournamentID); err != nil {
		return err
	}

	return tx.Commit()
}
