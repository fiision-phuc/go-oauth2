package oauth2

import "time"

// AuthClient descripts a characteristic of an authenticated client.
type AuthClient interface {
	GetClientID() string
	GetClientSecret() string
	GetGrantTypes() []string   // Server side only
	GetRedirectURIs() []string // Server side only
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// AuthUser descripts a characteristic of an authenticated user.
type AuthUser interface {
	GetUserID() string
	GetUsername() string
	GetPassword() string

	GetUserRoles() []string
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// Token descripts a characteristic of a token, it can be either access token or refresh token.
type Token interface {
	GetTokenID() string

	GetClientID() string
	GetUserID() string

	GetToken() string
	SetToken(token string)

	IsExpired() bool

	GetCreatedTime() time.Time
	SetCreatedTime(createdTime time.Time)

	GetExpiredTime() time.Time
	SetExpiredTime(expiredTime time.Time)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// TokenStore descripts a characteristic of a token storage.
type TokenStore interface {

	// User
	FindUserWithID(userID string) AuthUser
	FindUserWithClient(clientID string, clientSecret string) AuthUser
	FindUserWithCredential(username string, password string) AuthUser

	// Client
	FindClientWithID(clientID string) AuthClient
	FindClientWithCredential(clientID string, clientSecret string) AuthClient

	// Access Token
	FindAccessToken(token string) Token
	FindAccessTokenWithCredential(clientID string, userID string) Token
	CreateAccessToken(clientID string, userID string, token string, createdTime time.Time, expiredTime time.Time) Token
	DeleteAccessToken(token Token)
	SaveAccessToken(token Token)

	// Refresh Token
	FindRefreshToken(token string) Token
	FindRefreshTokenWithCredential(clientID string, userID string) Token
	CreateRefreshToken(clientID string, userID string, token string, createdTime time.Time, expiredTime time.Time) Token
	DeleteRefreshToken(token Token)
	SaveRefreshToken(token Token)

	//	// Authorization code
	//	FindAuthorizationCode(authorizationCode string)
	//	SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// TokenResponse descripts a granted response that will be returned to client.
type TokenResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
