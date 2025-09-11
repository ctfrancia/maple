package types

import (
	"time"
)

type RegistrationStatus string

const (
	RegistrationStatusClosed RegistrationStatus = "closed"
	RegistrationStatusOpen   RegistrationStatus = "open"
)

// Registration represents the registration information for the tournament
type Registration struct {
	Status         *RegistrationStatus `json:"status,omitempty"`
	StartTime      *time.Time          `json:"start_time,omitempty"`
	EndTime        *time.Time          `json:"end_time,omitempty"`
	PublicFee      *int64              `json:"fee,omitempty"`
	PrivateFee     *int64              `json:"private_fee,omitempty"`
	OtherFee       *int64              `json:"other_fee,omitempty"`
	PrizePool      *int64              `json:"prize_pool,omitempty"`
	Payment        *[]Payment          `json:"payment,omitempty"`
	Email          *Email              `json:"email,omitempty"`           // if they register through emailing
	Website        *string             `json:"website,omitempty"`         // if they register through website
	Phone          *string             `json:"phone,omitempty"`           // if they register through phone
	AdditionalInfo *string             `json:"additional_info,omitempty"` // additional info for registration
}
