package oauth2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// DefaultToken descripts a mongodb Token.
type DefaultToken struct {
	ID      bson.ObjectId `bson:"_id"`
	User    bson.ObjectId `bson:"user_id,omitempty"`
	Client  bson.ObjectId `bson:"client_id,omitempty"`
	Created time.Time     `bson:"created_time,omitempty"`
	Expired time.Time     `bson:"expired_time,omitempty"`
}

// ClientID returns client_id.
func (t *DefaultToken) ClientID() string { return t.Client.Hex() }

// UserID returns user_id.
func (t *DefaultToken) UserID() string { return t.User.Hex() }

// Token returns token.
func (t *DefaultToken) Token() string { return t.ID.Hex() }

// IsExpired validate if this token is expired or not.
func (t *DefaultToken) IsExpired() bool { return time.Now().Unix() >= t.Expired.Unix() }

// CreatedTime returns created_time.
func (t *DefaultToken) CreatedTime() time.Time { return t.Created }

// ExpiredTime returns expired_time.
func (t *DefaultToken) ExpiredTime() time.Time { return t.Expired }
