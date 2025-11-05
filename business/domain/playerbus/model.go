package playerbus

import (
	"net/mail"

	"github.com/google/uuid"
)

// Player represents an individual chess player.
type Player struct {
	ID uuid.UUID
}

// NewPlayer represents what we expect from clients when creating a new player.
type NewPlayer struct {
	FirstName string
	Name      string //name.Name
	Email     mail.Address
	Password  password.Password
}

// UpdatePlayer defines what information may be provided to modify an existting player.
// All fields are optional. Points are used so we may differentiate between data that is not
// provided and data that is explicitly set to a zero value.
type UpdatePlayer struct {
}
