package commands

import (
	"fmt"
	"strings"
)

type RegistrationStatus string

const (
	RegistrationStatusOpen   RegistrationStatus = "open"
	RegistrationStatusClosed RegistrationStatus = "closed"
)

type PaymentType string

const (
	PaymentTypeMonetary PaymentType = "monetary" // money
	PaymentTypePhysical PaymentType = "physical" // e.g. book/lesson/etc.
	PaymentTypeOther    PaymentType = "other"
)

// ValidationError represents multiple field validation errors
type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

func (ve ValidationError) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}

	var messages []string
	for field, msg := range ve.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", field, msg))
	}

	return fmt.Sprintf("validation failed: %s", strings.Join(messages, ", "))
}
