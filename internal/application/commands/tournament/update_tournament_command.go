package commands

import (
	"github.com/ctfrancia/maple/internal/application/commands/tournament/types"
	"github.com/ctfrancia/maple/internal/core/domain"
)

type UpdateTournamentCommand struct {
	Name                   *string                  `json:"name,omitempty"`
	Status                 *domain.TournamentStatus `json:"status,omitempty"`
	Description            *string                  `json:"description,omitempty"`
	AdditionalInfo         *string                  `json:"additional_info,omitempty"`
	Location               *types.Location          `json:"location,omitempty"`
	Contact                *types.Contact           `json:"contact,omitempty"`
	Registration           *types.Registration      `json:"registration,omitempty"`
	PairingMethod          *domain.PairingMethod    `json:"pairing_method,omitempty"`
	Arbitrator             *string                  `json:"arbitrator,omitempty"` // first + last name
	MaxPlayerParticipation *int                     `json:"max_player_participation,omitempty"`
	MinimumPlayers         *int                     `json:"minimum_players,omitempty"`
	OpenToPublic           *bool                    `json:"open_to_public,omitempty"`
	// Schedule       *[]types.Schedule `json:"schedule,omitempty"` // Separate endpoint for this
	// Roster         *types.Roster     `json:"roster,omitempty"` Linked to Matches
	// Matches        *[]types.Match      `json:"matches,omitempty"` Separate endpoint for this
}

/*
type Tournament struct {
	PublicID           uuid.UUID // this is the public ID
	Name               string
	Location           Location // address or location
	Creator            APIConsumer
	Contact            Contact
	Description        string
	OpenToPublic       bool
	OpenToSpectators   bool
	OpenToRegistration bool
	Registration       Registration
	Arbitrator         string
	PairingMethod      PairingMethod
	Matches            []Match // no more no less than 2 white/black
	Players            []Player
	NumberOfPlayers    int // how many are participating
	Schedule           []Schedule
	Results            []Result
	Status             TournamentStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
	SoftDeletedAt      time.Time
	DeletedAt          time.Time
}
*/

// Validate is where we handle the validation of the command
// TODO: all fields need to be scrubbed of swear words, racism, etc.
// also escaped for sql injection/etc.
func (cmd *UpdateTournamentCommand) Validate() error {
	errors := make(map[string]string)
	// TODO: expand, to include swears racism, etc.
	// fields also need to be in english
	if cmd.Name != nil && len(*cmd.Name) == 0 {
		errors["name"] = "cannot be empty to update"
	}

	if cmd.Description != nil && len(*cmd.Description) == 0 {
		errors["description"] = "a description is required"
	}

	if len(errors) > 0 {
		return ValidationError{Errors: errors}
	}

	return nil
}

// ToUpdateMap converts the command to a map of updates for partial updates in DB
func (u *UpdateTournamentCommand) ToUpdateMap() map[string]interface{} {
	updates := make(map[string]any)

	if u.Name != nil {
		updates["name"] = *u.Name
	}
	if u.Description != nil {
		updates["description"] = *u.Description
	}

	return updates
}
