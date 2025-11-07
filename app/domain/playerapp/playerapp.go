// Package playerapp maintains the app layer api for the player domain.
package playerapp

import (
	"context"
)

type app struct {
	playerBus playerbus.PlayerBus
}

