package name

type Name struct {
	value string
}

// String returns the value of the name.
func (n Name) String() string {
	return n.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (n Name) MarshalText() ([]byte, error) {
	return []byte(n.value), nil
}
