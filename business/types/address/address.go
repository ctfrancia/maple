// Package address contains types related to physical addresses.
package address

// Address represents a physical address.
type Address struct {
	value string
}

// String returns the value of the address.
func (a Address) String() string {
	return a.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.value), nil
}
