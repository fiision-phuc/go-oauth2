package oauth2

import "time"

// IStore descripts a token storage's characteristic.
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
