package oauth2

// User describes an user's characteristic.
type User interface {

	// Return user's ID.
	UserID() string

	// Return user's username.
	Username() string

	// Return user's password.
	Password() string

	// Return user's roles.
	UserRoles() []string
}
