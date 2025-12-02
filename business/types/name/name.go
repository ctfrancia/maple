package name

import (
	"fmt"
	"regexp"
)

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

// nameRegEx is the regular expression used to validate a name.
var nameRegEx = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9' -]{2,19}$")

// Parse parses the string value and returns a name if the value complies with the rules for a name.
func Parse(s string) (Name, error) {
	if !nameRegEx.MatchString(s) {
		return Name{}, fmt.Errorf("invalid name: %s", s)
	}

	return Name{s}, nil
}

// MustParse parses the string value and returns a name if the value complies with the rules for a name.
// It panics if the value is invalid.
func MustParse(s string) Name {
	n, err := Parse(s)
	if err != nil {
		panic(err)
	}

	return n
}
