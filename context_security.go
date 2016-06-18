package oauth2

// Security describes a user's security scope.
type Security struct {

	// Registered user. Always available.
	User IUser
	// Registered client. Always available.
	Client IClient
	// Access token that had been given to user. Always available.
	AccessToken IToken
	// Refresh token that had been given to user. Might not be available all the time.
	RefreshToken IToken
}
