package types

import (
	"time"
)

// Schedule represents the schedule for the tournament
type Schedule struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
