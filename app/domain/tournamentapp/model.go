package tournamentapp

import (
	"context"
	"errors"

	"github.com/ctfrancia/maple/app/sdk/mid"
	"github.com/ctfrancia/maple/business/domain/tournamentbus"
)

type Tournament struct {
}

// NewTournament represents the initial information needed to create and build a new tournament.
type NewTournament struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
}

func toBusNewTournament(ctx context.Context, nt NewTournament) (tournamentbus.NewTournament, error) {
	errs := make(map[string]string)
	ID, err := mid.GetUserID(ctx)
	if err != nil {
		errs["creator_id"] = err.Error()
	}

	if len(errs) != 0 {
		return tournamentbus.NewTournament{}, errors.New("implement me")
	}

	return tournamentbus.NewTournament{
		CreatedBy: ID,
	}, nil
}
