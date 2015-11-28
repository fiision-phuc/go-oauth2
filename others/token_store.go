package oauth2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type OAuth2User struct {
	//	Id       bson.ObjectId `bson:"_id,omitempty" json:"user_id,omitempty"`
	Username string `bson:"username,omitempty" json:"username,omitempty"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`

	//	FirstName string `bson:"first_name,omitempty" json:"first_name,omitempty"`
	//	LastName  string `bson:"last_name,omitempty" json:"last_name,omitempty"`
	//	Email     string `bson:"email,omitempty" json:"email,omitempty"`
}

// type OAuth2Client struct {
// 	ClientId     string `bson:"_id,omitempty" json:"client_id,omitempty"`
// 	ClientSecret string `bson:"secret,omitempty" json:"client_secret,omitempty"`

// 	RedirectUri []string `bson:"redirect_uri,omitempty" json:"redirect_uri,omitempty"`
// }
type OAuth2AccessToken struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"-"`
	UserId   bson.ObjectId `bson:"user_id,omitempty" json:"-"`
	ClientId bson.ObjectId `bson:"client_id,omitempty" json:"-"`

	AccessToken string    `bson:"access_token,omitempty" json:"access_token,omitempty"`
	ExpiresIn   uint      `bson:"expires_in,omitempty" json:"expires_in,omitempty"`
	CreatedTime time.Time `bson:"created_time,omitempty" json:"-"`
	ExpiredTime time.Time `bson:"expired_time,omitempty" json:"-"`
}
type OAuth2RefreshToken struct {
	Id       string `bson:"_id,omitempty" json:"-"`
	UserId   string `bson:"user_id,omitempty" json:"-"`
	ClientId string `bson:"client_id,omitempty" json:"-"`

	RefreshToken string    `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	ExpiresIn    uint      `bson:"expires_in,omitempty" json:"expires_in,omitempty"`
	CreatedTime  time.Time `bson:"created_time,omitempty" json:"-"`
	ExpiredTime  time.Time `bson:"expired_time,omitempty" json:"-"`
}

type TokenStore interface {
	FindUser(username string, password string) *OAuth2User
	FindClient(clientId string, clientSecret string, redirectUri string) *OAuth2Client

	FindAccessToken(accessToken string) *OAuth2AccessToken
	FindRefreshToken(refreshToken string) *OAuth2RefreshToken

	ValidateGrantType(clientId string, grantType string) bool

	FindAuthorizationCode(authorizationCode string)
	SaveAuthorizationCode(authorizationCode string, clientId string, expires time.Time)

	SaveAccessToken(token string, clientId string, expires time.Time, userId string)
	SaveRefreshToken(token string, clientId string, expires time.Time, userId string)
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
