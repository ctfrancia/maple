// Package commands - Represents the user's intent to perform an action
package commands

import (
	"strings"
	"time"
)

// CreateTournamentCommand represents the user's intent to create a tournament
// This represents all fields that are accepted by the API when creating a tournament
type CreateTournamentCommand struct {
	Name               string       `json:"name"`        //`json:"name" validate:"required,gte=3,lte=100"` look into this?
	Description        string       `json:"description"` // optional
	Schedule           []Schedule   `json:"schedule,omitempty"`
	AdditionalInfo     string       `json:"additional_info"`      // optional TODO: add this to the DTO
	LocationID         string       `json:"location_id"`          // need to revisit
	MaxPlayers         int          `json:"max_players"`          // optional when creating
	Contact            Contact      `json:"contact"`              // optional
	OpenToPublic       bool         `json:"open_to_public"`       // optional
	OpenToRegistration bool         `json:"open_to_registration"` // optional
	Registration       Registration `json:"registration"`         // optional
}

// Registration represents the registration information for the tournament
type Registration struct {
	Status     RegistrationStatus `json:"status"`
	StartTime  time.Time          `json:"start_time"`
	EndTime    time.Time          `json:"end_time"`
	PublicFee  int64              `json:"fee"`
	PrivateFee int64              `json:"private_fee"`
	OtherFee   int64              `json:"other_fee"`
	PrizePool  int64              `json:"prize_pool"`
	Payment    []Payment          `json:"payment"`
}

// Schedule represents the schedule for the tournament
type Schedule struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// Payment represents the payment information for the tournament
type Payment struct {
	Place  int         `json:"place"`  // 1st, 2nd, etc.
	Amount int64       `json:"amount"` // if type is monetary
	Type   PaymentType `json:"type"`
}

// Contact represents the contact information for the tournament
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// Validate is where we handle the validation of the command
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
func (cmd CreateTournamentCommand) validateDates(errors map[string]string) {
	// Check if dates are provided but zero (invalid state)
	if len(cmd.Schedule) == 0 {
		return
	}
	// TODO: check if dates are valid

	// organize the dates of start_time and end_time chronologically
}
