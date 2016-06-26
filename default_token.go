package oauth2

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"gopkg.in/mgo.v2/bson"
)

// DefaultToken describes a mongodb Token.
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
func (t *DefaultToken) Token() string {
	token := jwt.New(jwt.SigningMethodRS256)

	// Set some claims
	createdTime, _ := t.Created.MarshalText()
	expiredTime, _ := t.Expired.MarshalText()
	token.Claims = jwt.MapClaims{
		"_id":          t.ID.Hex(),
		"user_id":      t.User.Hex(),
		"client_id":    t.Client.Hex(),
		"created_time": string(createdTime),
		"expired_time": string(expiredTime),
	}

	// Generate token
	tokenString, _ := token.SignedString(privateKey)
	return tokenString
}

// IsExpired validate if this token is expired or not.
func (t *DefaultToken) IsExpired() bool { return time.Now().UTC().Unix() >= t.Expired.Unix() }

// CreatedTime returns created_time.
func (t *DefaultToken) CreatedTime() time.Time { return t.Created }

// ExpiredTime returns expired_time.
func (t *DefaultToken) ExpiredTime() time.Time { return t.Expired }
