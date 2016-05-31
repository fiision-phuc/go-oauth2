package oauth2

import "gopkg.in/mgo.v2/bson"

// DefaultClient descripts a mongodb client.
type DefaultClient struct {
	ID        bson.ObjectId `bson:"_id"`
	Secret    bson.ObjectId `bson:"client_secret"`
	Grants    []string      `bson:"grant_types,omitempty"`
	Redirects []string      `bson:"redirect_uris,omitempty"`
}

// ClientID returns client_id.
func (a *DefaultClient) ClientID() string { return a.ID.Hex() }

// ClientSecret returns client_secret.
func (a *DefaultClient) ClientSecret() string { return a.Secret.Hex() }

// GrantTypes returns grant_types.
func (a *DefaultClient) GrantTypes() []string { return a.Grants }

// RedirectURIs returns redirect_uris.
func (a *DefaultClient) RedirectURIs() []string { return a.Redirects }
