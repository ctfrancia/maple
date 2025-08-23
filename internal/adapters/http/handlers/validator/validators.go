// Package validator provides a set of functions for validating incomming requests
package validator

import (
	"regexp"
	"slices"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

// New returns a new pointer to the handler of Validator
func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid should be called after all validations are performed to check if there are indeed valid
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// ReturnErrors returns the errors map mainly for returning to the client
func (v *Validator) ReturnErrors() map[string]string {
	return v.Errors
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In returns true if a specific value is in a list of permitted values.
func (v *Validator) In(key string, permittedValues ...string) bool {
	return PermittedValue(key, permittedValues...)
}

// Generic function which returns true if a specific value is in a list of permitted
// values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches returns true if a string value matches a specific regexp pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Generic function which returns true if all values in a slice are unique.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
