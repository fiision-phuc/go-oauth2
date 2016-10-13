package oauth2

import "time"

////////////////////////////////////////////////////////////////////////////////////////////////////
// IClient describes a client's characteristic.
type IClient interface {

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
// IStore describes a token storage's characteristic.
type IStore interface {

	// User
	FindUserWithID(userID string) IUser
	FindUserWithClient(clientID string, clientSecret string) IUser
	FindUserWithCredential(username string, password string) IUser

	// Client
	FindClientWithID(clientID string) IClient
	FindClientWithCredential(clientID string, clientSecret string) IClient

	// Access Token
	FindAccessToken(token string) IToken
	FindAccessTokenWithCredential(clientID string, userID string) IToken
	CreateAccessToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) IToken
	DeleteAccessToken(token IToken)

	// Refresh Token
	FindRefreshToken(token string) IToken
	FindRefreshTokenWithCredential(clientID string, userID string) IToken
	CreateRefreshToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) IToken
	DeleteRefreshToken(token IToken)

	//	// Authorization code
	//	FindAuthorizationCode(authorizationCode string)
	//	SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// IToken describes a token's characteristic, it can be either access token or refresh token.
type IToken interface {

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

////////////////////////////////////////////////////////////////////////////////////////////////////
// TokenResponse describes a granted response that will be returned to client.
type TokenResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`

	Roles []string `json:"roles,omitempty"`
}
