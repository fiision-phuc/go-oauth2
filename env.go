package oauth2

import "os"

// Define keywords for OAuth2.0 flows.
const (
	AuthorizationCodeGrant = "authorization_code" // For apps running on a web server
	ClientCredentialsGrant = "client_credentials" // For application access
	PasswordGrant          = "password"           // For logging in with a username and password
	RefreshTokenGrant      = "refresh_token"      // Should allow refresh token or not
	//	ImplicitGrant          = "implicit"           // For browser-based or mobile apps
)

// Define table names.
const (
	TableAccessToken  = "access_token"
	TableRefreshToken = "refresh_token"
	TableClient       = "client"
	TableUser         = "user"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////////////////////////

// GetEnv retrieves value from environment.
func GetEnv(key string) string {
	if len(key) == 0 {
		return ""
	}
	return os.Getenv(key)
}

// SetEnv persists key-value to environment.
func SetEnv(key string, value string) {
	if len(key) == 0 || len(value) == 0 {
		return
	}
	os.Setenv(key, value)
}
