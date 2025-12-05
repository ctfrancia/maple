// Package federation contains types related to federations.
package federation

// Federation represents a federation.
type Federation struct {
	value   string // has to be a valid FIDE federation ID
	isValid bool   // is valid for value's ELO?
	remarks string // additional remarks
}

// String returns the value of the federation.
func (f Federation) String() string {
	return f.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (f Federation) MarshalText() ([]byte, error) {
	return []byte(f.value), nil
}
