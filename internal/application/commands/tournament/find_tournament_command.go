package commands

import (
	"github.com/google/uuid"
)

// FindTournamentCommand represents the user's intent to find a tournament
// This represents all fields that are accepted by the API when finding a tournament
type FindTournamentCommand struct {
	ID uuid.UUID `json:"id"` // public uuid
}

// Validate represents multiple field validation errors
// logic for validating the command
func (cmd FindTournamentCommand) Validate() error {
	errors := make(map[string]string)

	if cmd.ID == uuid.Nil {
		errors["id"] = "cannot be nil"
	}

	return nil
}
