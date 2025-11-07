package playerapp

import (
	"encoding/json"

	"github.com/ctfrancia/maple/business/domain/playerbus"
)

// Player represents information about an individual player.
type Player struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Club      string `json:"club"`
	Biography string `json:"biography"`
}

// Encode implements the encoder interface.
func (p Player) Encode() ([]byte, string, error) {
	data, err := json.Marshal(p)

	return data, "application/json", err
}

func toAppPlayer(bus playerbus.Player) Player {
	return Player{
		ID:        bus.ID.String(),
		FirstName: bus.FirstName.String(),
		LastName:  bus.LastName.String(),
		Email:     bus.Email.Address,
	}
}
