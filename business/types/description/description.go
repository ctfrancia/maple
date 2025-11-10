package description

// Description represents the description string
type Description struct {
	value string
}

// String returns the value of the description.
func (d Description) String() string {
	return d.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (d Description) MarshalText() ([]byte, error) {
	return []byte(d.value), nil
}
