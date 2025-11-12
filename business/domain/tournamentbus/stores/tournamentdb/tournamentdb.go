// Package tournamentdb contains tournament CRUD functinality.
package tournamentdb

import (
	"context"

	tb "github.com/ctfrancia/maple/business/domain/tournamentbus"
	"github.com/ctfrancia/maple/foundation/logger"

	"gorm.io/gorm"
)

type Storer interface {
	NewWithTx(tx *gorm.DB) (Storer, error)
	Create(ctx context.Context, t tb.Tournament) error
	Update(ctx context.Context, t tb.Tournament) error
	Delete(ctx context.Context, t tb.Tournament) error
	// Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]tb.Tournament, error)
	// Count(ctx context.Context, filter QueryFilter) (int, error)
	// QueryByID(ctx context.Context, userID uuid.UUID) (tb.Tournament, error)
	// QueryByEmail(ctx context.Context, email mail.Address) (tb.Tournament, error)
}

// Store manages the set of APIs for the tournament database access.
type Store struct {
	log *logger.Logger
	db  *gorm.DB
}

// NewStore contructs the API for data access.
func NewStore(log *logger.Logger, db *gorm.DB) *Store {
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

// Update updates an existing tournament.
// NOTE: here in the future then we will need to do pessimistic locking.
// https://gorm.io/docs/transactions.html#Optimistic-Locking
func (s *Store) Update(ctx context.Context, t tb.Tournament) error {
	// TODO: Implement this - correctly.
	return s.db.WithContext(ctx).Save(&t).Error
}

// Delete deletes an existing tournament.
func (s *Store) Delete(ctx context.Context, t tb.Tournament) error {
	return s.db.WithContext(ctx).Delete(ctx).Delete(&t).Error
}
