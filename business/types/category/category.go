// Package category represents a subcategory of a tournament.
// for example, Sub18 or Senior, something that is restricted to age, gender, etc.
// if there is no subcategory, then there are no restrictions.
package category

type Category struct {
	name string
	// TODO: add restrictions
}

func (c Category) String() string {
	return c.name
}

// MarshalText implements the encoding.TextMarshaler interface.
func (c Category) MarshalText() ([]byte, error) {
	return []byte(c.name), nil
}
