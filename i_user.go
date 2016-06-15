package oauth2

// IUser describes an user's characteristic.
type IUser interface {

	// Return user's ID.
	UserID() string

	// Return user's username.
	Username() string

	// Return user's password.
	Password() string

	// Return user's roles.
	UserRoles() []string
}
