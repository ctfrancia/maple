package playerbus

import (
	"net/mail"

	"github.com/ctfrancia/maple/business/types/password"
	"github.com/ctfrancia/maple/business/types/username"
	"github.com/google/uuid"
)

// Player represents information about an individual player.
type Player struct {
	ID       uuid.UUID
	Name     name.Name
	location Location
}

// NewPlayer represents what we expect from clients when creating a new player.
type NewPlayer struct {
	Username        username.Username
	Email           mail.Address
	Password        password.Password
	ConfirmPassword password.Password
}

// UpdatePlayer defines what information may be provided to modify an existting player.
// All fields are optional. Points are used so we may differentiate between data that is not
// provided and data that is explicitly set to a zero value.
type UpdatePlayer struct {
	FirstName *string
	LastName  *string
	Location
}
