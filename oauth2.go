package oauth2

import "time"

////////////////////////////////////////////////////////////////////////////////////////////////////
// Client describes a client's characteristic.
type Client interface {

	// Return client's ID.
	ClientID() string

	// Return client's secret.
	ClientSecret() string

	// Return client's allowed grant types.
	GrantTypes() []string

	// Return client's registered redirect URIs.
	RedirectURIs() []string
}

////////////////////////////////////////////////////////////////////////////////////////////////////
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

////////////////////////////////////////////////////////////////////////////////////////////////////
// Token describes a token's characteristic, it can be either access token or refresh token.
type Token interface {

	// Return client's ID.
	ClientID() string

	// Return user's ID.
	UserID() string

	// Return token.
	Token() string

	// Check if token is expired or not.
	IsExpired() bool

	// Return token's created time.
	CreatedTime() time.Time

	// Return token's expired time.
	ExpiredTime() time.Time
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// TokenStore describes a token store's characteristic.
type TokenStore interface {

	// User
	FindUserWithID(userID string) User
	FindUserWithClient(clientID string, clientSecret string) User
	FindUserWithCredential(username string, password string) User

	// Client
	FindClientWithID(clientID string) Client
	FindClientWithCredential(clientID string, clientSecret string) Client

	// Access Token
	FindAccessToken(token string) Token
	FindAccessTokenWithCredential(clientID string, userID string) Token
	CreateAccessToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token
	DeleteAccessToken(token Token)

	// Refresh Token
	FindRefreshToken(token string) Token
	FindRefreshTokenWithCredential(clientID string, userID string) Token
	CreateRefreshToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token
	DeleteRefreshToken(token Token)

	//	// Authorization code
	//	FindAuthorizationCode(authorizationCode string)
	//	SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// TokenResponse describes a granted response that will be returned to client.
type TokenResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`

	Roles []string `json:"roles,omitempty"`
}
