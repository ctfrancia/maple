package domain

const (
	// User Status - business states
	ConsumerStatusActive    UserStatus = "active"
	ConsumerStatusInactive  UserStatus = "inactive"
	ConsumerStatusSuspended UserStatus = "suspended"
	ConsumerStatusPending   UserStatus = "pending"
)
