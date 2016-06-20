package oauth2

// Define OAuth2 flows.
const (
	AuthorizationCodeGrant = "authorization_code" // For apps running on a web server
	ClientCredentialsGrant = "client_credentials" // For application access
	ImplicitGrant          = "implicit"           // For browser-based or mobile apps
	PasswordGrant          = "password"           // For logging in with a username and password
	RefreshTokenGrant      = "refresh_token"      // Should allow refresh token or not

)

// Define OAuth2 tables.
const (
	TableRefreshToken = "oauth_refresh_token"
	TableAccessToken  = "oauth_access_token"
	TableClient       = "oauth_client"
	TableUser         = "oauth_user"
)

// TokenResponse describes a granted response that will be returned to client.
type TokenResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`

	Roles []string `json:"roles,omitempty"`
}
