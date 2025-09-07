// Package types defines the types used for the tournament command
package types

// Contact represents the contact information for the tournament
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
