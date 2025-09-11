// Package types defines the types used for the tournament command
package types

// Contact represents the contact information for the tournament
type Contact struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Phone *string `json:"phone,omitempty"`
}
