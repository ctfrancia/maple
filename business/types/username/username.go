package username

// Username represents a username.
type Username struct {
	value string
}

// String returns the username as a string.
func (u Username) String() string {
	return u.value
}
