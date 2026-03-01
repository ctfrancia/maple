package playerdb

import (
	"context"
	"fmt"
	"time"

	"github.com/ctfrancia/maple/business/domain/playerbus"
	"github.com/ctfrancia/maple/foundation/logger"

	"gorm.io/gorm"
)

type Storer interface {
	// Related to scraping
	UpsertPlayers(ctx context.Context, tournamentID string, players []playerbus.Player) error
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

func (s *Store) UpsertPlayers(ctx context.Context, tournamentID string, players []playerbus.Player) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("tournament_id = ?", tournamentID).Delete(&playerbus.Player{}).Error; err != nil {
			return err
		}

		if len(players) > 0 {
			// Set tournament_id on each player
			for i := range players {
				players[i].TournamentID = tournamentID
			}

			if err := tx.Create(&players).Error; err != nil {
				return err
			}
		}

		return tx.Model(&Tournament{}).
			Where("id = ?", tournamentID).
			Updates(map[string]any{
				"status":       "scraped",
				"last_scraped": gorm.Expr("NOW()"),
				"updated_at":   gorm.Expr("NOW()"),
			}).Error
	})
}
