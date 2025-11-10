// Package country contains types related to countries.
package state

// State represents a state.
type State struct {
	value string
}

// String returns the value of the state.
func (s State) String() string {
	return s.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (s State) MarshalText() ([]byte, error) {
	return []byte(s.value), nil
}
