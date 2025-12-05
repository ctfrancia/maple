// package poster contains types related to the poster of a tournament.
package poster

// Poster represents the poster of a tournament.
type Poster struct {
	value string
}

// String returns the value of the poster.
func (p Poster) String() string {
	return p.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (p Poster) MarshalText() ([]byte, error) {
	return []byte(p.value), nil
}

