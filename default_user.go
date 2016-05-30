package oauth2

import "gopkg.in/mgo.v2/bson"

// DefaultUser descripts a mongodb user.
type DefaultUser struct {
	UserID     bson.ObjectId `bson:"_id,omitempty"`
	FacebookID string        `bson:"facebook_id,omitempty"`

	Username string   `bson:"username,omitempty"`
	Password string   `bson:"password,omitempty"`
	Roles    []string `bson:"roles,omitempty"`
}

// GetUserID returns user_id.
func (a *DefaultUser) GetUserID() string { return a.UserID.Hex() }

// GetUsername returns username.
func (a *DefaultUser) GetUsername() string { return a.Username }

// GetPassword returns password.
func (a *DefaultUser) GetPassword() string { return a.Password }

// GetUserRoles returns roles.
func (a *DefaultUser) GetUserRoles() []string { return a.Roles }
