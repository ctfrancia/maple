// Package tournament provides the business access to the tournament domain.
package tournamentbus

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ctfrancia/maple/foundation/logger"
	"github.com/ctfrancia/maple/foundation/otel"
)

// Storer is the interface for the persistence layer.
type Storer interface {
	Create(ctx context.Context, t Tournament) error
	Update(ctx context.Context, t Tournament, ut UpdateTournament) error
	Delete(ctx context.Context, t Tournament) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Tournament, error)
	QueryByID(ctx context.Context, tID uuid.UUID) (Tournament, error)
}

// Extension wraps additional business logic around an existing one.
type Extension func() any

// Business manages the set of APIs for Tournament access
type Business struct {
	log    *logger.Logger
	storer Storer
	//userBus  userbus.Business
	//delegate any
}

// NewBusiness constructs a tournament business API for use.
func NewBusiness(log *logger.Logger, storer Storer) *Business {
	b := &Business{
		log:    log,
		storer: storer,
	}

	return b
}

// Create adds a new tournament to the system.
func (b *Business) Create(ctx context.Context, nt NewTournament) (Tournament, error) {
	ctx, span := otel.AddSpan(ctx, "business.tournamentbus.Create")
	defer span.End()

	now := time.Now()
	t := Tournament{
		ID:          uuid.New(),
		Name:        nt.Name,
		Description: nt.Description,
		CreatedBy:   nt.CreatedBy,
		PlayerID:    nt.PlayerID,
		Enabled:     false,
		Status:      StatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := b.storer.Create(ctx, t); err != nil {
		return Tournament{}, fmt.Errorf("create: %w", err)
	}

	return t, nil
}

// Update updates an existing tournament.
func (b *Business) Update(ctx context.Context, t Tournament, ut UpdateTournament) (Tournament, error) {
	ctx, span := otel.AddSpan(ctx, "business.tournamentbus.Update")
	defer span.End()

	if err := b.storer.Update(ctx, t, ut); err != nil {
		return Tournament{}, fmt.Errorf("update: %w", err)
	}

	return Tournament{}, nil
}

// Delete removes a tournament from the system.
func (b *Business) Delete(ctx context.Context, t Tournament) error {
	ctx, span := otel.AddSpan(ctx, "business.tournamentbus.Delete")
	defer span.End()

	return b.storer.Delete(ctx, t)
}

// Query retrieves a list of tournaments.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Tournament, error) {
	tournaments, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return tournaments, nil
}

// QueryByID retrieves a Tournament by ID.
func (b *Business) QueryByID(ctx context.Context, tID uuid.UUID) (Tournament, error) {
	tournament, err := b.storer.QueryByID(ctx, tID)
	if err != nil {
		return Tournament{}, fmt.Errorf("query: tournamentID[%s]: %w", tID, err)
	}

	return tournament, nil
}
