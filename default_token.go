package oauth2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// DefaultToken descripts a mongodb Token.
type DefaultToken struct {
	TokenID     bson.ObjectId `bson:"_id,omitempty"`
	UserID      bson.ObjectId `bson:"user_id,omitempty"`
	ClientID    bson.ObjectId `bson:"client_id,omitempty"`
	CreatedTime time.Time     `bson:"created_time,omitempty"`
	ExpiredTime time.Time     `bson:"expired_time,omitempty"`
}

// GetClientID returns client_id.
func (t *DefaultToken) GetClientID() string { return t.ClientID.Hex() }

// GetUserID returns user_id.
func (t *DefaultToken) GetUserID() string { return t.UserID.Hex() }

// GetToken returns token.
func (t *DefaultToken) GetToken() string { return t.TokenID.Hex() }

// IsExpired validate if this token is expired or not.
func (t *DefaultToken) IsExpired() bool { return time.Now().Unix() >= t.ExpiredTime.Unix() }

// GetCreatedTime returns created_time.
func (t *DefaultToken) GetCreatedTime() time.Time { return t.CreatedTime }

// GetExpiredTime returns expired_time.
func (t *DefaultToken) GetExpiredTime() time.Time { return t.ExpiredTime }
