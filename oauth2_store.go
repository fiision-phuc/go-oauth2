package oauth2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Store interface {

	// User
	FindUserWithID(userID bson.ObjectId) *User
	FindUserWithCredential(username string, password string) *User

	// Client
	FindClientWithCredential(clientID string, clientSecret string) *Client

	// Access token
	FindAccessToken(accessToken string) *Token
	FindAccessTokenWithCredential(clientID string, userID bson.ObjectId) *Token
	DeleteAccessToken(accessToken *Token)
	SaveAccessToken(accessToken *Token)

	// Refresh token
	FindRefreshToken(refreshToken string) *Token
	FindRefreshTokenWithCredential(clientID string, userID bson.ObjectId) *Token
	DeleteRefreshToken(refreshToken *Token)
	SaveRefreshToken(refreshToken *Token)

	// Authorization code
	FindAuthorizationCode(authorizationCode string)
	SaveAuthorizationCode(authorizationCode string, clientID string, expires time.Time)
}

//func (m *Model) GetAccessToken(accessToken string) *AccessToken {
//	return nil
//}
//func (m *Model) GetRefreshToken(refreshToken string) *RefreshToken {
//	return nil
//}
//func (m *Model) GetClient(clientId string, clientSecret string) *Client {
//	return nil
//}
//func (m *Model) GrantTypeAllowed(clientId string, grantType string) bool {
//	return nil
//}
//func (m *Model) GetUser(username string, password string) *User {
//	return nil
//}

//func (m *Model) SaveAccessToken(token string, clientId string, expires time.Time, userId string) {
//	accessToken := AccessToken{
//		AccessToken: token,
//		ClientId:    clientId,
//		UserId:      userId,
//		Expires:     expires,
//	}
//	// accessToken.save(callback)
//}
//func (m *Model) SaveRefreshToken(token string, clientId string, expires time.Time, userId string) {
//	refreshToken := RefreshToken{
//		RefreshToken: token,
//		ClientId:     clientId,
//		UserId:       userId,
//		Expires:      expires,
//	}
//	// accessToken.save(callback)
//}
