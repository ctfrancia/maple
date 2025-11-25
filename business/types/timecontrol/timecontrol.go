// Package timecontrol contains types related to time controls.
package timecontrol

import "fmt"

// TimeControl represents a time control.
type TimeControl struct {
	value   string
	remarks string
}

// String returns the value of the time control.
func (t TimeControl) String() string {
	s, _ := t.MarshalText()
	return string(s)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (t TimeControl) MarshalText() ([]byte, error) {
	return fmt.Appendf(nil, "%s (%s)", t.value, t.remarks), nil
}
