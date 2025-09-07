package commands

import (
	"github.com/ctfrancia/maple/internal/application/commands/tournament/types"
)

type UpdateTournamentCommand struct {
	Name               *string             `json:"name,omitempty"`
	Description        *string             `json:"description,omitempty"`
	Schedule           *[]types.Schedule   `json:"schedule,omitempty"`
	AdditionalInfo     *string             `json:"additional_info,omitempty"`
	LocationID         *string             `json:"location_id,omitempty"`
	MaxPlayers         *int                `json:"max_players,omitempty"`
	Contact            *types.Contact      `json:"contact,omitempty"`
	OpenToPublic       *bool               `json:"open_to_public,omitempty"`
	OpenToRegistration *bool               `json:"open_to_registration,omitempty"`
	Registration       *types.Registration `json:"registration,omitempty"`
}

// Validate is where we handle the validation of the command
// TODO: all fields need to be scrubbed of swear words, racism, etc.
// also escaped for sql injection/etc.
func (cmd UpdateTournamentCommand) Validate() error {
	errors := make(map[string]string)
	// TODO: expand, to include swears racism, etc.
	// fields also need to be in english
	if cmd.Name != nil && len(*cmd.Name) == 0 {
		errors["name"] = "is required"
	}
	return nil
}
