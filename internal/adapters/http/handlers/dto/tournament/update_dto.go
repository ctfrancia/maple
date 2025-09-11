package dto

import (
	"strings"

	"github.com/ctfrancia/maple/internal/adapters/http/handlers/validator"
	commands "github.com/ctfrancia/maple/internal/application/commands/tournament"
	//"github.com/ctfrancia/maple/internal/application/commands/tournament/types"
	"github.com/ctfrancia/maple/internal/core/domain"
)

type UpdateTournamentRequest struct {
	Name   *string                  `json:"name,omitempty"`
	Status *domain.TournamentStatus `json:"status,omitempty"`
	// Description            *string                  `json:"description,omitempty"`
	// AdditionalInfo         *string                  `json:"additional_info,omitempty"`
	// Location               *types.Location          `json:"location,omitempty"`
	// Contact                *types.Contact        `json:"contact,omitempty"`
	// Registration           *types.Registration   `json:"registration,omitempty"`
	// PairingMethod          *domain.PairingMethod `json:"pairing_method,omitempty"`
	// Arbitrator             *string               `json:"arbitrator,omitempty"` // first + last name
	// MaxPlayerParticipation *int                  `json:"max_player_participation,omitempty"`
	// MinimumPlayers         *int                  `json:"minimum_players,omitempty"`
	// OpenToPublic           *bool                 `json:"open_to_public,omitempty"`

	// BELOW WILL BE SEPARATE ENDPOINT
	// Schedule       *[]types.Schedule `json:"schedule,omitempty"` // Separate endpoint for this
	// Roster         *types.Roster     `json:"roster,omitempty"` Linked to Matches
	// Matches        *[]types.Match      `json:"matches,omitempty"` Separate endpoint for this
}

func (u UpdateTournamentRequest) Validate() (*commands.UpdateTournamentCommand, error) {
	v := validator.New()
	if u.Name != nil {
		v.Check(strings.TrimSpace(*u.Name) != "", "name", "Tournament name cannot be empty")
	}
	if u.Status != nil {
		v.Check(IsValidTournamentStatus(u.Status), "status", "Tournament status is invalid")
	}

	if !v.Valid() {
		return nil, v
	}
	return mapToCommand(u), nil
}

func IsValidTournamentStatus(s *domain.TournamentStatus) bool {
	status := TournamentStatus(*s)
	switch status {
	case TournamentStatusActive,
		TournamentStatusDraft,
		TournamentStatusDeactive,
		TournamentStatusSuspended,
		TournamentStatusPending,
		TournamentStatusCompleted:
		return true
	default:
		return false
	}
}

func mapToCommand(u UpdateTournamentRequest) *commands.UpdateTournamentCommand {
	return &commands.UpdateTournamentCommand{
		Name:   u.Name,
		Status: u.Status,
	}
}
