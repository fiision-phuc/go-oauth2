package oauthRole

// User's roles.
const (
	Admin   = "r_admin"
	Android = "r_android"
	IOS     = "r_ios"
	Manager = "r_manager"
	User    = "r_user"
	Web     = "r_web"
	Windows = "r_windows"
)

// All returns an array that contains all roles.
func All() []string {
	a := []string{Admin, Android, IOS, Manager, User, Web, Windows}
	return a
}

// AllDevices returns an array that contains all device roles.
func AllDevices() []string {
	return []string{Admin, Android, IOS, Web, Windows}
}

// AllDevices returns an array that contains all user roles.
func AllUsers() []string {
	return []string{Admin, Manager, User}
}
