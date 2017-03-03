package oauth2

import "gopkg.in/mgo.v2/bson"

// MongoDBUser describes a mongodb user.
type MongoDBUser struct {
	ID    bson.ObjectId `bson:"_id"`
	User  string        `bson:"username,omitempty"`
	Pass  string        `bson:"password,omitempty"`
	Roles []string      `bson:"roles,omitempty"`

	FacebookID    string `bson:"facebook_id,omitempty"`
	FacebookToken string `bson:"facebook_token,omitempty"`
}

// UserID returns user_id.
func (a *MongoDBUser) UserID() string {
	return a.ID.Hex()
}

// Username returns user's username.
func (a *MongoDBUser) Username() string {
	return a.User
}

// Password returns password.
func (a *MongoDBUser) Password() string {
	return a.Pass
}

// UserRoles returns user's roles.
func (a *MongoDBUser) UserRoles() []string {
	return a.Roles
}
