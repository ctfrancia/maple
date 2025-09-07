package commands

import (
	"github.com/ctfrancia/maple/internal/application/commands/tournament/types"
)

type UpdateTournamentCommand struct {
	Name           *string             `json:"name,omitempty"`
	Description    *string             `json:"description,omitempty"`
	AdditionalInfo *string             `json:"additional_info,omitempty"`
	Schedule       *[]types.Schedule   `json:"schedule,omitempty"`
	Location       *types.Location     `json:"location,omitempty"`
	Roster         *types.Roster       `json:"roster,omitempty"`
	Contact        *types.Contact      `json:"contact,omitempty"`
	Registration   *types.Registration `json:"registration,omitempty"`
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

	if cmd.Description != nil && len(*cmd.Description) == 0 {
		errors["description"] = "a description is required"
	}

	if cmd.Schedule != nil && len(*cmd.Schedule) == 0 {
		errors["schedule"] = "a schedule cannot be empty"
	}

	if len(errors) > 0 {
		return ValidationError{Errors: errors}
	}

	return nil
}
