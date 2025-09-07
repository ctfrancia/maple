// Package commands - Represents the user's intent to perform an action
package commands

import (
	"strings"

	"github.com/ctfrancia/maple/internal/application/commands/tournament/types"
)

// CreateTournamentCommand represents the user's intent to create a tournament
// This represents all fields that are accepted by the API when creating a tournament
type CreateTournamentCommand struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Schedule       *[]types.Schedule `json:"schedule,omitempty"`
	AdditionalInfo *string           `json:"additional_info,omitempty"` // optional
	Location       *types.Location   `json:"location_id,omitempty"`     // need to revisit
	// MaxPlayers     *int                `json:"max_players,omitempty"`     // optional when creating
	Contact      *types.Contact      `json:"contact,omitempty"`        // optional
	OpenToPublic *bool               `json:"open_to_public,omitempty"` // optional
	Registration *types.Registration `json:"registration,omitempty"`   // optional
}

// Validate is where we handle the validation of the command
// TODO: all fields need to be scrubbed of swear words, racism, etc.
func (cmd CreateTournamentCommand) Validate() error {
	errors := make(map[string]string)

	// Name validation
	if strings.TrimSpace(cmd.Name) == "" {
		errors["name"] = "is required"
	} else if len(strings.TrimSpace(cmd.Name)) < 3 {
		errors["name"] = "must be at least 3 characters"
	} else if len(strings.TrimSpace(cmd.Name)) > 100 {
		errors["name"] = "must be less than 100 characters"
	}

	// Description validation (optional but if provided, check length)
	if len(cmd.Description) > 500 {
		errors["description"] = "must be less than 500 characters"
	}

	// TODO: implemeent schedule validation
	if cmd.Schedule != nil {
		// Schedule validation (optional but if provided, check relationship)
		cmd.validateDates(errors)
	}

	// Date validation (optional but if provided, check relationship)
	cmd.validateDates(errors)

	if len(errors) > 0 {
		return ValidationError{Errors: errors}
	}

	return nil
}

// IsValidationError Helper function to check if an error is a ValidationError
func IsValidationError(err error) (*ValidationError, bool) {
	if ve, ok := err.(ValidationError); ok {
		return &ve, true
	}
	return nil, false
}

// validateDates handles date validation logic
func (cmd *CreateTournamentCommand) validateDates(errors map[string]string) {
	// Check if dates are provided but zero (invalid state)
	if len(*cmd.Schedule) == 0 {
		return
	}
	// TODO: check if dates are valid

	// organize the dates of start_time and end_time chronologically
}
