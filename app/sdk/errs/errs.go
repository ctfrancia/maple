// Package errs provides types and support related to the web error handling.
package errs

import (
	"encoding/json"
	"fmt"
)

// FieldError represents a single validation error for a field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationError represents multiple field validation errors
type ValidationError struct {
	Fields map[string][]string `json:"fields"`
}

func (v ValidationError) Error() string {
	data, _ := json.Marshal(v.Fields)
	return fmt.Sprintf("validation failed: %s", string(data))
}

// NewValidationError creates a new ValidationError
func NewValidationError() *ValidationError {
	return &ValidationError{
		Fields: make(map[string][]string),
	}
}

// Add adds an error message for a field
func (v *ValidationError) Add(field, message string) {
	v.Fields[field] = append(v.Fields[field], message)
}

// HasErrors returns true if there are any validation errors
func (v *ValidationError) HasErrors() bool {
	return len(v.Fields) > 0
}
