package types

import "github.com/google/uuid"

type Roster struct {
	count        int
	registered   []uuid.UUID // Player
	unregistered []Email
}

func (p *Roster) Validate() error {
	return nil
}

func (p *Roster) SetPlayerCount(count int) {
	p.count = len(p.registered) + len(p.unregistered)
}
