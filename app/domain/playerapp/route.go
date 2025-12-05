package playerapp

import "github.com/ctfrancia/maple/business/domain/playerbus"

// Config contains all the required components for the handlers.
type Config struct {
	Log       any
	PlayerBus playerbus.Player
}

func Routes() {
	const version = "v1"
}
