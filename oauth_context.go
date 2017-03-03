package oauth2

// OAuthContext describes a user's oauth scope.
type OAuthContext struct {

	// Registered user. Always available.
	User User
	// Registered client. Always available.
	Client Client
	// Access token that had been given to user. Always available.
	AccessToken Token
	// Refresh token that had been given to user. Might not be available all the time.
	RefreshToken Token
}
