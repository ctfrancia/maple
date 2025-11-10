// Package city contains types related to cities.
package city

// City represents a city.
type City struct {
	value string
}

// String returns the value of the city.
func (c City) String() string {
	return c.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (c City) MarshalText() ([]byte, error) {
	return []byte(c.value), nil
}
