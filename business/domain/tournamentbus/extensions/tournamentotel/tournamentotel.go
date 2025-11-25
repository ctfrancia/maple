package tournamentotel

import (
	"context"

	"github.com/ctfrancia/maple/business/domain/tournamentbus"
	"github.com/ctfrancia/maple/business/sdk/order"
	"github.com/ctfrancia/maple/business/sdk/page"
	"github.com/ctfrancia/maple/foundation/otel"
)

// Extension provides a wrapper for otel functionality around the tournamentbus.
type Extension struct {
	bus tournamentbus.ExtBusiness
}

// NewExtension constructs a new extension that wraps the auditbus with otel.
func NewExtension() tournamentbus.Extension {
	return func(bus tournamentbus.ExtBusiness) tournamentbus.ExtBusiness {
		return &Extension{
			bus: bus,
		}
	}
}

func (e *Extension) Create(ctx context.Context, t tournamentbus.NewTournament) (tournamentbus.Tournament, error) {
	ctx, span := otel.AddSpan(ctx, "business.tournamentbus.Create")
	defer span.End()

	return e.bus.Create(ctx, t)
}

func (e *Extension) Update(ctx context.Context, t tournamentbus.Tournament, ut tournamentbus.UpdateTournament) error {
	ctx, span := otel.AddSpan(ctx, "business.tournamentbus.Update")
	defer span.End()
	return e.bus.Update(ctx, t, ut)
}

func (e *Extension) Query(ctx context.Context, filter tournamentbus.QueryFilter, orderBy order.By, page page.Page) ([]tournamentbus.Tournament, error) {
	ctx, span := otel.AddSpan(ctx, "business.tournamentbus.Query")
	defer span.End()
	return e.bus.Query(ctx, filter, orderBy, page)
}

func (e *Extension) Delete(ctx context.Context, t tournamentbus.Tournament) error {
	ctx, span := otel.AddSpan(ctx, "business.tournamentbus.Delete")
	defer span.End()

	return e.bus.Delete(ctx, t)
}
