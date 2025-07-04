package domain

// import "time"

type (
	EntityType   string
	Permission   string
	UserStatus   string
	ConsumerTier string
	AuthMethod   string
)

const (
	EntityTypeUser     EntityType = "user"
	EntityTypeConsumer EntityType = "consumer"
)

const (
	PermissionRead  Permission = "read"
	PermissionWrite Permission = "write"
	PermissionAdmin Permission = "admin"
)

// User represents a human user
type User struct {
	ID          string
	Username    string
	Email       string
	Permissions []Permission
}

// Consumer represents an API consumer (service, application, etc.)
type Consumer struct {
	ID          string
	Name        string
	APIKey      string
	Permissions []Permission
}

/*
might need to add this back in later
// AuthResult is the Authentication result
type AuthResult struct {
	Entity      AuthenticatedEntity
	Token       string
	ExpiresAt   time.Time
	Permissions []Permission
}
*/

const (
	// Authentication Methods - business ways to authenticate
	AuthMethodPassword AuthMethod = "password"
	AuthMethodAPIKey   AuthMethod = "api_key"
	AuthMethodOAuth    AuthMethod = "oauth"
	AuthMethodSSO      AuthMethod = "sso"
)

const (
	// Token expiration rules
	UserTokenExpiryHours     = 24
	ConsumerTokenExpiryHours = 168 // 7 days
	AdminTokenExpiryHours    = 8   // Shorter for security

	// Rate limiting rules
	BasicTierRequestsPerHour      = 1000
	PremiumTierRequestsPerHour    = 10000
	EnterpriseTierRequestsPerHour = 100000

	// Validation rules
	MinPasswordLength = 8
	MaxPasswordLength = 128
	MinAPIKeyLength   = 32
	MaxAPIKeyLength   = 64

	// Business constraints
	MaxPermissionsPerUser     = 10
	MaxPermissionsPerConsumer = 5
	MaxConsumersPerUser       = 3
)
