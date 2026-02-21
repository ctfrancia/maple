// Package tournamentdb contains tournament CRUD functinality.
package tournamentdb

import (
	"context"
	"time"

	"github.com/ctfrancia/maple/business/domain/tournamentbus"
	tb "github.com/ctfrancia/maple/business/domain/tournamentbus"
	"github.com/ctfrancia/maple/business/sdk/order"
	"github.com/ctfrancia/maple/business/sdk/page"
	"github.com/ctfrancia/maple/foundation/logger"

	"github.com/google/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
	QueryByID(ctx context.Context, tID uuid.UUID) (Tournament, error)

*/

type Storer interface {
	NewWithTx(tx *gorm.DB) (Storer, error)
	Create(ctx context.Context, t tb.Tournament) error
	Update(ctx context.Context, t tb.Tournament, ut tb.UpdateTournament) error
	Delete(ctx context.Context, t tb.Tournament) error
	// Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]tb.Tournament, error)
	// Count(ctx context.Context, filter QueryFilter) (int, error)
	// QueryByID(ctx context.Context, userID uuid.UUID) (tb.Tournament, error)
	// QueryByEmail(ctx context.Context, email mail.Address) (tb.Tournament, error)

	// Used for scraping
	KnownTournamentIDs(ctx context.Context, federation string) (map[uuid.UUID]bool, error)
	UpsertTournament(ctx context.Context, t tb.Tournament) (isnew bool, err error)
	UpdateStatus(ctx context.Context, tournamentID uuid.UUID, status string) error
	TournamentsToScrape(ctx context.Context, maxAge time.Duration) ([]tb.Tournament, error)
}

// Store manages the set of APIs for the tournament database access.
type Store struct {
	log *logger.Logger
	db  *gorm.DB
}

// NewStore contructs the API for data access.
func NewStore(log *logger.Logger, db *gorm.DB) *Store {
	psql, err := db.DB()
	if err != nil {
		return nil
	}

	psql.SetMaxOpenConns(10)
	psql.SetMaxIdleConns(5)
	psql.SetConnMaxLifetime(5 * time.Minute)

	return &Store{
		log: log,
		db:  db,
	}
}

// NewWithTx creates a new Store using the provided transaction.
func (s *Store) NewWithTx(tx *gorm.DB) (Storer, error) {
	return &Store{
		log: s.log,
		db:  tx,
	}, nil
}

// Create creates a new tournament.
func (s *Store) Create(ctx context.Context, t tb.Tournament) error {
	return s.db.WithContext(ctx).Create(&t).Error
}

// Update updates an existing tournament with pessimistic locking.
func (s *Store) Update(ctx context.Context, t tb.Tournament, ut tb.UpdateTournament) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing tb.Tournament
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&existing, t.ID).Error; err != nil {
			return err
		}

		// Now update with the new values
		return tx.Save(&t).Error
	})
}

// Delete deletes an existing tournament.
func (s *Store) Delete(ctx context.Context, t tb.Tournament) error {
	return s.db.WithContext(ctx).Delete(&t).Error
}

func (s *Store) Query(ctx context.Context, filter tournamentbus.QueryFilter, orderBy order.By, page page.Page) ([]tb.Tournament, error) {
	return nil, nil
}

func (s *Store) QueryByID(ctx context.Context, tID uuid.UUID) (tb.Tournament, error) {
	return tb.Tournament{}, nil
}

// ----------------- SCRAPING METHODS-----------------
func (s *Store) KnownTournamentIDs(ctx context.Context, federation string) (map[uuid.UUID]bool, error) {
	return nil, nil
}

func (s *Store) UpsertTournament(ctx context.Context, t tb.Tournament) (isnew bool, err error) {
	return false, nil
}

func (s *Store) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return nil
}

func (s *Store) TournamentsToScrape(_ context.Context, maxAge time.Duration) ([]tb.Tournament, error) {
	return nil, nil
}
