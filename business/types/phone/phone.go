// package phone contains types related to phone numbers.
package phone

// Phone represents a phone number.
type Phone struct {
	value string
}

// String returns the value of the phone number.
func (p Phone) String() string {
	return p.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (p Phone) MarshalText() ([]byte, error) {
	return []byte(p.value), nil
}
