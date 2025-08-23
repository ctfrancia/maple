package domain

type ConsumerStatus string

const (
	// User Status - business states
	ConsumerStatusActive    ConsumerStatus = "active"
	ConsumerStatusInactive  ConsumerStatus = "inactive"
	ConsumerStatusSuspended ConsumerStatus = "suspended"
	ConsumerStatusPending   ConsumerStatus = "pending"
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
