package oauth2

import "gopkg.in/mgo.v2/bson"

// DefaultUser describes a mongodb user.
type DefaultUser struct {
	ID    bson.ObjectId `bson:"_id"`
	User  string        `bson:"username,omitempty"`
	Pass  string        `bson:"password,omitempty"`
	Roles []string      `bson:"roles,omitempty"`

	FacebookID    string `bson:"facebook_id,omitempty"`
	FacebookToken string `bson:"facebook_token,omitempty"`
}

// UserID returns user_id.
func (a *DefaultUser) UserID() string { return a.ID.Hex() }

// Username returns user's username.
func (a *DefaultUser) Username() string { return a.User }

// Password returns password.
func (a *DefaultUser) Password() string { return a.Pass }

// UserRoles returns user's roles.
func (a *DefaultUser) UserRoles() []string { return a.Roles }
