package password

import (
	"fmt"
	"regexp"
)

// Password represents a password.
type Password struct {
	value string
}

// Strinng returns the password as a string.
func (p Password) String() string {
	return p.value
}

// Equal provides a means to compare two passwords, used for testing.
func (p Password) Equal(other Password) bool {
	return p.value == other.value
}

// MarshallText supports for logging and marshalling.
func (p Password) MarshalText() ([]byte, error) {
	return []byte(p.value), nil
}

var passwordRegexp = regexp.MustCompile("^(?=.*[!@#$%^&*])[a-zA-Z0-9!@#$%^&*]{8,16}$")

// Parse parses a password and returns a string if the value passes the password validation.
func Parse(value string) (Password, error) {
	if !passwordRegexp.MatchString(value) {
		return Password{}, fmt.Errorf("invalid password %q", value)
	}

	return Password{value}, nil
}
