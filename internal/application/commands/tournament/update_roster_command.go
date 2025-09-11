package commands

import (
	"github.com/ctfrancia/maple/internal/application/commands/tournament/types"
)

type Roster struct {
	count        int16
	registered   []types.Player // Player
	unregistered []string       // "chess_title first_name last_name"
}

func (p *Roster) Validate() error {
	return nil
}

func (p *Roster) SetPlayerCount(count int) {
	p.count = int16(len(p.registered) + len(p.unregistered))
}
