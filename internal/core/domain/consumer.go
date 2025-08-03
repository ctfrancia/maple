package domain

const (
	// User Status - business states
	ConsumerStatusActive    ConsumerStatus = "active"
	ConsumerStatusInactive  UserStatus     = "inactive"
	ConsumerStatusSuspended UserStatus     = "suspended"
	ConsumerStatusPending   UserStatus     = "pending"
)

type NewAPIConsumer struct {
	PublicID        string
	FirstName       string
	LastName        string
	Username        string
	Email           string
	Password        string
	Website         string
	ClubAffiliation string
}
