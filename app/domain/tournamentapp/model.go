package tournamentapp

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ctfrancia/maple/app/sdk/errs"
	"github.com/ctfrancia/maple/app/sdk/mid"
	"github.com/ctfrancia/maple/business/domain/tournamentbus"
	"github.com/ctfrancia/maple/business/types/description"
	"github.com/ctfrancia/maple/business/types/name"
)

type Tournament struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Poster      string `json:"poster"`
	CreatedBy   string `json:"created_by"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

// NewTournament represents the initial information needed to create and build a new tournament.
type NewTournament struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
}

func toBusNewTournament(ctx context.Context, nt NewTournament) (tournamentbus.NewTournament, error) {
	validationErr := errs.NewValidationError()
	ID, err := mid.GetUserID(ctx)
	if err != nil {
		validationErr.Add("created_by", "cannot be empty")
	}

	name, err := name.Parse(strings.TrimSpace(nt.Name))
	if err != nil {
		validationErr.Add("name", "cannot be empty")
	}

	desc, err := description.Parse(strings.TrimSpace(nt.Description))
	if err != nil {
		validationErr.Add("description", "cannot be empty")
	}

	if validationErr.HasErrors() {
		return tournamentbus.NewTournament{}, validationErr
	}

	return tournamentbus.NewTournament{
		Name:        name,
		CreatedBy:   ID,
		Description: desc,
	}, nil
}

// ============================================================================
func toAppTournament(t tournamentbus.Tournament) Tournament {
	return Tournament{
		ID:          t.ID.String(),
		Name:        t.Name.String(),
		Description: t.Description.String(),
		Poster:      t.Poster.String(),
		CreatedBy:   t.CreatedBy.String(),
		DateCreated: t.DateCreated.Format(time.RFC3339),
		DateUpdated: t.DateUpdated.Format(time.RFC3339),
	}
}
