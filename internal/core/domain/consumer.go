package domain

const (
	// User Status - business states
	ConsumerStatusActive    ConsumerStatus = "active"
	ConsumerStatusInactive  UserStatus     = "inactive"
	ConsumerStatusSuspended UserStatus     = "suspended"
	ConsumerStatusPending   UserStatus     = "pending"
)
