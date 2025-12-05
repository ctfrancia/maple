// Package pgn represents the PGN (POSet Game Notation) in the system.
package pgn

type PGN struct {
	value string
}

// String returns the value of the pgn.
func (p PGN) String() string {
	return p.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (p PGN) MarshalText() ([]byte, error) {
	return []byte(p.value), nil
}
