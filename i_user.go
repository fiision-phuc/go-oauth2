package oauth2

// IUser descripts an user's characteristic.
type IUser interface {
	GetUserID() string
	GetUsername() string
	GetPassword() string
	GetUserRoles() []string
}
