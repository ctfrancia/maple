// Package zip contains types related to zip codes.
package zip

// Zip represents a zip code.
type Zip struct {
	value string
}

// String returns the value of the zip code.
func (z Zip) String() string {
	return z.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (z Zip) MarshalText() ([]byte, error) {
	return []byte(z.value), nil
}
