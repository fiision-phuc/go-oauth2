package oauth2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

////////////////////////////////////////////////////////////////////////////////
// Access Token																  //
////////////////////////////////////////////////////////////////////////////////
type AccessToken struct {
	TokenID     bson.ObjectId `bson:"_id,omitempty" json:"-"`
	UserID      bson.ObjectId `bson:"user_id,omitempty" json:"-"`
	ClientID    bson.ObjectId `bson:"client_id,omitempty" json:"-"`
	AccessToken string        `bson:"access_token,omitempty" json:"access_token,omitempty"`
	CreatedTime time.Time     `bson:"created_time,omitempty" json:"-"`
	ExpiredTime time.Time     `bson:"expired_time,omitempty" json:"-"`
}

////////////////////////////////////////////////////////////////////////////////
// Access Token																  //
////////////////////////////////////////////////////////////////////////////////
type Client struct {
	ClientID     string   `bson:"_id,omitempty" json:"client_id,omitempty" inject:"client_id"`
	ClientSecret string   `bson:"secret,omitempty" json:"client_secret,omitempty" inject:"client_secret"`
	GrantTypes   []string `bson:"grant_types,omitempty" json:"-"`
	RedirectURIs []string `bson:"redirect_uris,omitempty" json:"-"`

	GrantType   string `bson:"-" json:"grant_type,omitempty" inject:"grant_type"`
	RedirectURI string `bson:"-" json:"redirect_uri,omitempty" inject:"redirect_uri"`
}

func createClient(c *Context) *Client {
	client := Client{}
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

////////////////////////////////////////////////////////////////////////////////
// Access Token																  //
////////////////////////////////////////////////////////////////////////////////
type RefreshToken struct {
	TokenID      bson.ObjectId `bson:"_id,omitempty"`
	UserID       bson.ObjectId `bson:"user_id,omitempty"`
	ClientID     bson.ObjectId `bson:"client_id,omitempty"`
	RefreshToken string        `bson:"refresh_token,omitempty"`
	CreatedTime  time.Time     `bson:"created_time,omitempty"`
	ExpiredTime  time.Time     `bson:"expired_time,omitempty"`
}

////////////////////////////////////////////////////////////////////////////////
// Access Token																  //
////////////////////////////////////////////////////////////////////////////////
type User struct {
	UserID   bson.ObjectId `bson:"_id,omitempty" json:"user_id,omitempty"`
	Username string        `bson:"username,omitempty" json:"username,omitempty" inject:"username"`
	Password string        `bson:"password,omitempty" json:"password,omitempty" inject:"password"`
	Roles    []string      `bson:"roles,omitempty" json:"roles,omitempty"`
}
