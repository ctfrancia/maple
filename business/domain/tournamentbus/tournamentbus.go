// Package tournament provides the business access to the tournament domain.
package tournamentbus

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ctfrancia/maple/foundation/logger"

	"github.com/ctfrancia/maple/foundation/otel"
)

// Storer is the interface for the persistence layer.
type Storer interface {
}

// Extension wraps additional business logic around an existing one.
type Extension func() any

// Business manages the set of APIs for user access
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

// Create creates a new tournament.
func (b *Business) Create(ctx context.Context, nt NewTournament) error {
	ctx, span := otel.AddSpan(ctx, "business.tournamentbus.Create")
	defer span.End()

	now := time.Now()

	t := Tournament{
		ID:          uuid.New(),
		Description: nt.Description,
		Name:        nt.Name,
		CreatedBy:   nt.CreatedBy,
		PlayerID:    nt.PlayerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return nil
}
