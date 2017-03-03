package oauth2

import "time"

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
