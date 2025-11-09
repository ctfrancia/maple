package round

import (
	"time"

	"github.com/ctfrancia/maple/business/types/pgn"
)

type Number struct {
	value int
}

type Round struct {
	Num     Number
	Date    time.Time
	Players [2]player.Player // 0: white, 1: black
	PEN     pgn.PGN
}
