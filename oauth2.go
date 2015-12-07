package oauth2

import "time"

// AuthClient descripts an authenticated client.
type AuthClient interface {
	GetClientID() string
	GetClientSecret() string

	GetGrantType() string      // Client side only
	GetRedirectURI() string    // Client side only
	GetGrantTypes() []string   // Server side only
	GetRedirectURIs() []string // Server side only
}

/////////////////////////////////////////////////////////////////////////////////////////////////

// AuthUser descripts an authenticated user.
type AuthUser interface {
	GetUserID() string
	GetUsername() string
	GetPassword() string

	GetUserRoles() []string
}

/////////////////////////////////////////////////////////////////////////////////////////////////

// Token descripts a token, it can be either access token or refresh token.
type Token interface {
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

/////////////////////////////////////////////////////////////////////////////////////////////////

// TokenStore descripts a token storage.
type TokenStore interface {

	// User
	FindUserWithID(userID string) AuthUser
	FindUserWithClient(clientID string, clientSecret string) AuthUser
	FindUserWithCredential(username string, password string) AuthUser

	// Client
	FindClientWithCredential(clientID string, clientSecret string) AuthClient

	// Token
	FindToken(accessToken string) Token
	FindTokenWithCredential(clientID string, userID string) Token
	CreateToken(clientID string, userID string, token string, createdTime time.Time, expiredTime time.Time) Token
	DeleteToken(accessToken Token)
	SaveToken(accessToken Token)

	//	// Authorization code
	//	FindAuthorizationCode(authorizationCode string)
	//	SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time)
}

/////////////////////////////////////////////////////////////////////////////////////////////////

// TokenResponse descripts a granted response that will be returned to client.
type TokenResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}