// Package country contains types related to countries.
package country

// Country represents a country.
type Country struct {
	value string
}

// String returns the value of the country.
func (c Country) String() string {
	return c.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (c Country) MarshalText() ([]byte, error) {
	return []byte(c.value), nil
}
