package oauth2

import "time"

// TokenStore describes a token store's characteristic.
type TokenStore interface {

	// FindUserWithID returns an user entity according to userID or null. A user entity can either
	// human or machine.
	//
	// @param
	// - userID {string} (userID that associated with user's entity)
	//
	// @return
	// - user {User} (an user entity or null)
	FindUserWithID(userID string) User

	// FindUserWithClient returns a machine user entity.
	//
	// @param
	// - clientID {string} (client's client_id)
	// - clientSecret {string} (client's client_secret)
	//
	// @return
	// - user {User} (a machine user entity or null)
	FindUserWithClient(clientID string, clientSecret string) User

	// FindUserWithCredential returns a human user entity.
	//
	// @param
	// - username {string} (user's username)
	// - password {string} (user's password)
	//
	// @return
	// - user {User} (a human user entity or null)
	FindUserWithCredential(username string, password string) User

	// FindClientWithID returns a client entity according to clientID or null.
	//
	// @param
	// - clientID {string} (client's client_id)
	//
	// @return
	// - client {Client} (a client entity or null)
	FindClientWithID(clientID string) Client

	// FindClientWithCredential returns a client entity according to clientID and clientSecret or
	// null.
	//
	// @param
	// - clientID {string} (client's client_id)
	// - clientSecret {string} (client's client_secret)
	//
	// @return
	// - client {Client} (a client entity or null)
	FindClientWithCredential(clientID string, clientSecret string) Client

	// FindAccessToken returns an access token entity according to token string or null.
	//
	// @param
	// - token {string} (user's access token in string form)
	//
	// @return
	// - token {Token} (a token's instance or null)
	FindAccessToken(token string) Token

	// FindAccessTokenWithCredential returns an access token entity according to clientID and
	// userID or null.
	//
	// @param
	// - clientID {string} (client's client_id)
	// - userID {string} (userID that associated with user's entity)
	//
	// @return
	// - token {Token} (a token's instance or null)
	FindAccessTokenWithCredential(clientID string, userID string) Token

	// CreateAccessToken create a token's instance.
	//
	// @param
	// - clientID {string} (client's client_id)
	// - userID {string} (userID that associated with user's entity)
	// - createdTime {time.Time} (token's issued time)
	// - expiredTime {time.Time} (token's expired time)
	//
	// @return
	// - token {Token} (an access token's instance)
	CreateAccessToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token

	// DeleteAccessToken deletes an access token from database.
	//
	// @param
	// - token {Token} (an access token's instance)
	DeleteAccessToken(token Token)

	// FindRefreshToken returns a refresh token entity according to token string or null.
	//
	// @param
	// - token {string} (user's refresh token in string form)
	//
	// @return
	// - token {Token} (a token's instance or null)
	FindRefreshToken(token string) Token

	// FindRefreshTokenWithCredential returns a refresh token entity according to clientID and
	// userID or null.
	//
	// @param
	// - clientID {string} (client's client_id)
	// - userID {string} (userID that associated with user's entity)
	//
	// @return
	// - token {Token} (a token's instance or null)
	FindRefreshTokenWithCredential(clientID string, userID string) Token

	// CreateRefreshToken create a token's instance.
	//
	// @param
	// - clientID {string} (client's client_id)
	// - userID {string} (userID that associated with user's entity)
	// - createdTime {time.Time} (token's issued time)
	// - expiredTime {time.Time} (token's expired time)
	//
	// @return
	// - token {Token} (a refresh token's instance)
	CreateRefreshToken(clientID string, userID string, createdTime time.Time, expiredTime time.Time) Token

	// DeleteRefreshToken deletes a refresh token from database.
	//
	// @param
	// - token {Token} (a refresh token's instance)
	DeleteRefreshToken(token Token)

	//	// Authorization code
	//	FindAuthorizationCode(authorizationCode string)
	//	SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time)
}
