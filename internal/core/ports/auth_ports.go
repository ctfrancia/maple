package ports

import (
	"context"

	"github.com/ctfrancia/maple/internal/core/domain"
)

// AuthenticationServicer defines the authentication use case
type AuthenticationServicer interface {
	// Login authenticates with username/password (typically for users)
	Login(ctx context.Context, username, password string) (*domain.AuthResult, error)

	// AuthenticateAPIKey authenticates with API key (typically for consumers)
	AuthenticateAPIKey(ctx context.Context, apiKey string) (*domain.AuthResult, error)

	// ValidateToken validates and parses a JWT token
	ValidateToken(ctx context.Context, token string) (*domain.AuthResult, error)

	// HasPermission checks if an entity has a specific permission
	HasPermission(entity domain.AuthenticatedEntity, permission domain.Permission) bool
}

// Repository interfaces
type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	ValidatePassword(ctx context.Context, username, password string) error
}

type ConsumerRepository interface {
	GetByAPIKey(ctx context.Context, apiKey string) (*domain.Consumer, error)
	GetByID(ctx context.Context, id string) (*domain.Consumer, error)
}

type TokenRepository interface {
	GenerateToken(ctx context.Context, entity domain.AuthenticatedEntity) (string, time.Time, error)
	ValidateToken(ctx context.Context, token string) (*domain.AuthResult, error)
}
