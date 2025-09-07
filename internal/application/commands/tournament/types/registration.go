package types

import (
	"time"
)

type RegistrationStatus int

const (
	RegistrationStatusOpen RegistrationStatus = iota
	RegistrationStatusClosed
)

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
