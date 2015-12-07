package oauth2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// AuthClientDefault descripts a mongodb AuthClient document.
type AuthClientDefault struct {
	ID bson.ObjectId `bson:"_id"`

	ClientID     string   `bson:"client_id" inject:"client_id"`
	ClientSecret string   `bson:"client_secret" inject:"client_secret"`
	GrantType    string   `bson:"-" inject:"grant_type"`
	RedirectURI  string   `bson:"-" inject:"redirect_uri"`
	GrantTypes   []string `bson:"grant_types"`
	RedirectURIs []string `bson:"redirect_uris"`
}

func createAuthClientDefault(c *RequestContext) AuthClient {
	client := AuthClientDefault{}
	c.BindForm(&client)

	if len(client.ClientID) == 0 || len(client.ClientSecret) == 0 {
		username, password, ok := c.BasicAuth()

		if ok {
			client.ClientID = username
			client.ClientSecret = password
		}
	}
	return &client
}

// GetClientID returns client_id.
func (a *AuthClientDefault) GetClientID() string { return a.ClientID }

// GetClientSecret returns client_secret.
func (a *AuthClientDefault) GetClientSecret() string { return a.ClientSecret }

// GetGrantType returns grant_type.
func (a *AuthClientDefault) GetGrantType() string { return a.GrantType }

// GetRedirectURI returns redirect_uri.
func (a *AuthClientDefault) GetRedirectURI() string { return a.RedirectURI }

// GetGrantTypes returns grant_types.
func (a *AuthClientDefault) GetGrantTypes() []string { return a.GrantTypes }

// GetRedirectURIs returns redirect_uris.
func (a *AuthClientDefault) GetRedirectURIs() []string { return a.RedirectURIs }

/////////////////////////////////////////////////////////////////////////////////////////////////

// AuthUserDefault descripts a mongodb AuthUser document.
type AuthUserDefault struct {
	UserID   bson.ObjectId `bson:"_id,omitempty"`
	Username string        `bson:"username" inject:"username"`
	Password string        `bson:"password" inject:"password"`
	Roles    []string      `bson:"roles"`
}

// GetUserID returns user_id.
func (a *AuthUserDefault) GetUserID() string { return a.UserID.Hex() }

// GetUsername returns username.
func (a *AuthUserDefault) GetUsername() string { return a.Username }

// GetPassword returns password.
func (a *AuthUserDefault) GetPassword() string { return a.Password }

// GetUserRoles returns roles.
func (a *AuthUserDefault) GetUserRoles() []string { return a.Roles }

/////////////////////////////////////////////////////////////////////////////////////////////////

// TokenDefault descripts a mongodb Token document.
type TokenDefault struct {
	TokenID     bson.ObjectId `bson:"_id"`
	UserID      bson.ObjectId `bson:"user_id"`
	ClientID    string        `bson:"client_id"`
	Token       string        `bson:"token"`
	CreatedTime time.Time     `bson:"created_time"`
	ExpiredTime time.Time     `bson:"expired_time"`
}

// GetClientID returns client_id.
func (t *TokenDefault) GetClientID() string { return t.TokenID.Hex() }

// GetUserID returns user_id.
func (t *TokenDefault) GetUserID() string { return t.UserID.Hex() }

// GetToken returns token.
func (t *TokenDefault) GetToken() string { return t.Token }

// SetToken updates token.
func (t *TokenDefault) SetToken(token string) { t.Token = token }

// isExpired validate if this token is expired or not.
func (t *TokenDefault) IsExpired() bool { return time.Now().Unix() >= t.ExpiredTime.Unix() }

// GetCreatedTime returns created_time.
func (t *TokenDefault) GetCreatedTime() time.Time { return t.CreatedTime }

// SetCreatedTime updates created_time.
func (t *TokenDefault) SetCreatedTime(createdTime time.Time) { t.CreatedTime = createdTime }

// GetExpiredTime returns expired_time.
func (t *TokenDefault) GetExpiredTime() time.Time { return t.ExpiredTime }

// SetExpiredTime updates expired_time.
func (t *TokenDefault) SetExpiredTime(expiredTime time.Time) { t.ExpiredTime = expiredTime }
